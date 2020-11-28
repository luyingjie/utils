package timer

import (
	"time"

	vtype "utils/container/type"
)

type Entry struct {
	wheel         *wheel
	job           JobFunc
	singleton     *vtype.Bool
	status        *vtype.Int
	times         *vtype.Int
	create        int64
	interval      int64
	createMs      int64
	intervalMs    int64
	rawIntervalMs int64
}

type JobFunc = func()

func (w *wheel) addEntry(interval time.Duration, job JobFunc, singleton bool, times int, status int) *Entry {
	if times <= 0 {
		times = gDEFAULT_TIMES
	}
	var (
		ms  = interval.Nanoseconds() / 1e6
		num = ms / w.intervalMs
	)
	if num == 0 {
		num = 1
	}
	nowMs := time.Now().UnixNano() / 1e6
	ticks := w.ticks.Val()
	entry := &Entry{
		wheel:         w,
		job:           job,
		times:         vtype.NewInt(times),
		status:        vtype.NewInt(status),
		create:        ticks,
		interval:      num,
		singleton:     vtype.NewBool(singleton),
		createMs:      nowMs,
		intervalMs:    ms,
		rawIntervalMs: ms,
	}

	w.slots[(ticks+num)%w.number].PushBack(entry)
	return entry
}

func (w *wheel) addEntryByParent(interval int64, parent *Entry) *Entry {
	num := interval / w.intervalMs
	if num == 0 {
		num = 1
	}
	nowMs := time.Now().UnixNano() / 1e6
	ticks := w.ticks.Val()
	entry := &Entry{
		wheel:         w,
		job:           parent.job,
		times:         parent.times,
		status:        parent.status,
		create:        ticks,
		interval:      num,
		singleton:     parent.singleton,
		createMs:      nowMs,
		intervalMs:    interval,
		rawIntervalMs: parent.rawIntervalMs,
	}
	w.slots[(ticks+num)%w.number].PushBack(entry)
	return entry
}

func (entry *Entry) Status() int {
	return entry.status.Val()
}

func (entry *Entry) SetStatus(status int) int {
	return entry.status.Set(status)
}

func (entry *Entry) Start() {
	entry.status.Set(STATUS_READY)
}

func (entry *Entry) Stop() {
	entry.status.Set(STATUS_STOPPED)
}

func (entry *Entry) Reset() {
	entry.status.Set(STATUS_RESET)
}

func (entry *Entry) Close() {
	entry.status.Set(STATUS_CLOSED)
}

func (entry *Entry) IsSingleton() bool {
	return entry.singleton.Val()
}

func (entry *Entry) SetSingleton(enabled bool) {
	entry.singleton.Set(enabled)
}

func (entry *Entry) SetTimes(times int) {
	entry.times.Set(times)
}

func (entry *Entry) Run() {
	entry.job()
}

func (entry *Entry) check(nowTicks int64, nowMs int64) (runnable, addable bool) {
	switch entry.status.Val() {
	case STATUS_STOPPED:
		return false, true
	case STATUS_CLOSED:
		return false, false
	case STATUS_RESET:
		return false, true
	}

	if diff := nowTicks - entry.create; diff > 0 && diff%entry.interval == 0 {

		if entry.wheel.level > 0 {
			diffMs := nowMs - entry.createMs
			switch {
			case diffMs < entry.wheel.timer.intervalMs:
				entry.wheel.slots[(nowTicks+entry.interval)%entry.wheel.number].PushBack(entry)
				return false, false
			case diffMs >= entry.wheel.timer.intervalMs:
				if leftMs := entry.intervalMs - diffMs; leftMs > entry.wheel.timer.intervalMs {
					entry.wheel.timer.doAddEntryByParent(leftMs, entry)
					return false, false
				}
			}
		}

		if entry.IsSingleton() {
			if entry.status.Set(STATUS_RUNNING) == STATUS_RUNNING {
				return false, true
			}
		}

		times := entry.times.Add(-1)
		if times <= 0 {
			if entry.status.Set(STATUS_CLOSED) == STATUS_CLOSED || times < 0 {
				return false, false
			}
		}

		if times < 2000000000 && times > 1000000000 {
			entry.times.Set(gDEFAULT_TIMES)
		}
		return true, true
	}
	return false, true
}
