/**
 * @description goroutinepool
 * @author ibingbo
 * @date 2018/12/1
 */
package woker

// 任务接口
type Workable interface {
	Work()
}

// 任务结构
type Worker struct {
	// 任务动作函数
	Action func()
}

/**
 * Work
 * 任务启动
 */
func (worker *Worker) Work() {
	worker.Action()
}
