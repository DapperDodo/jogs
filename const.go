package jogs

import (
	"html/template"
	"strings"

	"github.com/gopherjs/jquery"
)

/*
	usage:
	type Example struct {
		Fubr  string   `jogs:"CONST,Fuber,select fuber...,FOO BAR QUX FOX"`
	}
*/
type ConstSelector struct {
	skin *template.Template
}

func newConstSelector() *ConstSelector {
	return &ConstSelector{skin: template.Must(template.New("skin").Parse(string(`

		{{define "handle"}}
			<div class="form-group" id="{{.EditorId}}">
			{{if ne .Label ""}}
				<label class="control-label">{{.Label}}</label>
			{{end}}
			</div>
		{{end}}

		{{define "show"}}
			<input class="form-control" id="{{.EditorId}}-show" type="text" placeholder="{{if ne .Placeholder ""}}{{.Placeholder}}{{else}}string constant...{{end}}" value="{{.Object}}">
		{{end}}

		{{define "form"}}
			<div class="form-group input-group" id="{{ .EditorId }}-edit">
				<select class="form-control" id="{{ .EditorId }}-edit-input">
					<option value="">select constant...</option>
					<option value=""></option>
			{{range .List}}
				{{if .Selected}}
					<option value={{ .Val }} selected>{{ .Val }}</option>
				{{else}}
					<option value={{ .Val }}>{{ .Val }}</option>
				{{end}}
				}
			{{end}}
				</select>
			</div>
		{{end}}

	`)))}
}

func (h *ConstSelector) Handle(node Node, cb Callback) {
	J("#" + node.ContainerId).Append(Merge(h.skin, "handle", node))
	h.show(node, cb)
}

//////////////////////////////////////////////////////////////////////////////////

func (h *ConstSelector) show(node Node, cb Callback) {
	J("#" + node.EditorId).Append(Merge(h.skin, "show", node))
	J("#"+node.EditorId+"-show").On(jquery.FOCUS, func() {
		h.form(node, cb)
	})
}

type myconstant struct {
	Val      string
	Selected bool
}
type nodeconstants struct {
	*Node
	List []myconstant
}

func (h *ConstSelector) form(node Node, cb Callback) {

	value_config := node.GetTag(PARAM_1)
	value_fields := strings.Fields(value_config)
	values := make([]myconstant, len(value_fields))
	for idx, val := range value_fields {
		if val == node.Object.(string) {
			values[idx] = myconstant{val, true}
		} else {
			values[idx] = myconstant{val, false}
		}
	}

	ns := nodeconstants{Node: &node, List: values}

	J("#" + node.EditorId + "-show").Remove()
	J("#" + node.EditorId).Append(Merge(h.skin, "form", ns))
	J("#"+node.EditorId+"-edit-input").Focus().Select().SetAttr("size", len(ns.List)+2).On(jquery.CHANGE, func() {
		go h.save(node, cb)
	}).On(jquery.BLUR, func() {
		go h.save(node, cb)
	})
}

func (h *ConstSelector) save(node Node, cb Callback) {
	constant := J("#" + node.EditorId + "-edit-input").Val()
	J("#" + node.EditorId + "-edit").Remove()
	node.Object = constant
	go h.show(node, cb)
	cb(constant)
}
