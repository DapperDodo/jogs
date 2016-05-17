package jogs

import (
	"fmt"
	"reflect"
	"strconv"

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

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func handleUnsupported(node Node, cb Callback) {
	jQuery("#" + node.ContainerId).Append("<div class=\"form-group has-warning\" id=\"" + node.EditorId + "\"><label class=\"control-label\" for=\"inputWarning\">" + node.Label + "</label></div>")
	jQuery("#" + node.EditorId).Append("<input class=\"form-control\" id=\"inputWarning\" type=\"text\" placeholder=\"" + node.Object.(string) + "\" disabled>")
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func handleInt(node Node, cb Callback) {
	jQuery("#" + node.ContainerId).Append("<div class=\"form-group\" id=\"" + node.EditorId + "\"></div>")
	if node.Label != "" {
		jQuery("#" + node.EditorId).Append("<label class=\"control-label\">" + node.Label + "</label>")
	}
	showInt(node, cb)
}

func showInt(node Node, cb Callback) {
	jQuery("#" + node.EditorId).Append("<input class=\"form-control\" id=\"" + node.EditorId + "-show\" type=\"text\" placeholder=\"" + fmt.Sprint(node.Object) + "\">")
	jQuery("#"+node.EditorId+"-show").On(jquery.CLICK, func() {
		formInt(node, cb)
	})
}

func formInt(node Node, cb Callback) {
	jQuery("#" + node.EditorId + "-show").Remove()
	jQuery("#" + node.EditorId).Append("<div class=\"form-group input-group\" id=\"" + node.EditorId + "-edit\"><input class=\"form-control\" id=\"" + node.EditorId + "-edit-input\" type=\"text\" value=\"" + fmt.Sprint(node.Object) + "\"><span class=\"input-group-btn\"><button class=\"btn btn-default\"><i class=\"fa fa-save\"></i></button></span></div>")
	jQuery("#"+node.EditorId+"-edit-input").Focus().Select().On(jquery.BLUR, func() {
		saveInt(node, cb)
	})
}

func saveInt(node Node, cb Callback) {
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
	showInt(node, cb)

	cb(int(val32))
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func handleString(node Node, cb Callback) {
	jQuery("#" + node.ContainerId).Append("<div class=\"form-group\" id=\"" + node.EditorId + "\"></div>")
	if node.Label != "" {
		jQuery("#" + node.EditorId).Append("<label class=\"control-label\">" + node.Label + "</label>")
	}
	showString(node, cb)
}

func showString(node Node, cb Callback) {
	jQuery("#" + node.EditorId).Append("<input class=\"form-control\" id=\"" + node.EditorId + "-show\" type=\"text\" placeholder=\"" + fmt.Sprint(node.Object) + "\">")
	jQuery("#"+node.EditorId+"-show").On(jquery.CLICK, func() {
		formString(node, cb)
	})
}

func formString(node Node, cb Callback) {
	jQuery("#" + node.EditorId + "-show").Remove()
	jQuery("#" + node.EditorId).Append("<div class=\"form-group input-group\" id=\"" + node.EditorId + "-edit\"><input class=\"form-control\" id=\"" + node.EditorId + "-edit-input\" type=\"text\" value=\"" + fmt.Sprint(node.Object) + "\"><span class=\"input-group-btn\"><button class=\"btn btn-default\"><i class=\"fa fa-save\"></i></button></span></div>")
	jQuery("#"+node.EditorId+"-edit-input").Focus().Select().On(jquery.BLUR, func() {
		saveString(node, cb)
	})
}

func saveString(node Node, cb Callback) {
	str := jQuery("#" + node.EditorId + "-edit-input").Val()
	jQuery("#" + node.EditorId + "-edit").Remove()

	node.Object = str
	showString(node, cb)

	cb(str)
}
