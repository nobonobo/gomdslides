package main

import (
	"strconv"
	"syscall/js"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/prop"
)

// SearchLight ...
type SearchLight struct {
	vecty.Core
	tid     js.Value
	Enabled bool   `vecty:"prop"`
	Active  bool   `vecty:"prop"`
	Left    string `vecty:"prop"`
	Top     string `vecty:"prop"`
}

// OnMouseDown ...
func (c *SearchLight) OnMouseDown(event *vecty.Event) {
	if event.Value.Get("which").Int() == 2 {
		c.Enabled = !c.Enabled
		c.Active = c.Enabled
		vecty.Rerender(c)
		event.Value.Call("preventDefault")
		if sp := event.Value.Get("stopPropagation"); sp != js.Undefined() {
			event.Value.Call("stopPropagation")
		}
	}
}

// OnMouseMove ...
func (c *SearchLight) OnMouseMove(event *vecty.Event) {
	c.Active = c.Enabled
	x, y := event.Value.Get("pageX"), event.Value.Get("pageY")
	c.Left = strconv.Itoa(x.Int()-150) + "px"
	c.Top = strconv.Itoa(y.Int()-150) + "px"
	vecty.Rerender(c)
}

// OnMouseOut ...
func (c *SearchLight) OnMouseOut(event *vecty.Event) {
	if c.Active {
		c.Active = false
		vecty.Rerender(c)
	}
}

// OnWheel ...
func (c *SearchLight) OnWheel(event *vecty.Event) {
	if c.Enabled {
		// spotlight-onの時はwheelイベントを親に伝搬させない
		if sp := event.Value.Get("stopPropagation"); sp != js.Undefined() {
			event.Value.Call("stopPropagation")
		}
	}
}

// Render ...
func (c *SearchLight) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			prop.ID("searchlight"),
			vecty.ClassMap{
				"on": c.Active,
			},
			vecty.Style("margin-left", c.Left),
			vecty.Style("margin-top", c.Top),
		),
	)
}
