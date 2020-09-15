package timer

import (
	"time"

	"utils/container/list"
)

func (w *wheel) start() {
	go func() {
		ticker := time.NewTicker(time.Duration(w.intervalMs) * time.Millisecond)
		for {
			select {
			case <-ticker.C:
				switch w.timer.status.Val() {
				case STATUS_RUNNING:
					w.proceed()

				case STATUS_STOPPED:
				case STATUS_CLOSED:
					ticker.Stop()
					return
				}

			}
		}
	}()
}

func (w *wheel) proceed() {
	n := w.ticks.Add(1)
	l := w.slots[int(n%w.number)]
	length := l.Len()
	if length > 0 {
		go func(l *list.List, nowTicks int64) {
			entry := (*Entry)(nil)
			nowMs := time.Now().UnixNano() / 1e6
			for i := length; i > 0; i-- {
				if v := l.PopFront(); v == nil {
					break
				} else {
					entry = v.(*Entry)
				}

				runnable, addable := entry.check(nowTicks, nowMs)
				if runnable {
					go func(entry *Entry) {
						defer func() {
							if err := recover(); err != nil {
								if err != gPANIC_EXIT {
									panic(err)
								} else {
									entry.Close()
								}
							}
							if entry.Status() == STATUS_RUNNING {
								entry.SetStatus(STATUS_READY)
							}
						}()
						entry.job()
					}(entry)
				}

				if addable {

					if entry.Status() == STATUS_RESET {
						entry.SetStatus(STATUS_READY)
					}
					entry.wheel.timer.doAddEntryByParent(entry.rawIntervalMs, entry)
				}
			}
		}(l, n)
	}
}
