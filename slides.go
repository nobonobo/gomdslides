package main

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"syscall/js"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
)

// Fragments ...
const Fragments = ".active .fragment, .forwardIn .fragment"

// ParseHash ...
func ParseHash(s string) int {
	orig, _ := url.Parse(s)
	u, _ := url.Parse(orig.Fragment)
	if u.Path == "" {
		u.Path = "1"
	}
	v, _ := strconv.Atoi(u.Path)
	return v
}

// GetURL ...
func GetURL() int {
	return ParseHash(location.String())
}

// Slides ...
type Slides struct {
	vecty.Core
	Contents    []*Slide     `vecty:"prop"`
	Index       int          `vecty:"prop"`
	Last        int          `vecty:"prop"`
	Controller  *Controller  `vecty:"prop"`
	SearchLight *SearchLight `vecty:"prop"`
	Cursor      string       `vecty:"prop"`
}

// Mount ...
func (c *Slides) Mount() {
	c.Navigate(location.Call("toString").String())
	c.Last = c.Index
}

// Navigate ...
func (c *Slides) Navigate(u string) {
	newID := ParseHash(u)
	log.Println("navigate to:", u, newID)
	c.Index = newID - 1
	vecty.Rerender(c)
}

// Prev ...
func (c *Slides) Prev(event *vecty.Event) {
	if c.Index == 0 {
		return
	}
	location.Set("hash", fmt.Sprintf("#%d", c.Index))
}

// Next ...
func (c *Slides) Next(event *vecty.Event) {
	if e := document.Call("querySelector", Fragments); e != js.Null() {
		cs := e.Get("classList")
		cs.Call("remove", "fragment")
		cs.Call("add", "appeared")
		vecty.Rerender(c.Contents[c.Index])
		return
	}
	if c.Index < len(c.Contents)-1 {
		location.Set("hash", fmt.Sprintf("#%d", c.Index+2))
	}
}

// OnHashChange ...
func (c *Slides) OnHashChange(this js.Value, args []js.Value) interface{} {
	oldID := ParseHash(args[0].Get("oldURL").String())
	newID := ParseHash(args[0].Get("newURL").String())
	log.Println("change:", oldID, "->", newID)
	c.Index = newID - 1
	c.Last = oldID - 1
	vecty.Rerender(c)
	return nil
}

// OnKeyDown ...
func (c *Slides) OnKeyDown(this js.Value, args []js.Value) interface{} {
	event := args[0]
	switch v := event.JSValue().Get("code").String(); v {
	case "ArrowLeft":
		log.Println("prev:", v)
		c.Prev(nil)
	case "ArrowRight":
		log.Println("next:", v)
		c.Next(nil)
	case "KeyF":
		document.Get("body").Call("requestFullscreen")
	case "KeyS":
		c.SearchLight.Enabled = !c.SearchLight.Enabled
		vecty.Rerender(c.SearchLight)
	case "KeyR":
		go func() {
			c.Contents = LoadMarkdown()
			c.Last = c.Index
			c.Refresh()
		}()
	default:
		log.Println(v)
	}
	return nil
}

// OnMouseMove ...
func (c *Slides) OnMouseMove(event *vecty.Event) {
	c.SearchLight.OnMouseMove(event)
	cursor := "inherit"
	if c.SearchLight.Active {
		cursor = "none"
	}
	if cursor != c.Cursor {
		c.Cursor = cursor
		c.Last = c.Index
		vecty.Rerender(c)
	}
}

// Refresh ...
func (c *Slides) Refresh() {
	vecty.RenderBody(c)
}

// Render ...
func (c *Slides) Render() vecty.ComponentOrHTML {
	vecty.SetTitle(fmt.Sprintf("Slide(%d/%d)", c.Index+1, len(c.Contents)))
	var contents vecty.List
	for i, slide := range c.Contents {
		contents = append(contents, slide)
		switch {
		case i < c.Index-1:
			slide.State = "prev"
		case i == c.Index-1:
			switch {
			case i == c.Last:
				slide.State = "forwardOut"
			default:
				slide.State = "prev"
			}
		case i == c.Index:
			switch {
			case i-1 == c.Last:
				slide.State = "forwardIn"
			case i+1 == c.Last:
				slide.State = "reverseIn"
			default:
				slide.State = "active"
			}
		case i == c.Index+1:
			switch {
			case i == c.Last:
				slide.State = "reverseOut"
			default:
				slide.State = "next"
			}
		case i > c.Index+1:
			slide.State = "next"
		}
	}
	return elem.Body(
		elem.Div(
			vecty.Markup(
				vecty.Class("container"),
				event.MouseDown(c.SearchLight.OnMouseDown),
				event.MouseOut(c.SearchLight.OnMouseOut),
				event.MouseMove(c.OnMouseMove),
				event.Wheel(c.SearchLight.OnWheel),
			),
			contents,
		),
		c.Controller,
		c.SearchLight,
	)
}
