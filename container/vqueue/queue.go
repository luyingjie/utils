package vqueue

import (
	"math"

	"github.com/luyingjie/utils/container/vtype"

	"github.com/luyingjie/utils/container/vlist"
)

type Queue struct {
	limit  int
	list   *vlist.List
	closed *vtype.Bool
	events chan struct{}
	C      chan interface{}
}

const (
	gDEFAULT_QUEUE_SIZE     = 10000
	gDEFAULT_MAX_BATCH_SIZE = 10
)

func New(limit ...int) *Queue {
	q := &Queue{
		closed: vtype.NewBool(),
	}
	if len(limit) > 0 && limit[0] > 0 {
		q.limit = limit[0]
		q.C = make(chan interface{}, limit[0])
	} else {
		q.list = vlist.New(true)
		q.events = make(chan struct{}, math.MaxInt32)
		q.C = make(chan interface{}, gDEFAULT_QUEUE_SIZE)
		go q.asyncLoopFromListToChannel()
	}
	return q
}

func (q *Queue) asyncLoopFromListToChannel() {
	defer func() {
		if q.closed.Val() {
			_ = recover()
		}
	}()
	for !q.closed.Val() {
		<-q.events
		for !q.closed.Val() {
			if length := q.list.Len(); length > 0 {
				if length > gDEFAULT_MAX_BATCH_SIZE {
					length = gDEFAULT_MAX_BATCH_SIZE
				}
				for _, v := range q.list.PopFronts(length) {
					q.C <- v
				}
			} else {
				break
			}
		}

		for i := 0; i < len(q.events)-1; i++ {
			<-q.events
		}
	}

	close(q.C)
}

func (q *Queue) Push(v interface{}) {
	if q.limit > 0 {
		q.C <- v
	} else {
		q.list.PushBack(v)
		if len(q.events) < gDEFAULT_QUEUE_SIZE {
			q.events <- struct{}{}
		}
	}
}

func (q *Queue) Pop() interface{} {
	return <-q.C
}

func (q *Queue) Close() {
	q.closed.Set(true)
	if q.events != nil {
		close(q.events)
	}
	if q.limit > 0 {
		close(q.C)
	}
	for i := 0; i < gDEFAULT_MAX_BATCH_SIZE; i++ {
		q.Pop()
	}
}

func (q *Queue) Len() (length int) {
	if q.list != nil {
		length += q.list.Len()
	}
	length += len(q.C)
	return
}

func (q *Queue) Size() int {
	return q.Len()
}
