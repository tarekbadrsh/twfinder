// Label component interface and implementation.

package server

// Label interface defines a component which wraps a text into a component.
//
// Default style class: "gui-Label"
type Label interface {
	// Label is a component.
	Comp

	// Label has text.
	HasText
}

// Label implementation
type labelImpl struct {
	compImpl    // Component implementation
	hasTextImpl // Has text implementation
}

// NewLabel creates a new Label.
func NewLabel(text string) Label {
	c := &labelImpl{newCompImpl(nil), newHasTextImpl(text)}
	c.Style().AddClass("gui-Label")
	return c
}

func (c *labelImpl) Render(w Writer) {
	w.Write(strSpanOp)
	c.renderAttrsAndStyle(w)
	c.renderEHandlers(w)
	w.Write(strGT)

	c.renderText(w)

	w.Write(strSpanCl)
}
