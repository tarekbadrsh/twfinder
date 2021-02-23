// Button component interface and implementation.

package server

// Button interface defines a clickable button.
//
// Suggested event type to handle actions: ETypeClick
//
// Default style class: "gui-Button"
type Button interface {
	// Button is a component.
	Comp

	// Button has text.
	HasText

	// Button can be enabled/disabled.
	HasEnabled
}

// Button implementation.
type buttonImpl struct {
	compImpl       // Component implementation
	hasTextImpl    // Has text implementation
	hasEnabledImpl // Has enabled implementation
}

// NewButton creates a new Button.
func NewButton(text string) Button {
	c := newButtonImpl(nil, text)
	c.Style().AddClass("gui-Button")
	return &c
}

// newButtonImpl creates a new buttonImpl.
func newButtonImpl(valueProviderJs []byte, text string) buttonImpl {
	return buttonImpl{newCompImpl(valueProviderJs), newHasTextImpl(text), newHasEnabledImpl()}
}

var (
	strButtonOp = []byte(`<button type="button"`) // `<button type="button"`
	strButtonCl = []byte("</button>")             // "</button>"
)

func (c *buttonImpl) Render(w Writer) {
	w.Write(strButtonOp)
	c.renderAttrsAndStyle(w)
	c.renderEHandlers(w)
	c.renderEnabled(w)
	w.Write(strGT)

	c.renderText(w)

	w.Write(strButtonCl)
}
