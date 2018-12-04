/**
 * @description goroutinepool
 * @author ibingbo
 * @date 2018/12/1
 */
package goroutinepool

import (
	"context"
	"log"
	"testing"
	"time"
)

type MyWorker struct {
}

func (w *MyWorker) Work() {
	time.Sleep(1 * time.Second)
	log.Println("*****worker done.....")
}

func TestPool(t *testing.T) {
	pool := NewGoroutinePool(10, 30)
	for i := 0; i < 100; i++ {
		worker := &MyWorker{}
		pool.Execute(worker)
	}
	pool.AwaitTermination()

}

func TestPool1(t *testing.T) {
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

}

func TestPoolWithContext(t *testing.T) {
	context, cancel := context.WithCancel(context.Background())
	pool := NewGoroutinePool(10, 30)
	for i := 0; i < 500; i++ {
		id := i
		worker := &Worker{func() {
			time.Sleep(1 * time.Second)
			log.Println(id, "worker done...")
		}}
		pool.ExecuteWithContext(worker, context)
		if id%100 == 0 {
			cancel()
		}
	}

	pool.AwaitTermination()

}

func TestPoolWithContext2(t *testing.T) {
	context, _ := context.WithTimeout(context.Background(), 10*time.Second)
	pool := NewGoroutinePool(10, 30)
	for i := 0; i < 500; i++ {
		id := i
		worker := &Worker{func() {
			time.Sleep(1 * time.Second)
			log.Println(id, "worker done...")
		}}
		pool.ExecuteWithContext(worker, context)
	}

	pool.AwaitTermination()

}
