/**
 * @description goroutinepool
 * @author ibingbo
 * @date 2018/12/1
 */
package simple

import (
	"goroutinepool/woker"
	"log"
	"sync"
	"sync/atomic"
)

// 默认最大任务数
var defaultMaxWorkerCount int32 = 100

// 默认最大协程数
var defaultMaxGoroutineCount int32 = 10

// 协程池结构定义
type GoroutinePool struct {
	maxGoroutineCount     int32               // 最大协程数
	currentGoroutineCount int32               // 当前协程数
	maxWorkerCount        int32               // 最大任务数
	workerQueue           chan woker.Workable // 任务队列
	done                  chan struct{}       // 任务状态，是否所有任务已完成
	waitGroup             sync.WaitGroup
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
		maxGoroutineCount:     maxGoroutineCount,
		workerQueue:           make(chan woker.Workable, maxWorkerCount),
		done:                  make(chan struct{}, maxGoroutineCount), // 保证每个goroutine都能收到关系信号
		currentGoroutineCount: maxGoroutineCount,
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
	// log.Println("add one worker")
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
	for {
		// log.Println("goroutine count ", runtime.NumGoroutine())
		select {
		case <-pool.done:
			log.Println("close goroutine....")
			if atomic.CompareAndSwapInt32(&pool.currentGoroutineCount, 1, 0) {
				log.Println("------goroutine count == 0, and close")
				pool.clear()
			} else {
				pool.shutdown()
				atomic.AddInt32(&pool.currentGoroutineCount, -1)
			}
			return
		case worker := <-pool.workerQueue:
			if worker != nil {
				worker.Work()
				pool.waitGroup.Done()
			}
		}
	}
}

/**
 * 内部关闭协程操作
 */
func (pool *GoroutinePool) shutdown() {
	pool.done <- struct{}{}
	log.Println("close one time")
}

/**
 * 关系线程池
 */
func (pool *GoroutinePool) Close() {
	pool.shutdown()
}

/**
 * 收回所有资源
 */
func (pool *GoroutinePool) clear() {
	log.Println("close all channel")
	close(pool.workerQueue)
	close(pool.done)
}
