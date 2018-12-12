/**
 * @description goroutinepool
 * @author ibingbo
 * @date 2018/12/1
 */
package goroutinepool

import (
	"goroutinepool/woker"
	"log"
	"sync"
	"sync/atomic"
)

// 默认最大任务数
var defaultMaxWorkerCount = 100

// 默认最大协程数
var defaultMaxGoroutineCount = 10

// 协程池结构定义
type GoroutinePool struct {
	maxGoroutineCount  int            // 最大协程数
	maxWorkerCount     int            // 最大任务数
	workerQueue        chan woker.Workable  // 任务队列
	done               chan struct{}  // 任务状态，是否所有任务已完成
	waitGroup          sync.WaitGroup // 协程等待控制
	currentWorkerCount int32          // 当前要执行的任务数
	once               sync.Once
	stop               func()
}

/**
 * NewGoroutinePool
 * 实例化一个协程池
 */
func NewGoroutinePool(maxGoroutineCount int, maxWorkerCount int) *GoroutinePool {
	if maxGoroutineCount == 0 {
		maxGoroutineCount = defaultMaxGoroutineCount
	}
	if maxWorkerCount == 0 {
		maxWorkerCount = defaultMaxWorkerCount
	}
	pool := &GoroutinePool{
		maxGoroutineCount:  maxGoroutineCount,
		workerQueue:        make(chan woker.Workable, maxWorkerCount),
		done:               make(chan struct{}),
		currentWorkerCount: 0,
	}
	return pool
}

/**
 * 提交并执行一个任务
 */
func (pool *GoroutinePool) Execute(worker woker.Workable) {
	pool.workerQueue <- worker
	atomic.AddInt32(&pool.currentWorkerCount, 1)
	log.Println("add one worker")
	pool.once.Do(pool.start)
}

/**
 * 显式调用手动关闭协程池
 */
func (pool *GoroutinePool) Shutdown() {
	pool.shutdown()
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
	pool.waitGroup.Add(pool.maxGoroutineCount)
	for i := 0; i < pool.maxGoroutineCount; i++ {
		go pool.doWork()
	}
	go pool.watch()
}

/**
 * 协程池监控，关闭协程
 */
func (pool *GoroutinePool) watch() {
	for {
		select {
		case <-pool.done:
			pool.shutdown()
			return
		default:
			if atomic.LoadInt32(&pool.currentWorkerCount) == 0 {
				log.Println("all worker done....")
				pool.shutdown()
			}
		}
	}
}

/**
 * 开始处理所有的任务
 */
func (pool *GoroutinePool) doWork() {
	for {
		select {
		case <-pool.done:
			log.Println("close goroutine....")
			pool.waitGroup.Done()
			pool.shutdown()
			return
		case worker := <-pool.workerQueue:
			if worker != nil {
				worker.Work()
				atomic.AddInt32(&pool.currentWorkerCount, -1)
			}
		}
	}
}

/**
 * 内部关闭协程操作
 */
func (pool *GoroutinePool) shutdown() {
	pool.done <- struct{}{}
	pool.once.Do(pool.close)
}

func (pool *GoroutinePool) close() {
	close(pool.workerQueue)
}

// 并发访问MAP
type ConcurrentMap struct {
	sync.RWMutex
	Map map[string]interface{}
}

// 实例化一个线程安全MAP
func NewConcurrentMap() *ConcurrentMap {
	sm := &ConcurrentMap{}
	sm.Map = map[string]interface{}{}
	return sm
}

// 线程安全获取一个元素
func (sm *ConcurrentMap) Get(key string) interface{} {
	sm.RLock()
	defer sm.RUnlock()
	return sm.Map[key]
}

// 线程安全添加一个元素
func (sm *ConcurrentMap) Put(key string, value interface{}) {
	sm.Lock()
	defer sm.Unlock()
	sm.Map[key] = value
}

// 线程安全删除一个元素
func (sm *ConcurrentMap) Delete(key string) {
	sm.Lock()
	defer sm.Unlock()
	delete(sm.Map, key)
}
