package timer

import (
	"fmt"
	"time"

	"github.com/luyingjie/utils/container/vlist"
	"github.com/luyingjie/utils/container/vtype"
)

type Timer struct {
	status     *vtype.Int
	wheels     []*wheel
	length     int
	number     int
	intervalMs int64
}

type wheel struct {
	timer      *Timer
	level      int
	slots      []*vlist.List
	number     int64
	ticks      *vtype.Int64
	totalMs    int64
	createMs   int64
	intervalMs int64
}

func New(slot int, interval time.Duration, level ...int) *Timer {
	if slot <= 0 {
		panic(fmt.Sprintf("invalid slot number: %d", slot))
	}
	length := gDEFAULT_WHEEL_LEVEL
	if len(level) > 0 {
		length = level[0]
	}
	t := &Timer{
		status:     vtype.NewInt(STATUS_RUNNING),
		wheels:     make([]*wheel, length),
		length:     length,
		number:     slot,
		intervalMs: interval.Nanoseconds() / 1e6,
	}
	for i := 0; i < length; i++ {
		if i > 0 {
			n := time.Duration(t.wheels[i-1].totalMs) * time.Millisecond
			if n <= 0 {
				panic(fmt.Sprintf(`inteval is too large with level: %dms x %d`, interval, length))
			}
			w := t.newWheel(i, slot, n)
			t.wheels[i] = w
			t.wheels[i-1].addEntry(n, w.proceed, false, gDEFAULT_TIMES, STATUS_READY)
		} else {
			t.wheels[i] = t.newWheel(i, slot, interval)
		}
	}
	t.wheels[0].start()
	return t
}

func (t *Timer) newWheel(level int, slot int, interval time.Duration) *wheel {
	w := &wheel{
		timer:      t,
		level:      level,
		slots:      make([]*vlist.List, slot),
		number:     int64(slot),
		ticks:      vtype.NewInt64(),
		totalMs:    int64(slot) * interval.Nanoseconds() / 1e6,
		createMs:   time.Now().UnixNano() / 1e6,
		intervalMs: interval.Nanoseconds() / 1e6,
	}
	for i := int64(0); i < w.number; i++ {
		w.slots[i] = vlist.New(true)
	}
	return w
}

func (t *Timer) Add(interval time.Duration, job JobFunc) *Entry {
	return t.doAddEntry(interval, job, false, gDEFAULT_TIMES, STATUS_READY)
}

func (t *Timer) AddEntry(interval time.Duration, job JobFunc, singleton bool, times int, status int) *Entry {
	return t.doAddEntry(interval, job, singleton, times, status)
}

func (t *Timer) AddSingleton(interval time.Duration, job JobFunc) *Entry {
	return t.doAddEntry(interval, job, true, gDEFAULT_TIMES, STATUS_READY)
}

func (t *Timer) AddOnce(interval time.Duration, job JobFunc) *Entry {
	return t.doAddEntry(interval, job, true, 1, STATUS_READY)
}

func (t *Timer) AddTimes(interval time.Duration, times int, job JobFunc) *Entry {
	return t.doAddEntry(interval, job, true, times, STATUS_READY)
}

func (t *Timer) DelayAdd(delay time.Duration, interval time.Duration, job JobFunc) {
	t.AddOnce(delay, func() {
		t.Add(interval, job)
	})
}

func (t *Timer) DelayAddEntry(delay time.Duration, interval time.Duration, job JobFunc, singleton bool, times int, status int) {
	t.AddOnce(delay, func() {
		t.AddEntry(interval, job, singleton, times, status)
	})
}

func (t *Timer) DelayAddSingleton(delay time.Duration, interval time.Duration, job JobFunc) {
	t.AddOnce(delay, func() {
		t.AddSingleton(interval, job)
	})
}

func (t *Timer) DelayAddOnce(delay time.Duration, interval time.Duration, job JobFunc) {
	t.AddOnce(delay, func() {
		t.AddOnce(interval, job)
	})
}

func (t *Timer) DelayAddTimes(delay time.Duration, interval time.Duration, times int, job JobFunc) {
	t.AddOnce(delay, func() {
		t.AddTimes(interval, times, job)
	})
}

func (t *Timer) Start() {
	t.status.Set(STATUS_RUNNING)
}

func (t *Timer) Stop() {
	t.status.Set(STATUS_STOPPED)
}

func (t *Timer) Close() {
	t.status.Set(STATUS_CLOSED)
}

func (t *Timer) doAddEntry(interval time.Duration, job JobFunc, singleton bool, times int, status int) *Entry {
	return t.wheels[t.getLevelByIntervalMs(interval.Nanoseconds()/1e6)].addEntry(interval, job, singleton, times, status)
}

func (t *Timer) doAddEntryByParent(interval int64, parent *Entry) *Entry {
	return t.wheels[t.getLevelByIntervalMs(interval)].addEntryByParent(interval, parent)
}

func (t *Timer) getLevelByIntervalMs(intervalMs int64) int {
	pos, cmp := t.binSearchIndex(intervalMs)
	switch cmp {
	case 0:
		fallthrough
	case -1:
		i := pos
		for ; i > 0; i-- {
			if intervalMs > t.wheels[i].intervalMs && intervalMs <= t.wheels[i].totalMs {
				return i
			}
		}
		return i
	case 1:
		i := pos
		for ; i < t.length-1; i++ {
			if intervalMs > t.wheels[i].intervalMs && intervalMs <= t.wheels[i].totalMs {
				return i
			}
		}
		return i
	}
	return 0
}

func (t *Timer) binSearchIndex(n int64) (index int, result int) {
	min := 0
	max := t.length - 1
	mid := 0
	cmp := -2
	for min <= max {
		mid = int((min + max) / 2)
		switch {
		case t.wheels[mid].intervalMs == n:
			cmp = 0
		case t.wheels[mid].intervalMs > n:
			cmp = -1
		case t.wheels[mid].intervalMs < n:
			cmp = 1
		}
		switch cmp {
		case -1:
			max = mid - 1
		case 1:
			min = mid + 1
		case 0:
			return mid, cmp
		}
	}
	return mid, cmp
}
