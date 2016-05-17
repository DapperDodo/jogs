package jogs

import (
	"fmt"
	"reflect"
)

/*
	Root is the entry function of jogs. It creates the root node and calls the dispatcher,
	effectively starting the process of building the editor for our data

	TODO: maybe this is superfluous. Could clients create the root node and call the dispatcher themselves?
*/
func Root(d *Dispatcher, container string, obj interface{}, cb Callback) {
	node := Node{
		Object:      obj,
		ContainerId: container,
		EditorId:    "root",
		Label:       "",
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

func (n *root) handle(node Node, cb Callback) {

	typekind := reflect.TypeOf(node.Object).Kind()

	switch typekind {
	case reflect.Ptr:
		node.Handle = "PTR"
	default:
		node.Handle = "UNSUPPORTED"
		node.Object = fmt.Sprint("Unsupported root type (root must be PTR):", typekind, "->", node.Object)
	}

	n.dispatch(node, cb)
}
