package controller

import (
	"github.com/kataras/iris/v12"
	"github.com/zoulongbo/go-mall/types"
)

type ExampleController struct {

}

// GET: http://localhost:8888/books
func (c *ExampleController) Get() []types.Example {
	return []types.Example{
		{"Mastering Concurrency in Go"},
		{"Go Design Patterns"},
		{"Black Hat Go"},
	}
}

// POST: http://localhost:8888/books
func (c *ExampleController) Post(b types.Example) int {
	println("Received Book: " + b.Title)

	return iris.StatusCreated
}