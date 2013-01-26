package fimika
/*

Getting started
===============

::
	package hello
	
	import "fimika"
	
	var app = fimika.NewApplication()
	
	func init() {
		app.AddRoute("/", handleIndex)
		app.Start()
	}
	
	func handleIndex(c *fimika.Context) {
		c.Raw("Hello Fimika!")
	}

*/
