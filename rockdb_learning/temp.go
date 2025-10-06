package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Message represents a message from Kafka.
// Using a string key for simplicity, as map keys must be comparable.
type Message struct {
	Key   string
	Value string
}

// Batch represents a collection of messages to be written to the DB.
// The key is the Kafka message key, and the value is the message payload.
type Batch map[string]string

// BatchProcessor manages the two-map batching logic.
type BatchProcessor struct {
	// Configuration
	batchSize    int
	batchTimeout time.Duration

	// Internal state
	mu            sync.Mutex
	activeMap     Batch
	processingMap Batch
	batchChan     chan Batch
	wg            sync.WaitGroup
	shutdownCtx   context.Context
	cancelFunc    context.CancelFunc
}

// NewBatchProcessor creates and initializes a new processor.
func NewBatchProcessor(ctx context.Context, batchSize int, batchTimeout time.Duration) *BatchProcessor {
	// Use a buffered channel of size 1. This allows the consumer to send the
	// batch and continue immediately, even if the DB writer is busy with the previous one.
	batchChan := make(chan Batch, 1)

	// Create a cancellable context for graceful shutdown.
	shutdownCtx, cancelFunc := context.WithCancel(ctx)

	return &BatchProcessor{
		batchSize:     batchSize,
		batchTimeout:  batchTimeout,
		activeMap:     make(Batch, batchSize),
		processingMap: make(Batch, batchSize),
		batchChan:     batchChan,
		shutdownCtx:   shutdownCtx,
		cancelFunc:    cancelFunc,
	}
}

// Start begins the processor's goroutines for consuming and writing.
// It takes a channel of messages, simulating the Kafka consumer library's output.
func (p *BatchProcessor) Start(messageSource <-chan Message) {
	p.wg.Add(1)
	go p.dbWriterLoop()

	p.wg.Add(1)
	go p.consumerLoop(messageSource)
}

// consumerLoop is the main loop for collecting messages.
func (p *BatchProcessor) consumerLoop(messageSource <-chan Message) {
	defer p.wg.Done()
	ticker := time.NewTicker(p.batchTimeout)
	defer ticker.Stop()

	log.Println("Consumer loop started...")

	for {
		select {
		case msg, ok := <-messageSource:
			if !ok {
				log.Println("Message source channel closed. Flushing final batch.")
				p.flushBatch()
				return
			}

			// This is the only place we write to the activeMap
			p.mu.Lock()
			p.activeMap[msg.Key] = msg.Value
			shouldFlush := len(p.activeMap) >= p.batchSize
			p.mu.Unlock()

			if shouldFlush {
				// We reset the ticker because a size-based flush just occurred.
				ticker.Reset(p.batchTimeout)
				p.flushBatch()
			}

		case <-ticker.C:
			// Time-based flush
			p.flushBatch()

		case <-p.shutdownCtx.Done():
			// Shutdown signal received
			log.Println("Shutdown signal received in consumer. Flushing final batch.")
			p.flushBatch()
			return
		}
	}
}

// flushBatch performs the two-map swap and sends the batch to the writer.
func (p *BatchProcessor) flushBatch() {
	p.mu.Lock()
	defer p.mu.Unlock()

	// If there's nothing to flush, do nothing.
	if len(p.activeMap) == 0 {
		return
	}

	// The magic swap! This is extremely fast.
	p.activeMap, p.processingMap = p.processingMap, p.activeMap

	log.Printf("Flushing batch of size %d\n", len(p.processingMap))

	// Send the now-full processingMap to the DB writer.
	// This will block only if the writer is still busy AND the channel buffer is full.
	select {
	case p.batchChan <- p.processingMap:
		// The map has been sent. The consumer is free.
		// We re-assign a new empty map to p.processingMap for the next swap.
		// This ensures the DB writer has exclusive ownership of the batch data.
		p.processingMap = make(Batch, p.batchSize)
	case <-p.shutdownCtx.Done():
		log.Println("Shutdown during flush. Batch may be lost.")
		return
	}
}

// dbWriterLoop simulates the worker that writes batches to the database.
func (p *BatchProcessor) dbWriterLoop() {
	defer p.wg.Done()
	log.Println("DB writer loop started...")

	for batch := range p.batchChan {
		if len(batch) == 0 {
			continue
		}

		log.Printf("=> DB Writer: Received batch of %d records. Writing to DB...", len(batch))
		// Simulate DB work
		time.Sleep(200 * time.Millisecond)
		log.Printf("=> DB Writer: Write successful for %d records.", len(batch))

		// IMPORTANT: Clear the map to release memory. The map is now empty
		// and its underlying memory can be garbage collected. This map will NOT be reused.
		// The consumerLoop creates a new empty map for the next swap.
		for k := range batch {
			delete(batch, k)
		}
	}
	log.Println("DB writer loop finished.")
}

// Shutdown handles the graceful shutdown of the processor.
func (p *BatchProcessor) Shutdown() {
	log.Println("Initiating shutdown...")
	// Signal all goroutines to stop
	p.cancelFunc()

	// Close the channel. This will cause the dbWriterLoop's range to exit
	// once all buffered batches are processed.
	close(p.batchChan)

	// Wait for all goroutines to finish their cleanup.
	p.wg.Wait()
	log.Println("Shutdown complete.")
}

func main1() {
	// --- Setup ---
	// Create a context that listens for OS interrupt signals.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	batchSize := 10
	batchTimeout := 2 * time.Second

	processor := NewBatchProcessor(ctx, batchSize, batchTimeout)

	// --- Kafka Simulation ---
	// Create a channel to simulate Kafka messages coming in.
	kafkaMessages := make(chan Message, 100)
	processor.Start(kafkaMessages)

	// Goroutine to simulate a Kafka producer.
	go func() {
		defer close(kafkaMessages) // Close channel when done producing
		for i := 1; i <= 55; i++ {
			msg := Message{
				Key:   fmt.Sprintf("key-%d", i),
				Value: fmt.Sprintf("value-%d", i),
			}
			select {
			case kafkaMessages <- msg:
				fmt.Printf("Produced message %d\n", i)
				// Simulate variance in message arrival time
				time.Sleep(50 * time.Millisecond)
			case <-ctx.Done():
				fmt.Println("Producer stopping due to shutdown signal.")
				return
			}
		}
	}()

	// --- Wait for Shutdown ---
	<-ctx.Done() // Block here until a shutdown signal is received

	// --- Graceful Shutdown ---
	// We call stop() to remove the signal handler, allowing a second interrupt to force-exit
	stop()
	processor.Shutdown()
}
