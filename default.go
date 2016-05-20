package jogs

var DefaultDispatcher *Dispatcher

var defaultRoot handler

func init() {

	DefaultDispatcher = &Dispatcher{map[string]handler{}}

	defaultRoot = &root{DefaultDispatcher}

	DefaultDispatcher.Register("ROOT", defaultRoot)
	DefaultDispatcher.Register("LEAF", &leaf{DefaultDispatcher})
	DefaultDispatcher.Register("PTR", &ptr{DefaultDispatcher})
	DefaultDispatcher.Register("PTR_STRUCT", &ptr_struct{DefaultDispatcher})
	DefaultDispatcher.Register("SLICE", &slice{DefaultDispatcher})
	DefaultDispatcher.Register("UNSUPPORTED", newUnsupported())
	DefaultDispatcher.Register("LEAF_INT", newIntHandler())
	DefaultDispatcher.Register("LEAF_STRING", newStringHandler())
}
