package main

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
)

// Controller ...
type Controller struct {
	vecty.Core
	parent *Slides
}

// Render ...
func (c *Controller) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			vecty.Class("controller"),
		),
		elem.Button(
			vecty.Markup(
				vecty.ClassMap{
					"btn":      true,
					"btn-link": true,
					"btn-lg":   true,
				},
				event.Click(c.parent.Prev).PreventDefault(),
			),
			vecty.Text("<"),
		),
		elem.Button(
			vecty.Markup(
				vecty.ClassMap{
					"btn":      true,
					"btn-link": true,
					"btn-lg":   true,
				},
				event.Click(c.parent.Next).PreventDefault(),
			),
			vecty.Text(">"),
		),
	)
}
