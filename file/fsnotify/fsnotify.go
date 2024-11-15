package fsnotify

import (
	"container/list"
	"errors"
	"fmt"
	"time"

	"github.com/luyingjie/utils/container/vset"

	"github.com/luyingjie/utils/storage/cache"

	"github.com/luyingjie/utils/container/vmap"

	"github.com/luyingjie/utils/container/vqueue"

	"github.com/luyingjie/utils/container/vtype"

	gitFsnotify "github.com/fsnotify/fsnotify"
)

type Watcher struct {
	watcher   *gitFsnotify.Watcher
	events    *vqueue.Queue
	cache     *cache.Cache
	nameSet   *vset.StrSet
	callbacks *vmap.StrAnyMap
	closeChan chan struct{}
}

type Callback struct {
	Id        int
	Func      func(event *Event)
	Path      string
	name      string
	elem      *list.Element
	recursive bool
}

type Event struct {
	event   gitFsnotify.Event
	Path    string
	Op      Op
	Watcher *Watcher
}

type Op uint32

const (
	CREATE Op = 1 << iota
	WRITE
	REMOVE
	RENAME
	CHMOD
)

const (
	repeatEventFilterDuration = time.Millisecond
	callbackExitEventPanicStr = "exit"
)

var (
	defaultWatcher      *Watcher
	callbackIdMap       = vmap.NewIntAnyMap(true)
	callbackIdGenerator = vtype.NewInt()
)

func init() {
	var err error
	defaultWatcher, err = New()
	if err != nil {
		panic(fmt.Sprintf(`creating default fsnotify watcher failed: %s`, err.Error()))
	}
}

func New() (*Watcher, error) {
	w := &Watcher{
		cache:     cache.New(),
		events:    vqueue.New(),
		nameSet:   vset.NewStrSet(true),
		closeChan: make(chan struct{}),
		callbacks: vmap.NewStrAnyMap(true),
	}
	if watcher, err := gitFsnotify.NewWatcher(); err == nil {
		w.watcher = watcher
	} else {
		return nil, err
	}
	w.startWatchLoop()
	w.startEventLoop()
	return w, nil
}

func Add(path string, callbackFunc func(event *Event), recursive ...bool) (callback *Callback, err error) {
	return defaultWatcher.Add(path, callbackFunc, recursive...)
}

func AddOnce(name, path string, callbackFunc func(event *Event), recursive ...bool) (callback *Callback, err error) {
	return defaultWatcher.AddOnce(name, path, callbackFunc, recursive...)
}

func Remove(path string) error {
	return defaultWatcher.Remove(path)
}

func RemoveCallback(callbackId int) error {
	callback := (*Callback)(nil)
	if r := callbackIdMap.Get(callbackId); r != nil {
		callback = r.(*Callback)
	}
	if callback == nil {
		return errors.New(fmt.Sprintf(`callback for id %d not found`, callbackId))
	}
	defaultWatcher.RemoveCallback(callbackId)
	return nil
}

func Exit() {
	panic(callbackExitEventPanicStr)
}
