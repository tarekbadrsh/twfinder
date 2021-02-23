// Image component interface and implementation.

package server

// Image interface defines an image.
//
// Default style class: "gui-Image"
type Image interface {
	// Image is a component.
	Comp

	// Image has text which is its description (alternate text).
	HasText

	// Image has URL string.
	HasURL
}

// Image implementation
type imageImpl struct {
	compImpl    // Component implementation
	hasTextImpl // Has text implementation
	hasURLImpl  // Has text implementation
}

// NewImage creates a new Image.
// The text is used as the alternate text for the image.
func NewImage(text, url string) Image {
	c := &imageImpl{newCompImpl(nil), newHasTextImpl(text), newHasURLImpl(url)}
	c.Style().AddClass("gui-Image")
	return c
}

var (
	strImgOp = []byte("<img")   // "<img"
	strAlt   = []byte(` alt="`) // ` alt="`
	strImgCl = []byte(`">`)     // `">`
)

func (c *imageImpl) Render(w Writer) {
	w.Write(strImgOp)
	c.renderURL("src", w)
	c.renderAttrsAndStyle(w)
	c.renderEHandlers(w)
	w.Write(strAlt)
	c.renderText(w)
	w.Write(strImgCl)
}
