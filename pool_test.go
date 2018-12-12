/**
 * @description goroutinepool
 * @author ibingbo
 * @date 2018/12/1
 */
package goroutinepool

import (
	"context"
	"fmt"
	"goroutinepool/simple"
	"goroutinepool/woker"
	"log"
	"runtime"
	"strconv"
	"sync/atomic"
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
	result := NewConcurrentMap()
	pool := NewGoroutinePool(10, 30)
	for i := 0; i < 8000; i++ {
		id := i
		worker := &woker.Worker{func() {
			time.Sleep(1 * time.Second)
			result.Put(strconv.Itoa(id), true)
			log.Println(id, "worker done...")
		}}
		pool.Execute(worker)
	}
	pool.AwaitTermination()
	fmt.Println(result)
	fmt.Println(len(result.Map))
	var ids []string
	for key := range result.Map {
		ids = append(ids, key)
	}
	fmt.Println(len(ids))

}

func TestPoolWithContext(t *testing.T) {
	context, cancel := context.WithCancel(context.Background())
	pool := NewGoroutinePool(10, 30)
	pool.stop = cancel
	for i := 0; i < 500; i++ {
		id := i
		worker := &woker.Worker{func() {
			time.Sleep(1 * time.Second)
			log.Println(id, "worker done...")
		}}
		pool.ExecuteWithContext(worker, context)
		if id%450 == 0 {
			pool.stop()
		}
	}

	pool.AwaitTermination()

}

func TestPoolWithContext2(t *testing.T) {
	context, _ := context.WithTimeout(context.Background(), 10*time.Second)
	pool := NewGoroutinePool(10, 30)
	for i := 0; i < 500; i++ {
		id := i
		worker := &woker.Worker{func() {
			time.Sleep(1 * time.Second)
			log.Println(id, "worker done...")
		}}
		pool.ExecuteWithContext(worker, context)
	}

	pool.AwaitTermination()

}

func TestSom(t *testing.T) {
	var num int32 = 10
	for {
		fmt.Println("num = ", num)
		if atomic.CompareAndSwapInt32(&num, 0, 0) {
			fmt.Println("num == 0 already")
			break
		} else {
			atomic.AddInt32(&num, -1)
		}

		time.Sleep(1 * time.Second)

	}
}

func TestSimplePool(t *testing.T) {
	c := make(chan bool)
	go func() {
		for {
			select {
			case <-c:
				return
			default:
				log.Println("***goroutine num***", runtime.NumGoroutine())
			}
		}
	}()

	pool := simple.NewGoroutinePool(10, 100)
	for i := 0; i < 100; i++ {
		id := i
		worker := &woker.Worker{func() {
			time.Sleep(1 * time.Second)
			log.Println(id, "worker done...")
		}}
		pool.Submit(worker)
	}
	pool.AwaitTermination()
	pool.Close()

	time.Sleep(2 * time.Second)

	pool1 := simple.NewGoroutinePool(10, 100)
	for i := 0; i < 100; i++ {
		id := i
		worker := &woker.Worker{func() {
			time.Sleep(1 * time.Second)
			log.Println(id, "worker done...")
		}}
		pool1.Submit(worker)
	}
	pool1.AwaitTermination()
	pool1.Close()

	time.Sleep(2 * time.Second)

	pool2 := simple.NewGoroutinePool(10, 100)
	for i := 0; i < 100; i++ {
		id := i
		worker := &woker.Worker{func() {
			time.Sleep(1 * time.Second)
			log.Println(id, "worker done...")
		}}
		pool2.Submit(worker)
	}
	pool2.AwaitTermination()
	pool2.Close()

	time.Sleep(10 * time.Second)
	c <- true
	time.Sleep(1 * time.Second)
	close(c)

	time.Sleep(2 * time.Second)
	log.Println("***goroutine num***", runtime.NumGoroutine())
}
