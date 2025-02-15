package main

type Request interface {
	Object() Object // Object that triggered it, if any.
}

type RequestLevel struct {
	object    Object
	Level     string
	FromCellX int
	FromCellY int
}

func (r RequestLevel) Object() Object {
	return r.object
}

type RequestReset struct {
	object Object
}

func (r RequestReset) Object() Object {
	return r.object
}

type RequestDelete struct {
	object Object
}

func (r RequestDelete) Object() Object {
	return r.object
}

type RequestAdd struct {
	object Object
}

func (r RequestAdd) Object() Object {
	return r.object
}

type RequestSetCell struct {
	object Object
	x      int
	y      int
	flag   CellFlag
	image  int
}

func (r RequestSetCell) Object() Object {
	return r.object
}
