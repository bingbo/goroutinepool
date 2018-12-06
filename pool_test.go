/**
 * @description goroutinepool
 * @author ibingbo
 * @date 2018/12/1
 */
package goroutinepool

import (
	"context"
	"fmt"
	"log"
	"strconv"
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
	result:=NewConcurrentMap()
	pool := NewGoroutinePool(10, 30)
	for i := 0; i < 500; i++ {
		id := i
		worker := &Worker{func() {
			time.Sleep(1 * time.Second)
			result.Put(strconv.Itoa(id),true)
			log.Println(id, "worker done...")
		}}
		pool.Execute(worker)
	}
	pool.AwaitTermination()
	fmt.Println(result)
	fmt.Println(len(result.Map))
	var ids []string
	for key:=range result.Map{
		ids=append(ids,key)
	}
	fmt.Println(len(ids))

}

func TestPoolWithContext(t *testing.T) {
	context, cancel := context.WithCancel(context.Background())
	pool := NewGoroutinePool(10, 30)
	pool.stop = cancel
	for i := 0; i < 500; i++ {
		id := i
		worker := &Worker{func() {
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
		worker := &Worker{func() {
			time.Sleep(1 * time.Second)
			log.Println(id, "worker done...")
		}}
		pool.ExecuteWithContext(worker, context)
	}

	pool.AwaitTermination()

}
