package timer

import (
	"fmt"
	"math"
	"time"

	"utils/util/cmdenv"
)

const (
	STATUS_READY            = 0
	STATUS_RUNNING          = 1
	STATUS_STOPPED          = 2
	STATUS_RESET            = 3
	STATUS_CLOSED           = -1
	gPANIC_EXIT             = "exit"
	gDEFAULT_TIMES          = math.MaxInt32
	gDEFAULT_SLOT_NUMBER    = 10
	gDEFAULT_WHEEL_INTERVAL = 50
	gDEFAULT_WHEEL_LEVEL    = 6
	gCMDENV_KEY             = "gf.gtimer"
)

var (
	defaultSlots    = cmdenv.Get(fmt.Sprintf("%s.slots", gCMDENV_KEY), gDEFAULT_SLOT_NUMBER).Int()
	defaultLevel    = cmdenv.Get(fmt.Sprintf("%s.level", gCMDENV_KEY), gDEFAULT_WHEEL_LEVEL).Int()
	defaultInterval = cmdenv.Get(fmt.Sprintf("%s.interval", gCMDENV_KEY), gDEFAULT_WHEEL_INTERVAL).Duration() * time.Millisecond
	defaultTimer    = New(defaultSlots, defaultInterval, defaultLevel)
)

func SetTimeout(delay time.Duration, job JobFunc) {
	AddOnce(delay, job)
}

func SetInterval(interval time.Duration, job JobFunc) {
	Add(interval, job)
}

func Add(interval time.Duration, job JobFunc) *Entry {
	return defaultTimer.Add(interval, job)
}

func AddEntry(interval time.Duration, job JobFunc, singleton bool, times int, status int) *Entry {
	return defaultTimer.AddEntry(interval, job, singleton, times, status)
}

func AddSingleton(interval time.Duration, job JobFunc) *Entry {
	return defaultTimer.AddSingleton(interval, job)
}

func AddOnce(interval time.Duration, job JobFunc) *Entry {
	return defaultTimer.AddOnce(interval, job)
}

func AddTimes(interval time.Duration, times int, job JobFunc) *Entry {
	return defaultTimer.AddTimes(interval, times, job)
}

func DelayAdd(delay time.Duration, interval time.Duration, job JobFunc) {
	defaultTimer.DelayAdd(delay, interval, job)
}

func DelayAddEntry(delay time.Duration, interval time.Duration, job JobFunc, singleton bool, times int, status int) {
	defaultTimer.DelayAddEntry(delay, interval, job, singleton, times, status)
}

func DelayAddSingleton(delay time.Duration, interval time.Duration, job JobFunc) {
	defaultTimer.DelayAddSingleton(delay, interval, job)
}

func DelayAddOnce(delay time.Duration, interval time.Duration, job JobFunc) {
	defaultTimer.DelayAddOnce(delay, interval, job)
}

func DelayAddTimes(delay time.Duration, interval time.Duration, times int, job JobFunc) {
	defaultTimer.DelayAddTimes(delay, interval, times, job)
}

func Exit() {
	panic(gPANIC_EXIT)
}
