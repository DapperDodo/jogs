package jogs

import (
//"fmt"
)

/*
	Dispatcher decouples handles from handlers. This makes jogs a very flexible and extensible tool.
*/
type Dispatcher struct {
	registry map[string]Handler
}

func (d *Dispatcher) Register(handle string, plugin Handler) {
	d.registry[handle] = plugin
}

//////////////////////////////////////////////////////////////////////
// private parts
//////////////////////////////////////////////////////////////////////

func (d *Dispatcher) dispatch(node Node, cb Callback) {

	// fmt.Println("dispatch handle and arguments :", node.Handle, node.Options)

	if d.registry[node.Handle] == nil {
		// fmt.Println("Dispatcher has no Handler registered for handle", node.Handle)
		return
	}

	d.registry[node.Handle].Handle(node, cb)
}
