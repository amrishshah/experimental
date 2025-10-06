package main

import (
	"fmt"
	"log"

	"github.com/tecbot/gorocksdb"
)

func main() {
	// ---- Options ----
	opts := gorocksdb.NewDefaultOptions()
	defer opts.Destroy()
	opts.SetCreateIfMissing(true)

	// Block cache + Bloom filter (good for point-lookups)
	bbto := gorocksdb.NewDefaultBlockBasedTableOptions()
	defer bbto.Destroy()
	cache := gorocksdb.NewLRUCache(512 << 20) // 512MB
	defer cache.Destroy()
	filter := gorocksdb.NewBloomFilter(10)

	bbto.SetBlockCache(cache)
	bbto.SetFilterPolicy(filter)
	opts.SetBlockBasedTableFactory(bbto)

	// Compression (ZSTD if built with it)
	opts.SetCompression(gorocksdb.ZSTDCompression)

	// ---- Open DB ----
	db, err := gorocksdb.OpenDb(opts, "data/rocks-demo")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	wo := gorocksdb.NewDefaultWriteOptions()
	defer wo.Destroy()
	// Sync true => fsync on each write (safer, slower). Usually use default false + periodic Flush/WAL sync.
	wo.SetSync(false)

	ro := gorocksdb.NewDefaultReadOptions()
	defer ro.Destroy()

	// ---- Put ----
	if err := db.Put(wo, []byte("user:42"), []byte("Amrish")); err != nil {
		log.Fatal(err)
	}

	// ---- Get ----
	slice, err := db.Get(ro, []byte("user:42"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("GET user:42 =", string(slice.Data()))
	slice.Free()

	// ---- Batch (atomic) ----
	wb := gorocksdb.NewWriteBatch()
	wb.Put([]byte("k1"), []byte("v1"))
	wb.Put([]byte("k2"), []byte("v2"))
	//wb.Delete([]byte("user:42"))
	if err := db.Write(wo, wb); err != nil {
		log.Fatal(err)
	}
	wb.Destroy()

	// ---- Iterator (prefix/range scan) ----
	it := db.NewIterator(ro)
	defer it.Close()
	for it.Seek([]byte("k")); it.Valid(); it.Next() {
		fmt.Printf("%s => %s\n", it.Key().Data(), it.Value().Data())
		it.Key().Free()
		it.Value().Free()
	}
	if err := it.Err(); err != nil {
		log.Fatal(err)
	}

	// ---- Flush memtable to SST ----
	if err := db.Flush(gorocksdb.NewDefaultFlushOptions()); err != nil {
		log.Fatal(err)
	}
}
