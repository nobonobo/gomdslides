package main

import (
	"bytes"
	"log"
	"syscall/js"
	"time"

	"github.com/gopherjs/vecty"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

var (
	window   = js.Global()
	document = js.Global().Get("document")
	console  = js.Global().Get("console")
	location = js.Global().Get("location")
)

// Setup ...
func Setup() *Slides {
	document.Get("body").Get("parentElement").Get("style").Set("backgroundColor", "white")
	meta := document.Call("createElement", "meta")
	meta.Call("setAttribute", "name", "viewport")
	meta.Call("setAttribute", "content", "width=device-width,initial-scale=1")
	document.Get("head").Call("append", meta)
	vecty.AddStylesheet("css/spectre.min.css")
	vecty.AddStylesheet("css/spectre-exp.min.css")
	vecty.AddStylesheet("css/spectre-icons.min.css")
	vecty.AddStylesheet("css/app.css")
	contents := LoadMarkdown()
	slides := &Slides{
		Contents:    contents,
		Controller:  &Controller{},
		SearchLight: &SearchLight{},
	}
	slides.Controller.parent = slides
	js.Global().Call("addEventListener", "hashchange", js.FuncOf(slides.OnHashChange))
	js.Global().Call("addEventListener", "keydown", js.FuncOf(slides.OnKeyDown))
	return slides
}

// LoadMarkdown ...
func LoadMarkdown() []*Slide {
	md := goldmark.New(
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle("monokai"),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAttribute(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithUnsafe(),
		),
	)
	ch := make(chan string)
	window.Call("fetch", time.Now().Format("contents.md?t=20060102T150405")).Call("then",
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			response := args[0]
			response.Call("arrayBuffer").Call("then",
				js.FuncOf(func(this js.Value, args []js.Value) interface{} {
					buff := window.Get("Uint8Array").New(args[0])
					source := make([]byte, buff.Get("length").Int())
					js.CopyBytesToGo(source, buff)
					defer close(ch)
					for _, chunk := range bytes.Split(source, []byte("\n====\n")) {
						chunk = bytes.Trim(chunk, "\r\n\t ")
						var output bytes.Buffer
						if err := md.Convert(chunk, &output); err != nil {
							log.Fatal(err)
						}
						ch <- output.String()
					}
					return nil
				}),
			)
			return nil
		}),
	)
	var contents []*Slide
	for s := range ch {
		contents = append(contents, &Slide{State: "next", Content: s})
	}
	return contents
}
