package fsnotify

import (
	"fmt"

	"github.com/luyingjie/utils/container/vlist"
)

func (w *Watcher) startWatchLoop() {
	go func() {
		for {
			select {
			case <-w.closeChan:
				return

			case ev := <-w.watcher.Events:
				w.cache.SetIfNotExist(ev.String(), func() interface{} {
					w.events.Push(&Event{
						event:   ev,
						Path:    ev.Name,
						Op:      Op(ev.Op),
						Watcher: w,
					})
					return struct{}{}
				}, repeatEventFilterDuration)

			case err := <-w.watcher.Errors:
				fmt.Println(err)
			}
		}
	}()
}

func (w *Watcher) getCallbacks(path string) (callbacks []*Callback) {
	if v := w.callbacks.Get(path); v != nil {
		for _, v := range v.(*vlist.List).FrontAll() {
			callback := v.(*Callback)
			callbacks = append(callbacks, callback)
		}
	}
	dirPath := fileDir(path)
	if v := w.callbacks.Get(dirPath); v != nil {
		for _, v := range v.(*vlist.List).FrontAll() {
			callback := v.(*Callback)
			if callback.recursive {
				callbacks = append(callbacks, callback)
			}
		}
	}
	for {
		parentDirPath := fileDir(dirPath)
		if parentDirPath == dirPath {
			break
		}
		if v := w.callbacks.Get(parentDirPath); v != nil {
			for _, v := range v.(*vlist.List).FrontAll() {
				callback := v.(*Callback)
				if callback.recursive {
					callbacks = append(callbacks, callback)
				}
			}
		}
		dirPath = parentDirPath
	}
	return
}

func (w *Watcher) startEventLoop() {
	go func() {
		for {
			if v := w.events.Pop(); v != nil {
				event := v.(*Event)
				callbacks := w.getCallbacks(event.Path)
				if len(callbacks) == 0 {
					w.watcher.Remove(event.Path)
					continue
				}
				switch {
				case event.IsRemove():
					if fileExists(event.Path) {
						if err := w.watcher.Add(event.Path); err != nil {
							fmt.Println(err)
						} else {
							// intlog.Printf("fake remove event, watcher re-adds monitor for: %s", event.Path)
						}
						event.Op = RENAME
					}
				case event.IsRename():
					if fileExists(event.Path) {
						if err := w.watcher.Add(event.Path); err != nil {
							fmt.Println(err)
						} else {
							// intlog.Printf("fake rename event, watcher re-adds monitor for: %s", event.Path)
						}
						event.Op = CHMOD
					}
				case event.IsCreate():
					if fileIsDir(event.Path) {
						for _, subPath := range fileAllDirs(event.Path) {
							if fileIsDir(subPath) {
								if err := w.watcher.Add(subPath); err != nil {
									fmt.Println(err)
								} else {
									// intlog.Printf("folder creation event, watcher adds monitor for: %s", subPath)
								}
							}
						}
					} else {
						if err := w.watcher.Add(event.Path); err != nil {
							fmt.Println(err)
						} else {
							// intlog.Printf("file creation event, watcher adds monitor for: %s", event.Path)
						}
					}

				}
				for _, v := range callbacks {
					go func(callback *Callback) {
						defer func() {
							if err := recover(); err != nil {
								switch err {
								case callbackExitEventPanicStr:
									w.RemoveCallback(callback.Id)
								default:
									panic(err)
								}
							}
						}()
						callback.Func(event)
					}(v)
				}
			} else {
				break
			}
		}
	}()
}
