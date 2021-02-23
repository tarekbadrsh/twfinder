// Panel component interface and implementation.

package server

import (
	"bytes"
)

// Layout strategy type.
type Layout int

// Layout strategies.
const (
	LayoutNatural    Layout = iota // Natural layout: elements are displayed in their natural order.
	LayoutVertical                 // Vertical layout: elements are laid out vertically.
	LayoutHorizontal               // Horizontal layout: elements are laid out horizontally.
)

// PanelView interface defines a container which stores child components
// sequentially (one dimensional, associated with an index), and lays out
// its children in a row or column using TableView based on a layout strategy,
// but does not define the way how child components can be added.
//
// Default style class: "gui-Panel"
type PanelView interface {
	// PanelView is a TableView.
	TableView

	// Layout returns the layout strategy used to lay out components when rendering.
	Layout() Layout

	// SetLayout sets the layout strategy used to lay out components when rendering.
	SetLayout(layout Layout)

	// CompsCount returns the number of components added to the panel.
	CompsCount() int

	// CompAt returns the component at the specified index.
	// Returns nil if idx<0 or idx>=CompsCount().
	CompAt(idx int) Comp

	// CompIdx returns the index of the specified component in the panel.
	// -1 is returned if the component is not added to the panel.
	CompIdx(c Comp) int

	// CellFmt returns the cell formatter of the specified child component.
	// If the specified component is not a child, nil is returned.
	// Cell formatting has no effect if layout is LayoutNatural.
	CellFmt(c Comp) CellFmt
}

// Panel interface defines a container which stores child components
// associated with an index, and lays out its children based on a layout
// strategy.
// Default style class: "gui-Panel"
type Panel interface {
	// Panel is a PanelView.
	PanelView

	// Add adds a component to the panel.
	Add(c Comp)

	// Insert inserts a component at the specified index.
	// Returns true if the index was valid and the component is inserted
	// successfully, false otherwise. idx=CompsCount() is also allowed
	// in which case comp will be the last component.
	Insert(c Comp, idx int) bool

	// AddHSpace adds and returns a fixed-width horizontal space consumer.
	// Useful when layout is LayoutHorizontal.
	AddHSpace(width int) Comp

	// AddVSpace adds and returns a fixed-height vertical space consumer.
	// Useful when layout is LayoutVertical.
	AddVSpace(height int) Comp

	// AddSpace adds and returns a fixed-size space consumer.
	AddSpace(width, height int) Comp

	// AddHConsumer adds and returns a horizontal (free) space consumer.
	// Useful when layout is LayoutHorizontal.
	//
	// Tip: When adding a horizontal space consumer, you may set the
	// white space style attribute of other components in the the panel
	// to WhiteSpaceNowrap to avoid texts getting wrapped to multiple lines.
	AddHConsumer() Comp

	// AddVConsumer adds and returns a vertical (free) space consumer.
	// Useful when layout is LayoutVertical.
	AddVConsumer() Comp
}

// Panel implementation.
type panelImpl struct {
	tableViewImpl // TableView implementation

	layout   Layout              // Layout strategy
	comps    []Comp              // Components added to this panel
	cellFmts map[ID]*cellFmtImpl // Lazily initialized cell formatters of the child components
}

// NewPanel creates a new Panel.
// Default layout strategy is LayoutVertical,
// default horizontal alignment is HADefault,
// default vertical alignment is VADefault.
func NewPanel() Panel {
	c := newPanelImpl()
	c.Style().AddClass("gui-Panel")
	return &c
}

// NewNaturalPanel creates a new Panel initialized with
// LayoutNatural layout.
// Default horizontal alignment is HADefault,
// default vertical alignment is VADefault.
func NewNaturalPanel() Panel {
	p := NewPanel()
	p.SetLayout(LayoutNatural)
	return p
}

// NewHorizontalPanel creates a new Panel initialized with
// LayoutHorizontal layout.
// Default horizontal alignment is HADefault,
// default vertical alignment is VADefault.
func NewHorizontalPanel() Panel {
	p := NewPanel()
	p.SetLayout(LayoutHorizontal)
	return p
}

// NewVerticalPanel creates a new Panel initialized with
// LayoutVertical layout.
// Default horizontal alignment is HADefault,
// default vertical alignment is VADefault.
func NewVerticalPanel() Panel {
	return NewPanel()
}

// newPanelImpl creates a new panelImpl.
func newPanelImpl() panelImpl {
	return panelImpl{tableViewImpl: newTableViewImpl(), layout: LayoutVertical, comps: make([]Comp, 0, 2)}
}

func (c *panelImpl) Remove(c2 Comp) bool {
	i := c.CompIdx(c2)
	if i < 0 {
		return false
	}

	// Remove associated cell formatter
	if c.cellFmts != nil {
		delete(c.cellFmts, c2.ID())
	}

	c2.setParent(nil)
	// When removing, also reference must be cleared to allow the comp being gc'ed, also to prevent memory leak.
	oldComps := c.comps
	// Copy the part after the removable comp, backward by 1:
	c.comps = append(oldComps[:i], oldComps[i+1:]...)
	// Clear the reference that becomes unused:
	oldComps[len(oldComps)-1] = nil

	return true
}

func (c *panelImpl) ByID(id ID) Comp {
	if c.id == id {
		return c
	}

	for _, c2 := range c.comps {
		if c2.ID() == id {
			return c2
		}

		if c3, isContainer := c2.(Container); isContainer {
			if c4 := c3.ByID(id); c4 != nil {
				return c4
			}
		}
	}
	return nil
}

func (c *panelImpl) Clear() {
	// Clear cell formatters
	if c.cellFmts != nil {
		c.cellFmts = nil
	}

	for _, c2 := range c.comps {
		c2.setParent(nil)
	}
	c.comps = nil
}

func (c *panelImpl) Layout() Layout {
	return c.layout
}

func (c *panelImpl) SetLayout(layout Layout) {
	c.layout = layout
}

func (c *panelImpl) CompsCount() int {
	return len(c.comps)
}

func (c *panelImpl) CompAt(idx int) Comp {
	if idx < 0 || idx >= len(c.comps) {
		return nil
	}
	return c.comps[idx]
}

func (c *panelImpl) CompIdx(c2 Comp) int {
	for i, c3 := range c.comps {
		if c2.Equals(c3) {
			return i
		}
	}
	return -1
}

func (c *panelImpl) CellFmt(c2 Comp) CellFmt {
	if c.CompIdx(c2) < 0 {
		return nil
	}

	if c.cellFmts == nil {
		c.cellFmts = make(map[ID]*cellFmtImpl)
	}

	cf := c.cellFmts[c2.ID()]
	if cf == nil {
		cf = newCellFmtImpl()
		c.cellFmts[c2.ID()] = cf
	}
	return cf
}

func (c *panelImpl) Add(c2 Comp) {
	c2.makeOrphan()
	c.comps = append(c.comps, c2)
	c2.setParent(c)
}

func (c *panelImpl) Insert(c2 Comp, idx int) bool {
	if idx < 0 || idx > len(c.comps) {
		return false
	}

	c2.makeOrphan()

	// Make sure we have room for the extra component:
	c.comps = append(c.comps, nil)
	copy(c.comps[idx+1:], c.comps[idx:len(c.comps)-1])
	c.comps[idx] = c2

	c2.setParent(c)

	return true
}

func (c *panelImpl) AddHSpace(width int) Comp {
	l := NewLabel("")
	l.Style().SetDisplay(DisplayBlock).SetWidthPx(width)
	c.Add(l)
	return l
}

func (c *panelImpl) AddVSpace(height int) Comp {
	l := NewLabel("")
	l.Style().SetDisplay(DisplayBlock).SetHeightPx(height)
	c.Add(l)
	return l
}

func (c *panelImpl) AddSpace(width, height int) Comp {
	l := NewLabel("")
	l.Style().SetDisplay(DisplayBlock).SetSizePx(width, height)
	c.Add(l)
	return l
}

func (c *panelImpl) AddHConsumer() Comp {
	l := NewLabel("")
	c.Add(l)
	c.CellFmt(l).Style().SetFullWidth()
	return l
}

func (c *panelImpl) AddVConsumer() Comp {
	l := NewLabel("")
	c.Add(l)
	c.CellFmt(l).Style().SetFullHeight()
	return l
}

func (c *panelImpl) Render(w Writer) {
	switch c.layout {
	case LayoutNatural:
		c.layoutNatural(w)
	case LayoutHorizontal:
		c.layoutHorizontal(w)
	case LayoutVertical:
		c.layoutVertical(w)
	}
}

// layoutNatural renders the panel and the child components
// using the natural layout strategy.
func (c *panelImpl) layoutNatural(w Writer) {
	// No wrapper table but we still need a wrapper tag for attributes...
	w.Write(strSpanOp)
	c.renderAttrsAndStyle(w)
	c.renderEHandlers(w)
	w.Write(strGT)

	for _, c2 := range c.comps {
		c2.Render(w)
	}

	w.Write(strSpanCl)
}

// layoutHorizontal renders the panel and the child components
// using the horizontal layout strategy.
func (c *panelImpl) layoutHorizontal(w Writer) {
	w.Write(strTableOp)
	c.renderAttrsAndStyle(w)
	c.renderEHandlers(w)
	w.Write(strGT)

	c.renderTr(w)

	for _, c2 := range c.comps {
		c.renderTd(c2, w)
		c2.Render(w)
	}

	w.Write(strTableCl)
}

// layoutVertical renders the panel and the child components
// using the vertical layout strategy.
func (c *panelImpl) layoutVertical(w Writer) {
	w.Write(strTableOp)
	c.renderAttrsAndStyle(w)
	c.renderEHandlers(w)
	w.Write(strGT)

	// There is the same TR tag for each cell:
	trWriter := bytes.NewBuffer(nil)
	c.renderTr(NewWriter(trWriter))
	tr := trWriter.Bytes()

	for _, c2 := range c.comps {
		w.Write(tr)
		c.renderTd(c2, w)
		c2.Render(w)
	}

	w.Write(strTableCl)
}

// renderTd renders the formatted HTML TD tag for the specified child component.
func (c *panelImpl) renderTd(c2 Comp, w Writer) {
	if cf := c.cellFmts[c2.ID()]; cf == nil {
		w.Write(strTD)
	} else {
		cf.render(strTDOp, w)
	}
}
