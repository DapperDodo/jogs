package jogs

import (
	"reflect"

	"github.com/gopherjs/jquery"
)

/*
	Root is the entry function of jogs. It creates the root node and calls the dispatcher,
	effectively starting the process of building the editor for our data

	TODO: maybe this is superfluous. Could clients create the root node and call the dispatcher themselves?
*/
func Root(d *Dispatcher, container string, obj interface{}, cb Callback) {

	jquery.NewJQuery("#" + container).Empty()

	node := Node{
		Object:      obj,
		ContainerId: container,
		EditorId:    "root",
		Handle:      "ROOT",
	}
	d.dispatch(node, cb)
}

///////////////////////////////////////////////////////////////////////
// private parts
///////////////////////////////////////////////////////////////////////

type root struct {
	*Dispatcher
}

func (n *root) Handle(node Node, cb Callback) {

	typekind := reflect.TypeOf(node.Object).Kind()

	switch typekind {
	case reflect.Ptr:
		node.Handle = "PTR"
	default:
		node.Handle = "UNSUPPORTED"
		node.Object = "Unsupported root type (root must be PTR)"
	}

	n.dispatch(node, cb)
}
