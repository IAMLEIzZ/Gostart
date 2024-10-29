package main

import (
	"gee"
	"net/http"
)

func main() {
	// r := gee.Default()
	r := gee.NewEngine()
	r.Use(gee.Logger(), gee.Recovery())
	r.Get("/", func(c *gee.Context) {
		c.String(http.StatusOK, "Hello LEI\n")
	})
	// index out of range for testing Recovery()
	r.Get("/panic", func(c *gee.Context) {
		names := []string{"IAMLEIzZ"}
		c.String(http.StatusOK, names[100])
	})

	r.Run(":9999")
}