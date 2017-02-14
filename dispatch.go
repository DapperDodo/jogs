package jogs

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

	if d.registry[node.Handle] == nil {
		// TODO: log this
		return
	}

	d.registry[node.Handle].Handle(node, cb)
}
