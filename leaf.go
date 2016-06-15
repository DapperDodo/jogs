package jogs

import (
	"bytes"
	"fmt"
	"html/template"
	"reflect"
	"strconv"

	"github.com/gopherjs/jquery"
)

type leaf struct {
	*Dispatcher
}

func (n *leaf) Handle(node Node, cb Callback) {

	typekind := reflect.TypeOf(node.Object).Kind()

	switch typekind {
	case reflect.Int:
		node.Handle = "LEAF_INT"
	case reflect.Float64:
		node.Handle = "LEAF_FLOAT"
	case reflect.String:
		node.Handle = "LEAF_STRING"
	case reflect.Bool:
		node.Handle = "LEAF_BOOL"
	default:
		node.Handle = "UNSUPPORTED"
		node.Object = fmt.Sprint("Unsupported leaf type (leaf must be one of int|float|string|bool):", typekind, "->", node.Object)
	}

	n.dispatch(node, cb)
}

func Merge(skin *template.Template, tpl string, data interface{}) string {
	var buf bytes.Buffer
	err := skin.ExecuteTemplate(&buf, tpl, data)
	if err != nil {
		return "" // TODO: log this...
	}
	return buf.String()
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type unsupported struct {
	skin *template.Template
}

func newUnsupported() *unsupported {
	return &unsupported{template.Must(template.New("skin").Parse(string(`
		{{define "handle"}}
			<div class="form-group has-warning" id="{{.EditorId}}">
				<label class="control-label" for="inputWarning">{{.Label}}</label>
				<input class="form-control" id="inputWarning" type="text" placeholder="{{.Object}}" disabled></input>
			</div>
		{{end}}
	`)))}
}

func (h *unsupported) Handle(node Node, cb Callback) {
	J("#" + node.ContainerId).Append(Merge(h.skin, "handle", node))
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type intHandler struct {
	skin *template.Template
}

func newIntHandler() *intHandler {
	return &intHandler{template.Must(template.New("skin").Parse(string(`
		{{define "handle"}}
			<div class="form-group" id="{{.EditorId}}">
			{{if ne .Label ""}}
				<label class="control-label">{{.Label}}</label>
			{{end}}
			</div>
		{{end}}

		{{define "show"}}
			<input class="form-control" id="{{.EditorId}}-show" type="text" placeholder="int..." value="{{.Object}}">
		{{end}}

		{{define "form"}}
			<div class="form-group input-group" id="{{.EditorId}}-edit">
				<input class="form-control" id="{{.EditorId}}-edit-input" type="text" placeholder="int..." value="{{.Object}}">
				<span class="input-group-btn">
					<button class="btn btn-default">
						<i class="fa fa-save"></i>
					</button>
				</span>
			</div>
		{{end}}

		{{define "form-error"}}
			<div id="{{.EditorId}}-error" class="alert alert-danger">
				Please fill in a number!
			</div>
		{{end}}
	`)))}
}

func (h *intHandler) Handle(node Node, cb Callback) {
	J("#" + node.ContainerId).Append(Merge(h.skin, "handle", node))
	h.show(node, cb)
}

func (h *intHandler) show(node Node, cb Callback) {
	J("#" + node.EditorId).Append(Merge(h.skin, "show", node))
	J("#"+node.EditorId+"-show").On(jquery.CLICK, func() {
		h.form(node, cb)
	})
}

func (h *intHandler) form(node Node, cb Callback) {
	J("#" + node.EditorId + "-show").Remove()
	J("#" + node.EditorId).Append(Merge(h.skin, "form", node))
	J("#"+node.EditorId+"-edit-input").Focus().Select().On(jquery.BLUR, func() {
		h.save(node, cb)
	})
}

func (h *intHandler) save(node Node, cb Callback) {

	val32, err := strconv.ParseInt(J("#"+node.EditorId+"-edit-input").Val(), 10, 32)
	if err != nil {
		J("#" + node.EditorId).AddClass("has-error")
		J("#" + node.EditorId + "-edit-input").Focus().Select()
		J("#" + node.EditorId + "-error").Remove()
		J("#" + node.EditorId).Append(Merge(h.skin, "form-error", node))
		return
	}
	J("#" + node.EditorId).RemoveClass("has-error")
	J("#" + node.EditorId + "-edit").Remove()
	J("#" + node.EditorId + "-error").Remove()

	node.Object = int(val32)
	h.show(node, cb)

	cb(int(val32))
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type floatHandler struct {
	skin *template.Template
}

func newFloatHandler() *floatHandler {
	return &floatHandler{template.Must(template.New("skin").Parse(string(`
		{{define "handle"}}
			<div class="form-group" id="{{.EditorId}}">
			{{if ne .Label ""}}
				<label class="control-label">{{.Label}}</label>
			{{end}}
			</div>
		{{end}}

		{{define "show"}}
			<input class="form-control" id="{{.EditorId}}-show" type="text" placeholder="float..." value="{{.Object}}">
		{{end}}

		{{define "form"}}
			<div class="form-group input-group" id="{{.EditorId}}-edit">
				<input class="form-control" id="{{.EditorId}}-edit-input" type="text" placeholder="float..." value="{{.Object}}">
				<span class="input-group-btn">
					<button class="btn btn-default">
						<i class="fa fa-save"></i>
					</button>
				</span>
			</div>
		{{end}}

		{{define "form-error"}}
			<div id="{{.EditorId}}-error" class="alert alert-danger">
				Please fill in a (floating point) number! i.e. "0", "3.1415" or "10"
			</div>
		{{end}}
	`)))}
}

func (h *floatHandler) Handle(node Node, cb Callback) {
	J("#" + node.ContainerId).Append(Merge(h.skin, "handle", node))
	h.show(node, cb)
}

func (h *floatHandler) show(node Node, cb Callback) {
	J("#" + node.EditorId).Append(Merge(h.skin, "show", node))
	J("#"+node.EditorId+"-show").On(jquery.CLICK, func() {
		h.form(node, cb)
	})
}

func (h *floatHandler) form(node Node, cb Callback) {
	J("#" + node.EditorId + "-show").Remove()
	J("#" + node.EditorId).Append(Merge(h.skin, "form", node))
	J("#"+node.EditorId+"-edit-input").Focus().Select().On(jquery.BLUR, func() {
		h.save(node, cb)
	})
}

func (h *floatHandler) save(node Node, cb Callback) {

	valf64, err := strconv.ParseFloat(J("#"+node.EditorId+"-edit-input").Val(), 64)
	if err != nil {
		J("#" + node.EditorId).AddClass("has-error")
		J("#" + node.EditorId + "-edit-input").Focus().Select()
		J("#" + node.EditorId + "-error").Remove()
		J("#" + node.EditorId).Append(Merge(h.skin, "form-error", node))
		return
	}
	J("#" + node.EditorId).RemoveClass("has-error")
	J("#" + node.EditorId + "-edit").Remove()
	J("#" + node.EditorId + "-error").Remove()

	node.Object = valf64
	h.show(node, cb)
	cb(valf64)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type stringHandler struct {
	skin *template.Template
}

func newStringHandler() *stringHandler {
	return &stringHandler{template.Must(template.New("skin").Parse(string(`
		{{define "handle"}}
			<div class="form-group" id="{{.EditorId}}">
			{{if ne .Label ""}}
				<label class="control-label">{{.Label}}</label>
			{{end}}
			</div>
		{{end}}

		{{define "show"}}
			<input class="form-control" id="{{.EditorId}}-show" type="text" placeholder="{{if ne .Placeholder ""}}{{.Placeholder}}{{else}}string...{{end}}" value="{{.Object}}">
		{{end}}

		{{define "form"}}
			<div class="form-group input-group" id="{{.EditorId}}-edit">
				<input class="form-control" id="{{.EditorId}}-edit-input" type="text" placeholder="{{if ne .Placeholder ""}}{{.Placeholder}}{{else}}string...{{end}}" value="{{.Object}}">
				<span class="input-group-btn">
					<button class="btn btn-default">
						<i class="fa fa-save"></i>
					</button>
				</span>
			</div>
		{{end}}
	`)))}
}

func (h *stringHandler) Handle(node Node, cb Callback) {
	J("#" + node.ContainerId).Append(Merge(h.skin, "handle", node))
	h.show(node, cb)
}

func (h *stringHandler) show(node Node, cb Callback) {
	J("#" + node.EditorId).Append(Merge(h.skin, "show", node))
	J("#"+node.EditorId+"-show").On(jquery.CLICK, func() {
		h.form(node, cb)
	})
}

func (h *stringHandler) form(node Node, cb Callback) {
	J("#" + node.EditorId + "-show").Remove()
	J("#" + node.EditorId).Append(Merge(h.skin, "form", node))
	J("#"+node.EditorId+"-edit-input").Focus().Select().On(jquery.BLUR, func() {
		h.save(node, cb)
	})
}

func (h *stringHandler) save(node Node, cb Callback) {
	str := J("#" + node.EditorId + "-edit-input").Val()
	J("#" + node.EditorId + "-edit").Remove()
	node.Object = str
	h.show(node, cb)
	cb(str)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type boolHandler struct {
	skin *template.Template
}

func newBoolHandler() *boolHandler {
	return &boolHandler{template.Must(template.New("skin").Parse(string(`
		{{define "handle"}}
			<div class="form-group" id="{{.EditorId}}">
			{{if ne .Label ""}}
				<label class="control-label">{{.Label}}</label>
			{{end}}
			</div>
		{{end}}

		{{define "show"}}
			<input class="form-control" id="{{.EditorId}}-show" type="text" placeholder="boolean..." value="{{.Object}}">
		{{end}}

		{{define "form"}}
			<div class="form-group input-group" id="{{.EditorId}}-edit">
				<input class="form-control" id="{{.EditorId}}-edit-input" type="text" placeholder="boolean..." value="{{.Object}}">
				<span class="input-group-btn">
					<button class="btn btn-default">
						<i class="fa fa-save"></i>
					</button>
				</span>
			</div>
		{{end}}

		{{define "form-error"}}
			<div id="{{.EditorId}}-error" class="alert alert-danger">
				Please fill in a boolean! i.e. "0", "1", "true", "True", "TRUE", "false", "False", "FALSE", "t", "f", "T", "F"
			</div>
		{{end}}		
	`)))}
}

func (h *boolHandler) Handle(node Node, cb Callback) {
	J("#" + node.ContainerId).Append(Merge(h.skin, "handle", node))
	h.show(node, cb)
}

func (h *boolHandler) show(node Node, cb Callback) {
	J("#" + node.EditorId).Append(Merge(h.skin, "show", node))
	J("#"+node.EditorId+"-show").On(jquery.CLICK, func() {
		h.form(node, cb)
	})
}

func (h *boolHandler) form(node Node, cb Callback) {
	J("#" + node.EditorId + "-show").Remove()
	J("#" + node.EditorId).Append(Merge(h.skin, "form", node))
	J("#"+node.EditorId+"-edit-input").Focus().Select().On(jquery.BLUR, func() {
		h.save(node, cb)
	})
}

func (h *boolHandler) save(node Node, cb Callback) {
	boolval, err := strconv.ParseBool(J("#" + node.EditorId + "-edit-input").Val())
	if err != nil {
		J("#" + node.EditorId).AddClass("has-error")
		J("#" + node.EditorId + "-edit-input").Focus().Select()
		J("#" + node.EditorId + "-error").Remove()
		J("#" + node.EditorId).Append(Merge(h.skin, "form-error", node))
		return
	}
	J("#" + node.EditorId).RemoveClass("has-error")
	J("#" + node.EditorId + "-edit").Remove()
	J("#" + node.EditorId + "-error").Remove()

	node.Object = boolval
	h.show(node, cb)
	cb(boolval)
}
