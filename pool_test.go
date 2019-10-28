/**
 * @description goroutinepool
 * @author ibingbo
 * @date 2018/12/1
 */
package goroutinepool

import (
	"context"
	"fmt"
	"github.com/panjf2000/ants"
	"goroutinepool/channel"
	"goroutinepool/simple"
	"goroutinepool/woker"
	"log"
	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"sync"
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
func TestChannelPool(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("don't worry, i can take care this")
		}
	}()
	c := make(chan bool)
	go func() {
		for {
			select {
			case <-c:
				return
			default:
				log.Println("***goroutine num***", runtime.NumGoroutine())
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	time.Sleep(1 * time.Second)

	pool := channel.NewGoroutinePoolWithCatch(10, 100, true)
	for i := 0; i < 100; i++ {
		id := i
		worker := &woker.Worker{func() {
			time.Sleep(1 * time.Second)
			log.Println("***GoroutineID", GoID(), id, "worker done...")
			if id == 50 {
				panic("test error")
			}
		}}
		pool.Submit(worker)
	}
	pool.AwaitTermination()
	pool.Close()

	// time.Sleep(2 * time.Second)
	//
	// pool1 := channel.NewGoroutinePool(10, 100)
	// for i := 0; i < 100; i++ {
	// 	id := i
	// 	worker := &woker.Worker{func() {
	// 		time.Sleep(1 * time.Second)
	// 		log.Println(id, "worker done...")
	// 	}}
	// 	pool1.Submit(worker)
	// }
	// pool1.AwaitTermination()
	// pool1.Close()
	// time.Sleep(2 * time.Second)

	c <- true
	time.Sleep(1 * time.Second)
	close(c)

	time.Sleep(2 * time.Second)
	log.Println("***goroutine num***", runtime.NumGoroutine())
}

func TestSimpleChannelPool(t *testing.T) {
	pool := channel.NewGoroutinePool(4, 100)
	defer pool.Close()
	var arr []string
	titleArr := []string{"111", "222"}
	descArr := []string{"****", "---"}
	rand.Seed(time.Now().Unix())
	for i := 0; i < 500; i++ {
		arr = append(arr, fmt.Sprintf("%d", i))
		if len(arr) == 20 {
			arr1 := arr
			idx := rand.Intn(2)
			title := titleArr[idx]
			desc := descArr[idx]
			worker := &woker.Worker{func() {
				err := sendMessage(strings.Join(arr1, ","), title, desc)
				log.Println("err:", err)
			}}
			pool.Submit(worker)
			arr = []string{}
		}

	}
	pool.AwaitTermination()
}

func sendMessage(uids string, title string, description string) error {
	time.Sleep(1 * time.Second)
	log.Printf("title:%s description:%s uids:%s", title, description, uids)
	return nil
}

func TestSomething(t *testing.T) {
	titleArr := []string{"111", "222"}
	descArr := []string{"****", "---"}
	rand.Seed(time.Now().Unix())
	for i := 0; i < 100; i++ {
		idx := rand.Intn(2)
		fmt.Println(time.Now().Format("2006/01/02 15:04:05"), idx, titleArr[idx], descArr[idx])
	}

}

func GoID() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}

var sum int32

func myFunc(i interface{}) {
	n := i.(int32)
	atomic.AddInt32(&sum, n)
	log.Printf("run with %d", n)
}
func demoFunc() {
	time.Sleep(10 * time.Microsecond)
	log.Println("hello world!")
}
func TestAnts(t *testing.T) {
	defer ants.Release()
	runTimes := 1000
	var wg sync.WaitGroup
	syncCalculateSum := func() {
		demoFunc()
		wg.Done()
	}
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		ants.Submit(syncCalculateSum)
	}
	wg.Wait()
	log.Printf("running goroutines: %d", ants.Running())
	log.Println("finish all tasks")

	p, _ := ants.NewPoolWithFunc(10, func(i interface{}) {
		myFunc(i)
		wg.Done()
	})
	defer p.Release()
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		p.Invoke(int32(i))
	}
	wg.Wait()
	log.Printf("running goroutines: %d", p.Running())
	log.Printf("finish all tasks,result is %d", sum)

	type TestStruct struct {
		Name string
		T    time.Time
	}
	ts := TestStruct{Name: "bill"}
	log.Println(ts, ts.T.Format("2006-01-02 15:01-02"))
}

func TestSimplePoolForAsync(t *testing.T) {
	type Item struct {
		IDs   []int
		Type  int
		Title string
	}
	inputChan := make(chan *Item, 2)

	go func() {
		// defer wg.Done()
		var arr []int
		for i := 0; i < 1000; i++ {
			arr = append(arr, i)
			if len(arr)%10 == 0 {
				arrCopy := arr
				t := rand.Intn(2)
				if t == 1 {
					inputChan <- &Item{IDs: arrCopy, Type: t, Title: "aa"}
				} else {
					inputChan <- &Item{IDs: arrCopy, Type: t, Title: "bb"}
				}
				// log.Println("********** input one number ", arrCopy)
				arr = []int{}
			}
		}
		close(inputChan)
		log.Println("----------- input all done ....")
	}()
	pool1 := simple.NewGoroutinePool(5, 10)
	// var wg sync.WaitGroup
	defer pool1.Close()
	// wg.Add(2)
	go func() {
		// defer wg.Done()
		for item := range inputChan {
			// log.Println("********** output one number ", id)
			idCopy := item
			worker := &woker.Worker{func() {
				time.Sleep(1 * time.Second)
				log.Println("*********** take one number ", idCopy)
			}}
			pool1.Submit(worker)
		}
		log.Println("---------- all take done")
	}()
	time.Sleep(1 * time.Second)
	// wg.Wait()
	pool1.AwaitTermination()
	log.Println("*********** all number taked")
}
