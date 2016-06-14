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
	DefaultDispatcher.Register("SLICE", newSlice(DefaultDispatcher))
	DefaultDispatcher.Register("UNSUPPORTED", newUnsupported())
	DefaultDispatcher.Register("LEAF_INT", newIntHandler())
	DefaultDispatcher.Register("LEAF_STRING", newStringHandler())
	DefaultDispatcher.Register("LEAF_FLOAT", newFloatHandler())
	DefaultDispatcher.Register("LEAF_BOOL", newBoolHandler())
}
