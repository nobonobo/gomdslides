package main

import (
	"log"

	"github.com/gopherjs/vecty"
)

func main() {
	log.SetFlags(log.Lshortfile)
	slides := Setup()
	vecty.RenderBody(slides)
	select {}
}
