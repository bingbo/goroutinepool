/**
 * @description goroutinepool
 * @author ibingbo
 * @date 2018/12/1
 */
package channel

import (
	"fmt"
	"goroutinepool/woker"
	"log"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

// 默认最大任务数
var defaultMaxWorkerCount int32 = 100

// 默认最大协程数
var defaultMaxGoroutineCount int32 = 10

// 协程池结构定义
type GoroutinePool struct {
	maxGoroutineCount int32               // 最大协程数
	maxWorkerCount    int32               // 最大任务数
	workerQueue       chan woker.Workable // 任务队列
	waitGroup         sync.WaitGroup
	catchedPanic      bool // 是否捕获每个线程中的panic，如果不捕获出现panic时会导致整个程序异常退出
}

/**
 * NewGoroutinePool
 * 实例化一个协程池
 */
func NewGoroutinePool(maxGoroutineCount int32, maxWorkerCount int32) *GoroutinePool {
	if maxGoroutineCount == 0 {
		maxGoroutineCount = defaultMaxGoroutineCount
	}
	if maxWorkerCount == 0 {
		maxWorkerCount = defaultMaxWorkerCount
	}
	pool := &GoroutinePool{
		maxGoroutineCount: maxGoroutineCount,
		workerQueue:       make(chan woker.Workable, maxWorkerCount),
	}
	pool.start()
	return pool
}

func NewGoroutinePoolWithCatch(maxGoroutineCount int32, maxWorkerCount int32, catchedPanic bool) *GoroutinePool {
	if maxGoroutineCount == 0 {
		maxGoroutineCount = defaultMaxGoroutineCount
	}
	if maxWorkerCount == 0 {
		maxWorkerCount = defaultMaxWorkerCount
	}
	pool := &GoroutinePool{
		maxGoroutineCount: maxGoroutineCount,
		workerQueue:       make(chan woker.Workable, maxWorkerCount),
		catchedPanic:      catchedPanic,
	}
	pool.start()
	return pool
}

/**
 * 提交并执行一个任务
 */
func (pool *GoroutinePool) Submit(worker woker.Workable) {
	pool.waitGroup.Add(1)
	pool.workerQueue <- worker
	log.Println("add one worker")
}

/**
 * 等待所有的任务执行完成
 */
func (pool *GoroutinePool) AwaitTermination() {
	pool.waitGroup.Wait()
}

/**
 * 协程池开始执行
 */
func (pool *GoroutinePool) start() {
	for i := 0; int32(i) < pool.maxGoroutineCount; i++ {
		go pool.doWork()
	}
}

/**
 * 开始处理所有的任务
 */
func (pool *GoroutinePool) doWork() {
	log.Println(fmt.Sprintf("goroutine %d is running", pool.goroutineID()))
	if pool.catchedPanic {
		// 捕获当前自己线程中的异常，因为不能被主线程捕获到，如果不处理panic异常则会导致整个程序异常退出
		defer pool.catchPanic()
	}
	for worker := range pool.workerQueue {
		worker.Work()
		pool.waitGroup.Done()
	}
}

/**
 * 关系线程池
 */
func (pool *GoroutinePool) Close() {
	close(pool.workerQueue)
}

/**
 * 捕获每个线程的panic异常，由于每个线程只能捕获自己的异常，而主线程捕获不到子线程异常，如果不处理panic异常则会导致整个程序panic异常退出
 */
func (pool *GoroutinePool) catchPanic() {
	if err := recover(); err != nil {
		buf := make([]byte, 10000)
		runtime.Stack(buf, false)
		log.Println(fmt.Sprintf("panic defered [%v] stack trace : %v", err, string(buf)))
		pool.waitGroup.Done()
	}
}

/**
 * 获取当前协程ID
 */
func (pool *GoroutinePool) goroutineID() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		log.Println(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}
