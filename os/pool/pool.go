// 协程池解决的主要问题实际上就是goroutine的数量过多导致资源侵占，那要解决这个问题就要限制运行的goroutine数量，合理复用，节省资源。
// 具体就是——goroutine池使用协程池将让多余的请求进入排队状态，等待池中有空闲协程的时候来处理，这样就可以通过控制协程池的大小，来控制内存的消耗，让其在极端状态下，也尽可能的保证服务的可用性。
// goroutine的初衷就是轻量级的线程，大部分的时候不需要pool，但是比如网关服务器和DDOS这类场景的时候，goroutine积压会让new开销大于正常，
// 这时对内存的消耗是极大的，协程越多，GC的负担也越大的，一旦内存耗尽，服务也就基本崩溃了...

package pool

import (
	"fmt"
	"time"
)

// Job : 定义一个任务
type Job func([]interface{})
type taskWork struct {
	Run       Job
	startBool bool
	params    []interface{}
}

var WorkMaxTask int
var WorkTaskPool chan taskWork
var WorkTaskReturn chan []interface{}

//启动任务
func (t *taskWork) start() {
	go func() {
		for {
			select {
			case funcRun := <-WorkTaskPool:
				if funcRun.startBool == true {
					funcRun.Run(funcRun.params)
				} else {
					fmt.Println("task  stop!")
					return
				}
			case <-time.After(time.Millisecond * 1000):
				fmt.Print("time out")
			}
		}
	}()
}

func (t *taskWork) stop() {
	fmt.Println("t stop ")
	t.startBool = false
}
func createTask() taskWork {
	var funcJob Job
	var paramSlice []interface{}
	return taskWork{funcJob, true, paramSlice}
}

//循环启动协程池
func StartPool(maxTask int) {
	WorkMaxTask = maxTask
	WorkTaskPool = make(chan taskWork, maxTask)
	WorkTaskReturn = make(chan []interface{}, maxTask)

	for i := 0; i < maxTask; i++ {
		var t = createTask()
		fmt.Println("start task:", i)
		t.start()
	}
}

//消费任务
func Dispatch(funcJob Job, params ...interface{}) {
	WorkTaskPool <- taskWork{funcJob, true, params}
}

//停止协程池
func StopPool() {
	var funcJob Job
	var paramSlice []interface{}
	for i := 0; i < WorkMaxTask; i++ {
		WorkTaskPool <- taskWork{funcJob, false, paramSlice}
	}
}
