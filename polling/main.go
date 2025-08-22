package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func newCon() *sql.DB {
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:6306)/demo_12")
	if err != nil {
		log.Panic(err)
	}
	return db
}

type conn struct {
	db *sql.DB
}

type cpool struct {
	conns    []*conn
	mu       *sync.Mutex
	channel  chan interface{}
	maxConns int
}

func NewCPool(maxConns int) (*cpool, error) {
	var mu = sync.Mutex{}
	pool := cpool{
		mu:       &mu,
		conns:    make([]*conn, 0, maxConns),
		maxConns: maxConns,
		channel:  make(chan interface{}, maxConns),
	}

	for i := 1; i <= maxConns; i++ {
		pool.conns = append(pool.conns, &conn{newCon()})
		pool.channel <- nil
	}
	return &pool, nil
}

func (pool *cpool) close() {
	close(pool.channel)
	for i := range pool.conns {
		pool.conns[i].db.Close()
	}
}

func (pool *cpool) get() (*conn, error) {
	<-pool.channel
	pool.mu.Lock()
	c := pool.conns[0]
	pool.conns = pool.conns[1:]
	pool.mu.Unlock()
	return c, nil
}

func (pool *cpool) put(c *conn) {
	pool.mu.Lock()
	pool.conns = append(pool.conns, c)
	pool.mu.Unlock()
	pool.channel <- nil
}

func nonpoolbenchmark() {
	startTime := time.Now()
	var wg sync.WaitGroup
	for i := 1; i <= 200; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			db := newCon()
			_, err := db.Exec("SELECT SLEEP(0.01) ")
			if err != nil {
				log.Panic(err)
			}
			db.Close()
		}()
	}
	wg.Wait()
	fmt.Println("Benchmark time", time.Since(startTime))
}

func poolbenchmark() {
	startTime := time.Now()
	var wg sync.WaitGroup
	var pool, err = NewCPool(10)
	if err != nil {
		log.Panic(err)
	}

	for i := 1; i <= 200; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c, err := pool.get()
			if err != nil {
				log.Panic(err)
			}
			_, err = c.db.Exec("SELECT SLEEP(0.01) ")
			if err != nil {
				log.Panic(err)
			}
			pool.put(c)
		}()
	}
	wg.Wait()
	fmt.Println("Benchmark time", time.Since(startTime))
	pool.close()
}

func main() {
	//nonpoolbenchmark()
	poolbenchmark()
}
