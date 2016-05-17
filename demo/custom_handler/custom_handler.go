package custom_handler

import (
	"fmt"

	"github.com/gopherjs/jquery"

	"github.com/DapperDodo/jogs"
)

var jQuery = jquery.NewJQuery

func RegisterAll(d *jogs.Dispatcher) {
	d.Register("custom", jogs.HandlerFunc(HandleCustom))
}

// HandleCustom is an example plugin function for jogs
// it implements a green textarea editor instead of the standard string editor shipped with jogs
func HandleCustom(node jogs.Node, cb jogs.Callback) {
	jQuery("#" + node.ContainerId).Append("<div class=\"form-group has-success\" id=\"" + node.EditorId + "\"></div>")
	if node.Label != "" {
		jQuery("#" + node.EditorId).Append("<label class=\"control-label\">" + node.Label + "</label><p class=\"help-block\">This is a custom string editor.</p>")
	}
	showString(node, cb)
}

func showString(node jogs.Node, cb jogs.Callback) {
	jQuery("#" + node.EditorId).Append("<input class=\"form-control\" id=\"" + node.EditorId + "-show\" type=\"text\" placeholder=\"" + fmt.Sprint(node.Object) + "\">")
	jQuery("#"+node.EditorId+"-show").On(jquery.CLICK, func() {
		editString(node, cb)
	})
}

func editString(node jogs.Node, cb jogs.Callback) {
	jQuery("#" + node.EditorId + "-show").Remove()
	jQuery("#" + node.EditorId).Append("<div class=\"form-group input-group\" id=\"" + node.EditorId + "-edit\"><textarea id=\"" + node.EditorId + "-edit-input\" class=\"form-control\" rows=\"3\">" + fmt.Sprint(node.Object) + "</textarea><span class=\"input-group-btn\"><button class=\"btn btn-default\"><i class=\"fa fa-save\"></i></button></span></div>")
	jQuery("#"+node.EditorId+"-edit-input").Focus().Select().On(jquery.BLUR, func() {
		saveString(node, cb)
	})
}

func saveString(node jogs.Node, cb jogs.Callback) {
	str := jQuery("#" + node.EditorId + "-edit-input").Val()
	jQuery("#" + node.EditorId + "-edit").Remove()
	node.Object = str
	showString(node, cb)
	cb(str)
}
