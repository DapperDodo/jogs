// +build js
package main

import (
	"encoding/json"

	"github.com/DapperDodo/jogs"
	"github.com/gopherjs/jquery"

	"backend/api"
)

func main() {

	data := &api.Aardvark{}

	// show the empty object that we will edit
	showData(data)

	// start editing
	jogs.Root(
		jogs.DefaultDispatcher,
		"panel-content", // the target element id where the editor must be rendered
		data,            // the object to edit
		showData,        // the function that is called when editing is done
	)
}

// showData is a helper function that will display the data in json format
func showData(data interface{}) {
	d, _ := json.MarshalIndent(data, "<br/>", "&nbsp;&nbsp;&nbsp;&nbsp;")
	jquery.NewJQuery("#panel-title").SetHtml(string(d))
}
