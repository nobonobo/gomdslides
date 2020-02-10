package main

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
)

// Slide ...
type Slide struct {
	vecty.Core
	State   string `vecty:"prop"`
	Content string `vecty:"prop"`
}

// OnAnimationEnd ...
func (c *Slide) OnAnimationEnd(event *vecty.Event) {
	switch c.State {
	case "forwardOut":
		c.State = "prev"
	case "forwardIn", "reverseIn":
		c.State = "active"
	case "reverseOut":
		c.State = "next"
	}
	vecty.Rerender(c)
}

// Render ...
func (c *Slide) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			event.AnimationEnd(c.OnAnimationEnd),
			vecty.Class("card", c.State),
		),
		elem.Div(
			vecty.Markup(
				vecty.Class("card-body"),
			),
			elem.Div(
				vecty.Markup(
					vecty.Class("content"),
					vecty.UnsafeHTML(c.Content),
				),
			),
		),
	)
}
