// Link component interface and implementation.

package server

// Link interface defines a clickable link pointing to a URL.
// Links are usually used with a text, although Link is a
// container, and allows to set a child component
// which if set will also be a part of the clickable link.
//
// Default style class: "gui-Link"
type Link interface {
	// Link is a Container.
	Container

	// Link has text.
	HasText

	// Link has URL string.
	HasURL

	// Target returns the target of the link.
	Target() string

	// SetTarget sets the target of the link.
	// Tip: pass "_blank" if you want the URL to open in a new window
	// (this is the default).
	SetTarget(target string)

	// Comp returns the optional child component, if set.
	Comp() Comp

	// SetComp sets the only child component
	// (which can be a Container of course).
	SetComp(c Comp)
}

// Link implementation.
type linkImpl struct {
	compImpl    // Component implementation
	hasTextImpl // Has text implementation
	hasURLImpl  // Has text implementation

	comp Comp // Optional child component
}

// NewLink creates a new Link.
// By default links open in a new window (tab)
// because their target is set to "_blank".
func NewLink(text, url string) Link {
	c := &linkImpl{newCompImpl(nil), newHasTextImpl(text), newHasURLImpl(url), nil}
	c.SetTarget("_blank")
	c.Style().AddClass("gui-Link")
	return c
}

func (c *linkImpl) Remove(c2 Comp) bool {
	if c.comp == nil || !c.comp.Equals(c2) {
		return false
	}

	c2.setParent(nil)
	c.comp = nil

	return true
}

func (c *linkImpl) ByID(id ID) Comp {
	if c.id == id {
		return c
	}

	if c.comp != nil {
		if c.comp.ID() == id {
			return c.comp
		}
		if c2, isContainer := c.comp.(Container); isContainer {
			if c3 := c2.ByID(id); c3 != nil {
				return c3
			}
		}

	}

	return nil
}

func (c *linkImpl) Clear() {
	if c.comp != nil {
		c.comp.setParent(nil)
		c.comp = nil
	}
}

func (c *linkImpl) Target() string {
	return c.attrs["target"]
}

func (c *linkImpl) SetTarget(target string) {
	if len(target) == 0 {
		delete(c.attrs, "target")
	} else {
		c.attrs["target"] = target
	}
}

func (c *linkImpl) Comp() Comp {
	return c.comp
}

func (c *linkImpl) SetComp(c2 Comp) {
	c.comp = c2
}

var (
	strAOp = []byte("<a")   // "<a"
	strACL = []byte("</a>") // "</a>"
)

func (c *linkImpl) Render(w Writer) {
	w.Write(strAOp)
	c.renderURL("href", w)
	c.renderAttrsAndStyle(w)
	c.renderEHandlers(w)
	w.Write(strGT)

	c.renderText(w)

	if c.comp != nil {
		c.comp.Render(w)
	}

	w.Write(strACL)
}
