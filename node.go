package jogs

import (
	"fmt"
	"html/template"
	"reflect"
	"strings"
)

/*
	Node is a structure containing all necessary info for building an editor node.
*/
type Node struct {
	Object      interface{} // the actual data object that we want to edit
	ContainerId string      // the id of the dom node in which we want this editor placed, a.k.a. parent
	EditorId    string      // the id of the dom node of this editor
	Handle      string      // a key to the handler that will build this node's editor
	Label       string      // the name of this editor, also the form field name
	Placeholder string      // text to display when the field is empty
	Idx         int         // the index of the node in the struct field list or slice
	Tags        []string    // (optional) arguments for the handler, filled by parsing a struct tag called 'jogs' with comma separated strings (see OVERRIDE_ constants)
}

const (
	OVERRIDE_HANDLER     int = iota // tag field 0 is for custom handlers
	OVERRIDE_LABEL                  // tag field 1 is for overriding the field label (name)
	OVERRIDE_PLACEHOLDER            // tag field 2 is for overriding the placeholder (text to display when the field is empty)
)

//////////////////////////////////////////////////////////////////////////////////////
// private parts
//////////////////////////////////////////////////////////////////////////////////////

/*
	ptr is a handler for nodes of type pointer
*/
type ptr struct {
	*Dispatcher
}

func (n *ptr) Handle(node Node, cb Callback) {

	valueofkind := reflect.ValueOf(node.Object).Elem().Kind()

	switch valueofkind {
	case reflect.Struct:
		node.Handle = "PTR_STRUCT"
	default:
		node.Handle = "UNSUPPORTED"
		node.Object = fmt.Sprint("Unsupported ptr to <T> (must be ptr to struct)  :", valueofkind, "->", node.Object)
	}

	n.dispatch(node, cb)
}

/*
	ptr_struct is a handler for nodes of type pointer to struct
*/
type ptr_struct struct {
	*Dispatcher
}

func (n *ptr_struct) field_sep(c rune) bool {
	return c == ','
}

func (n *ptr_struct) Handle(node Node, cb Callback) {

	e := reflect.ValueOf(node.Object).Elem()

	for i := 0; i < e.NumField(); i++ {
		field_value := e.Field(i)
		field_name := e.Type().Field(i).Name
		//fmt.Println("struct field", i, field_name, ":", field_value.Kind(), "->", field_value.Interface())
		if !field_value.CanSet() {
			fmt.Println("field is not settable")
			continue
		}

		node_row := node
		node_row.Idx = i
		node_row.Label = field_name

		// parse optional tags
		tag := e.Type().Field(i).Tag.Get("jogs")
		if tag != "" {
			node_row.Tags = strings.FieldsFunc(tag, n.field_sep)
		}

		// case 2: standard nodes
		switch e.Field(i).Type().Kind() {
		case reflect.Struct:
			node_nested := n.nest(node_row)
			node_nested.Object = field_value.Addr().Interface()
			node_nested.Handle = "PTR"
			n.dispatch(node_nested, func(out interface{}) {
				cb(node.Object)
			})
		case reflect.Ptr:
			node_nested := n.nest(node_row)
			node_nested.Object = field_value.Interface()
			node_nested.Handle = "PTR"
			n.dispatch(node_nested, func(out interface{}) {
				cb(node.Object)
			})
		case reflect.Slice:
			node_nested := n.nest(node_row)
			node_nested.Object = field_value.Interface()
			node_nested.Handle = "SLICE"
			n.dispatch(node_nested, func(out interface{}) {
				field_value.Set(reflect.ValueOf(out))
				cb(node.Object)
			})
		default:
			node_row.EditorId += "-" + node_row.Label
			node_row.Object = field_value.Interface()
			node_row.Handle = "LEAF"
			n.dispatch(node_row, func(out interface{}) {
				field_value.Set(reflect.ValueOf(out))
				cb(node.Object)
			})
		}
	}
}

/////////////////////////////////////////////////////////////////////////////

var nest_tpl = template.Must(template.New("skin").Parse(string(`
	{{define "nest"}}
		<div class="row" id="{{.EditorId}}-slice">
			<div class="col-lg-1" id="{{.EditorId}}-margin">
				<label class="control-label">{{.Label}}</label>
			</div>
			<div class="col-lg-11" id="{{.EditorId}}-content">
			</div>
		</div>
	{{end}}
`)))

func (n *ptr_struct) nest(node Node) Node {
	child := node
	// intercept and clear Label override
	if len(child.Tags) > OVERRIDE_LABEL && child.Tags[OVERRIDE_LABEL] != "" {
		child.Label = child.Tags[OVERRIDE_LABEL]
		child.Tags[OVERRIDE_LABEL] = ""
	}
	child.EditorId += "-" + child.Label
	J("#" + child.ContainerId).Append(Merge(nest_tpl, "nest", child))
	child.Label = ""
	child.ContainerId = child.EditorId
	child.EditorId += "-content"
	child.ContainerId += "-content"
	return child
}
