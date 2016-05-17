# jogs

generate editors (web forms) from go structs, specifically meant for developing single page web apps

### installation

	go get -u github.com/DapperDodo/jogs

### simple use: main.go

	// +build js
	package main

	import (
		"encoding/json"

		"github.com/gopherjs/jquery"

		"github.com/DapperDodo/jogs"
	)

	type Data struct {
		Id     int
		GameId string
		Dialog []string
	}

	func main() {

		data := &Data{
			Id:     1,
			GameId: "444-555-666",
			Dialog: []string{"foo", "bar", "qux"},
		}

		showData(data)

		// start editing
		jogs.Root(jogs.DefaultDispatcher, "panel-content", data, func(data_out interface{}) {
			showData(data)
		})
	}

	func showData(data interface{}) {
		d, _ := json.MarshalIndent(data, "<br/>", "&nbsp;&nbsp;&nbsp;&nbsp;")
		jquery.NewJQuery("#panel-title").SetHtml(string(d))
	}

run this example using gopherjs' built in server:

	gopherjs serve

and navigate a browser to:

	http://localhost:8080/github.com/DapperDodo/jogs/demo/main/

### advanced use: main.go

Jogs can handle complex data structures, for example nesting and slices are fully supported. (TODO: maps, floats and bools)

Also, Jogs sports a powerful plugin structure that makes overriding and extending easy. 
Custom handlers can be registered, after which struct tags can be used to tell Jogs where to use these custom handlers.


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
