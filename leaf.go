package jogs

import (
	"html/template"
	"fmt"
	"reflect"
	"strconv"
	"bytes"

	"github.com/gopherjs/jquery"
)

var jQuery = jquery.NewJQuery

type leaf struct {
	*Dispatcher
}

func (n *leaf) handle(node Node, cb Callback) {

	typekind := reflect.TypeOf(node.Object).Kind()

	switch typekind {
	case reflect.Int:
		node.Handle = "LEAF_INT"
	case reflect.String:
		node.Handle = "LEAF_STRING"
	default:
		node.Handle = "UNSUPPORTED"
		node.Object = fmt.Sprint("Unsupported leaf type (leaf must be one of int|string|float|bool):", typekind, "->", node.Object)
	}

	n.dispatch(node, cb)
}

func merge(skin *template.Template, tpl string, data interface{}) string {
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

	skin_tpl := `
		{{define "handle"}}
			<div class="form-group has-warning" id="{{.EditorId}}">
				<label class="control-label" for="inputWarning">{{.Label}}</label>
				<input class="form-control" id="inputWarning" type="text" placeholder="{{.Object}}" disabled></input>
			</div>
		{{end}}
	`
	skin := template.Must(template.New("skin").Parse(string(skin_tpl)))
	return &unsupported{skin}
}

func (h *unsupported) handle(node Node, cb Callback) {
	jQuery("#" + node.ContainerId).Append(merge(h.skin, "handle", node))
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type intHandler struct {
	skin *template.Template
}

func newIntHandler() *intHandler {
	skin_tpl := `
		{{define "handle"}}
			<div class="form-group" id="{{.EditorId}}">
			{{if ne .Label ""}}
				<label class="control-label">{{.Label}}</label>
			{{end}}
			</div>
		{{end}}

		{{define "show"}}
			<input class="form-control" id="{{.EditorId}}-show" type="text" placeholder="{{.Object}}">
		{{end}}

		{{define "form"}}
			<div class="form-group input-group" id="{{.EditorId}}-edit">
				<input class="form-control" id="{{.EditorId}}-edit-input" type="text" value="{{.Object}}">
				<span class="input-group-btn">
					<button class="btn btn-default">
						<i class="fa fa-save"></i>
					</button>
				</span>
			</div>
		{{end}}
	`	
	skin := template.Must(template.New("skin").Parse(string(skin_tpl)))
	return &intHandler{skin}
}

func (h *intHandler) handle(node Node, cb Callback) {
	jQuery("#" + node.ContainerId).Append(merge(h.skin, "handle", node))
	h.show(node, cb)
}

func (h *intHandler) show(node Node, cb Callback) {
	jQuery("#" + node.EditorId).Append(merge(h.skin, "show", node))
	jQuery("#"+node.EditorId+"-show").On(jquery.CLICK, func() {
		h.form(node, cb)
	})
}

func (h *intHandler) form(node Node, cb Callback) {
	jQuery("#" + node.EditorId + "-show").Remove()
	jQuery("#" + node.EditorId).Append(merge(h.skin, "form", node))
	jQuery("#"+node.EditorId+"-edit-input").Focus().Select().On(jquery.BLUR, func() {
		h.save(node, cb)
	})
}

func (h *intHandler) save(node Node, cb Callback) {

	val32, err := strconv.ParseInt(jQuery("#"+node.EditorId+"-edit-input").Val(), 10, 32)
	if err != nil {
		jQuery("#" + node.EditorId).AddClass("has-error")
		jQuery("#" + node.EditorId + "-edit-input").Focus().Select()
		jQuery("#" + node.EditorId + "-help").Remove()
		jQuery("#" + node.EditorId).Append("<p class=\"help-block\" id=\"" + node.EditorId + "-help\">Please fill in a number!</p>")
		return
	}
	jQuery("#" + node.EditorId).RemoveClass("has-error")
	jQuery("#" + node.EditorId + "-edit").Remove()
	jQuery("#" + node.EditorId + "-help").Remove()

	node.Object = int(val32)
	h.show(node, cb)

	cb(int(val32))
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type stringHandler struct {
	skin *template.Template
}

func newStringHandler() *stringHandler {
	skin_tpl := `
		{{define "handle"}}
			<div class="form-group" id="{{.EditorId}}">
			{{if ne .Label ""}}
				<label class="control-label">{{.Label}}</label>
			{{end}}
			</div>
		{{end}}

		{{define "show"}}
			<input class="form-control" id="{{.EditorId}}-show" type="text" placeholder="{{.Object}}">
		{{end}}

		{{define "form"}}
			<div class="form-group input-group" id="{{.EditorId}}-edit">
				<input class="form-control" id="{{.EditorId}}-edit-input" type="text" value="{{.Object}}">
				<span class="input-group-btn">
					<button class="btn btn-default">
						<i class="fa fa-save"></i>
					</button>
				</span>
			</div>
		{{end}}
	`	
	skin := template.Must(template.New("skin").Parse(string(skin_tpl)))
	return &stringHandler{skin}
}

func (h *stringHandler) handle(node Node, cb Callback) {
	jQuery("#" + node.ContainerId).Append(merge(h.skin, "handle", node))
	h.show(node, cb)
}

func (h *stringHandler) show(node Node, cb Callback) {
	jQuery("#" + node.EditorId).Append(merge(h.skin, "show", node))
	jQuery("#"+node.EditorId+"-show").On(jquery.CLICK, func() {
		h.form(node, cb)
	})
}

func (h *stringHandler) form(node Node, cb Callback) {
	jQuery("#" + node.EditorId + "-show").Remove()
	jQuery("#" + node.EditorId).Append(merge(h.skin, "form", node))
	jQuery("#"+node.EditorId+"-edit-input").Focus().Select().On(jquery.BLUR, func() {
		h.save(node, cb)
	})
}

func (h *stringHandler) save(node Node, cb Callback) {
	str := jQuery("#" + node.EditorId + "-edit-input").Val()
	jQuery("#" + node.EditorId + "-edit").Remove()
	node.Object = str
	h.show(node, cb)
	cb(str)
}
