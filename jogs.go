package jogs

import (
	"github.com/gopherjs/jquery"
)

/*
	Convenience handle to jQuery
*/
var J = jquery.NewJQuery

/*
	Callback propagates changes to the value object
*/
type Callback func(interface{})

/*
	handler is the interface that node editors should implement
*/
type Handler interface {
	Handle(node Node, cb Callback)
}

/*
	HandlerFunc is a wrapper that makes any function with signature func(Node, Callback) implement jogs' handler interface
*/
type HandlerFunc func(node Node, cb Callback)

func (f HandlerFunc) Handle(node Node, cb Callback) {
	f(node, cb)
}
