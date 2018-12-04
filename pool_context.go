/**
 * @description goroutinepool
 * @author ibingbo
 * @date 2018/12/1
 */
package goroutinepool

import (
	"context"
	"log"
	"sync/atomic"
)

/**
 * 提交并执行一个任务
 */
func (pool *GoroutinePool) ExecuteWithContext(worker Workable, ctx context.Context) {
	select {
	case <-ctx.Done():
		log.Println("goroutine had canceled...")
	default:
		pool.workerQueue <- worker
		atomic.AddInt32(&pool.currentWorkerCount, 1)
		log.Println("add one worker")
		pool.once.Do(func() {
			pool.startWithContext(ctx)
		})
	}

}

/**
 * 协程池开始执行
 */
func (pool *GoroutinePool) startWithContext(ctx context.Context) {
	pool.waitGroup.Add(pool.maxGoroutineCount)
	for i := 0; i < pool.maxGoroutineCount; i++ {
		go pool.doWorkWithContext(ctx)
	}
	go pool.watchWithContext(ctx)
}

/**
 * 协程池监控，关闭协程
 */
func (pool *GoroutinePool) watchWithContext(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			pool.shutdown()
			return
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
func (pool *GoroutinePool) doWorkWithContext(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("cancel goroutine....")
			pool.waitGroup.Done()
			pool.shutdown()
			return
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
