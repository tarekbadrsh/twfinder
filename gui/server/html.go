// Defines the Html component.

package server

// HTML interface defines a component which wraps an HTML text into a component.
//
// Default style class: "gui-HTML"
type HTML interface {
	// HTML is a component.
	Comp

	// HTML returns the HTML text.
	HTML() string

	// SetHTML sets the HTML text.
	SetHTML(html string)
}

// HTML implementation
type htmlImpl struct {
	compImpl // Component implementation

	html string // HTML text
}

// NewHTML creates a new HTML.
func NewHTML(html string) HTML {
	c := &htmlImpl{newCompImpl(nil), html}
	c.Style().AddClass("gui-Html")
	return c
}

func (c *htmlImpl) HTML() string {
	return c.html
}

func (c *htmlImpl) SetHTML(html string) {
	c.html = html
}

func (c *htmlImpl) Render(w Writer) {
	w.Write(strSpanOp)
	c.renderAttrsAndStyle(w)
	c.renderEHandlers(w)
	w.Write(strGT)

	w.Writes(c.html)

	w.Write(strSpanCl)
}
