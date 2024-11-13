package fsnotify

func (e *Event) String() string {
	return e.event.String()
}

func (e *Event) IsCreate() bool {
	return e.Op == 1 || e.Op&CREATE == CREATE
}

func (e *Event) IsWrite() bool {
	return e.Op&WRITE == WRITE
}

func (e *Event) IsRemove() bool {
	return e.Op&REMOVE == REMOVE
}

func (e *Event) IsRename() bool {
	return e.Op&RENAME == RENAME
}

func (e *Event) IsChmod() bool {
	return e.Op&CHMOD == CHMOD
}
