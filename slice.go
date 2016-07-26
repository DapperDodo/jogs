package jogs

import (
	"html/template"
	"reflect"
	"strconv"

	"github.com/gopherjs/jquery"
)

type slice struct {
	*Dispatcher
	skin *template.Template
}

func newSlice(d *Dispatcher) *slice {
	return &slice{d, template.Must(template.New("skin").Parse(string(`
		{{define "row"}}
			<div class="row" id="{{.EditorId}}-row\">
				<div class="col-lg-1" id="{{.EditorId}}-col-L">
					<button type="button" id="{{.EditorId}}-del" class="btn btn-danger btn-md" tabindex="-1">
						<i class="fa fa-trash"></i>
					</button>
				</div>				
				<div class="col-lg-11" id="{{.ContainerId}}">
				</div>
			</div>
		{{end}}

		{{define "add"}}
			<div class="row" id="{{.EditorId}}-add-row">
				<div class="col-lg-1" id="{{.EditorId}}-add-col-L">
					<button type="button" id="{{.EditorId}}-add" class="btn btn-success btn-md" tabindex="-1">
						<i class="fa fa-plus"></i>
					</button>
				</div>
				<div class="col-lg-11" id="{{.EditorId}}-add-col-R">
					<hr />
				</div>
			</div>
		{{end}}
	`)))}
}

func (s *slice) Handle(node Node, cb Callback) {

	typ := reflect.TypeOf(node.Object)
	val := reflect.ValueOf(node.Object)

	J("#" + node.EditorId).Empty()

	e := typ.Elem()

	for i := 0; i < val.Len(); i++ {

		valrow := val.Index(i)

		noderow := node
		noderow.Idx = i
		noderow.EditorId += "-" + strconv.Itoa(i)
		noderow.ContainerId = noderow.EditorId + "-col-R"

		J("#" + node.EditorId).Append(Merge(s.skin, "row", noderow))

		switch valrow.Kind() {
		case reflect.String, reflect.Int, reflect.Float64, reflect.Bool:
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
			J("#" + node.EditorId).Empty()
			noderow.Handle = "UNSUPPORTED"
			noderow.Object = "unsupported slice type"
			s.dispatch(noderow, nil)
			return
		}

		i_closed_over := i
		J("#"+node.EditorId+"-"+strconv.Itoa(i)+"-del").On(jquery.CLICK, func() {

			svless := reflect.Zero(typ)
			for j := 0; j < val.Len(); j++ {
				if j == i_closed_over {
					continue
				}
				svless = reflect.Append(svless, val.Index(j))
			}

			node.Object = svless.Interface()
			s.Handle(node, cb)
			cb(svless.Interface())
		})
	}

	J("#" + node.EditorId).Append(Merge(s.skin, "add", node))
	J("#"+node.EditorId+"-add").On(jquery.CLICK, func() {
		if e.Kind() == reflect.Ptr {
			val = reflect.Append(val, reflect.New(e.Elem()))
		} else {
			val = reflect.Append(val, reflect.Zero(e))
		}

		node.Object = val.Interface()
		s.Handle(node, cb)
		cb(val.Interface())
	})
}
