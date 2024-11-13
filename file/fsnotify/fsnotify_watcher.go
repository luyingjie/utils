package fsnotify

import (
	"errors"
	"fmt"

	mylist "github.com/luyingjie/utils/container/vlist"
)

func (w *Watcher) Add(path string, callbackFunc func(event *Event), recursive ...bool) (callback *Callback, err error) {
	return w.AddOnce("", path, callbackFunc, recursive...)
}

func (w *Watcher) AddOnce(name, path string, callbackFunc func(event *Event), recursive ...bool) (callback *Callback, err error) {
	w.nameSet.AddIfNotExistFuncLock(name, func() bool {
		callback, err = w.addWithCallbackFunc(name, path, callbackFunc, recursive...)
		if err != nil {
			return false
		}
		if fileIsDir(path) && (len(recursive) == 0 || recursive[0]) {
			for _, subPath := range fileAllDirs(path) {
				if fileIsDir(subPath) {
					if err := w.watcher.Add(subPath); err != nil {
						fmt.Println(err)
					} else {
						// intlog.Printf("watcher adds monitor for: %s", subPath)
					}
				}
			}
		}
		if name == "" {
			return false
		}
		return true
	})
	return
}

func (w *Watcher) addWithCallbackFunc(name, path string, callbackFunc func(event *Event), recursive ...bool) (callback *Callback, err error) {
	if t := fileRealPath(path); t == "" {
		return nil, errors.New(fmt.Sprintf(`"%s" does not exist`, path))
	} else {
		path = t
	}

	callback = &Callback{
		Id:        callbackIdGenerator.Add(1),
		Func:      callbackFunc,
		Path:      path,
		name:      name,
		recursive: true,
	}
	if len(recursive) > 0 {
		callback.recursive = recursive[0]
	}

	w.callbacks.LockFunc(func(m map[string]interface{}) {
		list := (*mylist.List)(nil)
		if v, ok := m[path]; !ok {
			list = mylist.New(true)
			m[path] = list
		} else {
			list = v.(*mylist.List)
		}
		callback.elem = list.PushBack(callback)
	})

	if err := w.watcher.Add(path); err != nil {
		fmt.Println(err)
	} else {
		// intlog.Printf("watcher adds monitor for: %s", path)
	}

	callbackIdMap.Set(callback.Id, callback)

	return
}

func (w *Watcher) Close() {
	w.events.Close()
	if err := w.watcher.Close(); err != nil {
		fmt.Println(err)
	}
	close(w.closeChan)
}

func (w *Watcher) Remove(path string) error {
	if r := w.callbacks.Remove(path); r != nil {
		list := r.(*mylist.List)
		for {
			if r := list.PopFront(); r != nil {
				callbackIdMap.Remove(r.(*Callback).Id)
			} else {
				break
			}
		}
	}

	if subPaths, err := fileScanDir(path, "*", true); err == nil && len(subPaths) > 0 {
		for _, subPath := range subPaths {
			if w.checkPathCanBeRemoved(subPath) {
				if err := w.watcher.Remove(subPath); err != nil {
					fmt.Println(err)
				}
			}
		}
	}

	return w.watcher.Remove(path)
}

func (w *Watcher) checkPathCanBeRemoved(path string) bool {
	if v := w.callbacks.Get(path); v != nil {
		return false
	}

	dirPath := fileDir(path)
	if v := w.callbacks.Get(dirPath); v != nil {
		for _, c := range v.(*mylist.List).FrontAll() {
			if c.(*Callback).recursive {
				return false
			}
		}
		return false
	}

	parentDirPath := ""
	for {
		parentDirPath = fileDir(dirPath)
		if parentDirPath == dirPath {
			break
		}
		if v := w.callbacks.Get(parentDirPath); v != nil {
			for _, c := range v.(*mylist.List).FrontAll() {
				if c.(*Callback).recursive {
					return false
				}
			}
			return false
		}
		dirPath = parentDirPath
	}
	return true
}

func (w *Watcher) RemoveCallback(callbackId int) {
	callback := (*Callback)(nil)
	if r := callbackIdMap.Get(callbackId); r != nil {
		callback = r.(*Callback)
	}
	if callback != nil {
		if r := w.callbacks.Get(callback.Path); r != nil {
			r.(*mylist.List).Remove(callback.elem)
		}
		callbackIdMap.Remove(callbackId)
		if callback.name != "" {
			w.nameSet.Remove(callback.name)
		}
	}
}
