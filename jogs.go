package jogs

/*
	Callback propagates changes to the value object
*/
type Callback func(interface{})

/*
	HandlerFunc is a wrapper that makes any function with signature func(Node, Callback) implement jogs' handler interface
*/
type HandlerFunc func(node Node, cb Callback)

//////////////////////////////////////////////////////////////////////
// private parts
//////////////////////////////////////////////////////////////////////

/*
	handler is the interface that node editors should implement
	TODO: is there any reason to make this interface public?
*/
type handler interface {
	handle(node Node, cb Callback)
}

func (f HandlerFunc) handle(node Node, cb Callback) {
	f(node, cb)
}
