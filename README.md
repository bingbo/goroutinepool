# goroutinepool

```go
pool := NewGoroutinePool(10, 30)
	for i := 0; i < 500; i++ {
		id := i
		worker := &Worker{func() {
			time.Sleep(1 * time.Second)
			log.Println(id, "worker done...")
		}}
		pool.Execute(worker)
	}
	pool.AwaitTermination()
```
