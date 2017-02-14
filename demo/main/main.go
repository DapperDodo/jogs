// +build js
package main

import (
	"encoding/json"

	"github.com/gopherjs/jquery"

	"github.com/DapperDodo/jogs"
	"github.com/DapperDodo/jogs/demo/custom_handler"
)

type Data struct {
	Id     int
	GameId string `jogs:"custom"`
	Dialog []string
	Nested Nested
	Sa     string
	Sb     string
	Sc     string
	Aa     []string
	Bb     []string
	Cc     []string
}

type Nested struct {
	Id   int
	Name string
	Nums []int
}

func main() {

	data := &Data{
		Id:     1,
		GameId: "444-555-666",
		Dialog: []string{"foo", "bar", "qux"},
		Nested: Nested{2, "dodo", []int{10, 20}},
	}

	showData(data)

	// register custom handlers
	custom_handler.RegisterAll(jogs.DefaultDispatcher)

	// start editing
	jogs.Root(jogs.DefaultDispatcher, "panel-content", data, func(data_out interface{}) {
		showData(data)
	})
}

func showData(data interface{}) {
	d, _ := json.MarshalIndent(data, "<br/>", "&nbsp;&nbsp;&nbsp;&nbsp;")
	jquery.NewJQuery("#panel-title").SetHtml(string(d))
}
