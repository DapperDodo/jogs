package jogs

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/gopherjs/jquery"
)

type slice struct {
	*Dispatcher
}

func (s *slice) handle(node Node, cb Callback) {

	typ := reflect.TypeOf(node.Object)
	val := reflect.ValueOf(node.Object)

	jQuery("#" + node.EditorId).Empty()

	e := typ.Elem()

	for i := 0; i < val.Len(); i++ {

		valrow := val.Index(i)
		//fmt.Println("field", i, ":", valrow.Kind(), "->", valrow.Interface())

		jQuery("#" + node.EditorId).Append("<div class=\"row\" id=\"" + node.EditorId + "-" + strconv.Itoa(i) + "-row\"></div>")
		jQuery("#" + node.EditorId + "-" + strconv.Itoa(i) + "-row").Append("<div class=\"col-lg-1\" id=\"" + node.EditorId + "-" + strconv.Itoa(i) + "-col-L\"></div>")
		jQuery("#" + node.EditorId + "-" + strconv.Itoa(i) + "-row").Append("<div class=\"col-lg-11\" id=\"" + node.EditorId + "-" + strconv.Itoa(i) + "-col-R\"></div>")
		jQuery("#" + node.EditorId + "-" + strconv.Itoa(i) + "-col-L").Append("<button type=\"button\" id=\"" + node.EditorId + "-" + strconv.Itoa(i) + "-del\" class=\"btn btn-danger btn-xs\"><i class=\"fa fa-trash\"></i></button>")

		noderow := node
		noderow.EditorId += "-" + strconv.Itoa(i)
		noderow.Label = ""
		noderow.ContainerId = noderow.EditorId + "-col-R"

		switch valrow.Kind() {
		case reflect.String, reflect.Int:
			noderow.Handle = "LEAF"
			noderow.Object = valrow.Interface()
			s.dispatch(noderow, func(out interface{}) {
				valrow.Set(reflect.ValueOf(out))
				cb(val.Interface())
			})
		case reflect.Struct:
			noderow.Handle = "PTR"
			noderow.Object = valrow.Addr().Interface()
			s.dispatch(noderow, func(out interface{}) {
				cb(val.Interface())
			})
		case reflect.Ptr:
			noderow.Handle = "PTR"
			noderow.Object = valrow.Interface()
			s.dispatch(noderow, func(out interface{}) {
				cb(val.Interface())
			})
		default:
			jQuery("#" + node.EditorId).Empty()
			noderow.Handle = "UNSUPPORTED"
			noderow.Object = fmt.Sprint("unsupported slice type :", e.Kind(), "->", valrow.Interface())
			node.Label = ""
			s.dispatch(noderow, nil)
			return
		}

		i_closed_over := i
		jQuery("#"+node.EditorId+"-"+strconv.Itoa(i)+"-del").On(jquery.CLICK, func() {

			svless := reflect.Zero(typ)
			for j := 0; j < val.Len(); j++ {
				fmt.Println("j :", j)
				if j == i_closed_over {
					fmt.Println("index skipped :", j)
					continue
				}
				svless = reflect.Append(svless, val.Index(j))
			}

			node.Object = svless.Interface()
			s.handle(node, cb)
			cb(svless.Interface())
		})
	}

	jQuery("#" + node.EditorId).Append("<div class=\"row\" id=\"" + node.EditorId + "-add-row\"></div>")
	jQuery("#" + node.EditorId + "-add-row").Append("<div class=\"col-lg-1\" id=\"" + node.EditorId + "-add-col-L\"></div>")
	jQuery("#" + node.EditorId + "-add-row").Append("<div class=\"col-lg-11\" id=\"" + node.EditorId + "-add-col-R\"><hr /> </div>")
	jQuery("#" + node.EditorId + "-add-col-L").Append("<button type=\"button\" id=\"" + node.EditorId + "-add\" class=\"btn btn-success btn-xs\"><i class=\"fa fa-plus\"></i></button>")

	jQuery("#"+node.EditorId+"-add").On(jquery.CLICK, func() {
		if e.Kind() == reflect.Ptr {
			val = reflect.Append(val, reflect.New(e.Elem()))
		} else {
			val = reflect.Append(val, reflect.Zero(e))
		}

		node.Object = val.Interface()
		s.handle(node, cb)
		cb(val.Interface())
	})
}
