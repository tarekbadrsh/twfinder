// TabPanel component interface and implementation.

package server

// TabBar interface defines the tab bar for selecting the visible
// component of a TabPanel.
//
// Note: Removing a tab component through the tab bar also
// removes the content component from the tab panel of the tab bar.
//
// Default style classes: "gui-TabBar", "gui-TabBar-Top", "gui-TabBar-Bottom",
// "gui-TabBar-Left", "gui-TabBar-Right", "gui-TabBar-NotSelected",
// "gui-TabBar-Selected"
type TabBar interface {
	// TabBar is a PanelView.
	PanelView
}

// TabBar implementation.
type tabBarImpl struct {
	panelImpl // panel implementation
}

// newTabBarImpl creates a new tabBarImpl.
func newTabBarImpl() *tabBarImpl {
	c := &tabBarImpl{newPanelImpl()}
	return c
}

func (c *tabBarImpl) Remove(c2 Comp) bool {
	i := c.CompIdx(c2)
	if i < 0 {
		return false
	}

	// Removing a tab component also needs removing the
	// associated content component. Call parent's (TabPanel) Remove()
	// method which takes care of everything:
	if parent := c.parent; parent != nil {
		if tabPanel, isTabPanel := parent.(TabPanel); isTabPanel {
			return tabPanel.Remove(tabPanel.CompAt(i))
		}
	}

	return c.panelImpl.Remove(c2)
}

// TabBarPlacement is the Tab bar placement type.
type TabBarPlacement int

// Tab bar placements.
const (
	TbPlacementTop    TabBarPlacement = iota // Tab bar placement to Top
	TbPlacementBottom                        // Tab bar placement to Bottom
	TbPlacementLeft                          // Tab bar placement to Left
	TbPlacementRight                         // Tab bar placement to Right
)

// TabPanel interface defines a PanelView which has multiple child components
// but only one is visible at a time. The visible child can be visually selected
// using an internal TabBar component.
//
// Both the tab panel and its internal tab bar component are PanelViews.
// This gives high layout configuration possibilities.
// Usually you only need to set the tab bar placement with the SetTabBarPlacement()
// method which also sets other reasonable internal layout defaults.
// But you have many other options to override the layout settings.
// If the content component is bigger than the tab bar, you can set the tab bar
// horizontal and the vertical alignment, e.g. with the TabBar().SetAlignment() method.
// To apply cell formatting to individual content components, you can simply use the
// CellFmt() method. If the tab bar is bigger than the content component, you can set
// the content alignment, e.g. with the SetAlignment() method. You can also set different
// alignments for individual tab components using TabBar().CellFmt(). You can also set
// other cell formatting applied to the tab bar using TabBarFmt() method.
//
// You can register ETypeStateChange event handlers which will be called when the user
// changes tab selection by clicking on a tab. The event source will be the tab panel.
// The event will have a parent event whose source will be the clicked tab and will
// contain the mouse coordinates.
//
// Default style classes: "gui-TabPanel", "gui-TabPanel-Content"
type TabPanel interface {
	// TabPanel is a Container.
	PanelView

	// TabBar returns the tab bar.
	TabBar() TabBar

	// TabBarPlacement returns the tab bar placement.
	TabBarPlacement() TabBarPlacement

	// SetTabBarPlacement sets tab bar placement.
	// Also sets the alignment of the tab bar according
	// to the tab bar placement.
	SetTabBarPlacement(tabBarPlacement TabBarPlacement)

	// TabBarFmt returns the cell formatter of the tab bar.
	TabBarFmt() CellFmt

	// Add adds a new tab (component) and an associated (content) component
	// to the tab panel.
	Add(tab, content Comp)

	// Add adds a new tab (string) and an associated (content) component
	// to the tab panel.
	// This is a shorthand for
	// 		Add(NewLabel(tab), content)
	AddString(tab string, content Comp)

	// Selected returns the selected tab idx.
	// Returns -1 if no tab is selected.
	Selected() int

	// PrevSelected returns the previous selected tab idx.
	// Returns -1 if no tab was previously selected.
	PrevSelected() int

	// SetSelected sets the selected tab idx.
	// If idx < 0, no tabs will be selected.
	// If idx > CompsCount(), this is a no-op.
	SetSelected(idx int)
}

// TabPanel implementation.
type tabPanelImpl struct {
	panelImpl // panel implementation: TabPanel is a Panel, but only PanelView's methods are exported.

	tabBarImpl      *tabBarImpl     // Tab bar implementation
	tabBarPlacement TabBarPlacement // Tab bar placement
	tabBarFmt       *cellFmtImpl    // Tab bar cell formatter

	selected     int // The selected tab idx
	prevSelected int // Previous selected tab idx
}

// NewTabPanel creates a new TabPanel.
// Default tab bar placement is TbPlacementTop,
// default horizontal alignment is HADefault,
// default vertical alignment is VADefault.
func NewTabPanel() TabPanel {
	c := &tabPanelImpl{panelImpl: newPanelImpl(), tabBarImpl: newTabBarImpl(), tabBarFmt: newCellFmtImpl(), selected: -1, prevSelected: -1}
	c.tabBarFmt.Style().AddClass("gui-TabBar")
	c.tabBarImpl.setParent(c)
	c.SetTabBarPlacement(TbPlacementTop)
	c.tabBarFmt.SetAlign(HALeft, VATop)
	c.Style().AddClass("gui-TabPanel")
	return c
}

func (c *tabPanelImpl) Remove(c2 Comp) bool {
	i := c.CompIdx(c2)
	if i < 0 {
		// Try the tab bar:
		i = c.tabBarImpl.CompIdx(c2)
		if i < 0 {
			return false
		}

		// It's a tab component
		return c.Remove(c.panelImpl.CompAt(i))
	}

	// It's a content component
	c.tabBarImpl.panelImpl.Remove(c.tabBarImpl.CompAt(i))
	c.panelImpl.Remove(c2)

	// Update the previous selected
	if c.prevSelected >= 0 {
		if i < c.prevSelected {
			c.prevSelected-- // Keep the same previous selected by decreasing its index by 1
		} else if i == c.prevSelected { // Previous selected tab was removed...
			c.prevSelected = -1
		}
	}

	// Update the current selected
	if i < c.selected {
		c.selected-- // Keep the same tab selected by decreasing its index by 1
	} else if i == c.selected { // Selected tab was removed...
		// Store previous selected as it will be implicitly changed here
		prevSelected := c.prevSelected
		if i < c.CompsCount() {
			c.SetSelected(i) // There is next tab, select it
		} else if i > 0 { // Last was selected and removed but there are previous tabs...
			c.SetSelected(i - 1) // ...select the "new" last one
		} else { // Last was selected and removed and no previous tabs...
			c.SetSelected(-1) // No tabs remained.
		}
		// Restore previous selected
		c.prevSelected = prevSelected
	}

	return true
}

func (c *tabPanelImpl) ByID(id ID) Comp {
	// panelImpl.ById() also checks our own id first
	c2 := c.panelImpl.ByID(id)
	if c2 != nil {
		return c2
	}

	c2 = c.tabBarImpl.ByID(id)
	if c2 != nil {
		return c2
	}

	return nil
}

func (c *tabPanelImpl) Clear() {
	c.tabBarImpl.Clear()
	c.panelImpl.Clear()

	c.SetSelected(-1)
}

func (c *tabPanelImpl) TabBar() TabBar {
	return c.tabBarImpl
}

func (c *tabPanelImpl) TabBarPlacement() TabBarPlacement {
	return c.tabBarPlacement
}

func (c *tabPanelImpl) SetTabBarPlacement(tabBarPlacement TabBarPlacement) {
	style := c.tabBarFmt.Style()

	// Remove old style class
	switch c.tabBarPlacement {
	case TbPlacementTop:
		style.RemoveClass("gui-TabBar-Top")
	case TbPlacementBottom:
		style.RemoveClass("gui-TabBar-Bottom")
	case TbPlacementLeft:
		style.RemoveClass("gui-TabBar-Left")
	case TbPlacementRight:
		style.RemoveClass("gui-TabBar-Right")
	}

	c.tabBarPlacement = tabBarPlacement

	switch tabBarPlacement {
	case TbPlacementTop:
		c.tabBarImpl.SetLayout(LayoutHorizontal)
		c.tabBarImpl.SetAlign(HALeft, VABottom)
		style.AddClass("gui-TabBar-Top")
	case TbPlacementBottom:
		c.tabBarImpl.SetLayout(LayoutHorizontal)
		c.tabBarImpl.SetAlign(HALeft, VATop)
		style.AddClass("gui-TabBar-Bottom")
	case TbPlacementLeft:
		c.tabBarImpl.SetLayout(LayoutVertical)
		c.tabBarImpl.SetAlign(HARight, VATop)
		style.AddClass("gui-TabBar-Left")
	case TbPlacementRight:
		c.tabBarImpl.SetLayout(LayoutVertical)
		c.tabBarImpl.SetAlign(HALeft, VATop)
		style.AddClass("gui-TabBar-Right")
	}
}

func (c *tabPanelImpl) TabBarFmt() CellFmt {
	return c.tabBarFmt
}

func (c *tabPanelImpl) Add(tab, content Comp) {
	c.tabBarImpl.Add(tab)
	c.panelImpl.Add(content)
	c.tabBarImpl.CellFmt(tab).Style().AddClass("gui-TabBar-NotSelected")
	c.CellFmt(content).Style().AddClass("gui-TabPanel-Content")

	if c.CompsCount() == 1 {
		c.SetSelected(0)
	}

	// TODO would be nice to remove this internal handler func when the tab is removed!
	tab.AddEHandlerFunc(func(e Event) {
		c.SetSelected(c.CompIdx(content))
		e.MarkDirty(c)
		if c.handlers[ETypeStateChange] != nil {
			c.dispatchEvent(e.forkEvent(ETypeStateChange, c))
		}
	}, ETypeClick)
}

func (c *tabPanelImpl) AddString(tab string, content Comp) {
	tabc := NewLabel(tab)
	tabc.Style().SetDisplay(DisplayBlock) // Display: block - so the whole cell of the tab is clickable
	c.Add(tabc, content)
}

func (c *tabPanelImpl) Selected() int {
	return c.selected
}

func (c *tabPanelImpl) PrevSelected() int {
	return c.prevSelected
}

func (c *tabPanelImpl) SetSelected(idx int) {
	if idx >= c.CompsCount() {
		return
	}

	if c.selected >= 0 {
		// Deselect current selected
		style := c.tabBarImpl.CellFmt(c.tabBarImpl.CompAt(c.selected)).Style()
		style.RemoveClass("gui-TabBar-Selected")
		style.AddClass("gui-TabBar-NotSelected")
	}

	c.prevSelected = c.selected
	c.selected = idx

	if c.selected >= 0 {
		// Select new selected
		style := c.tabBarImpl.CellFmt(c.tabBarImpl.CompAt(c.selected)).Style()
		style.RemoveClass("gui-TabBar-NotSelected")
		style.AddClass("gui-TabBar-Selected")
	}
}

func (c *tabPanelImpl) Render(w Writer) {
	w.Write(strTableOp)
	c.renderAttrsAndStyle(w)
	c.renderEHandlers(w)
	w.Write(strGT)

	switch c.tabBarPlacement {
	case TbPlacementTop:
		w.Write(strTR)
		c.tabBarFmt.render(strTDOp, w)
		c.tabBarImpl.Render(w)
		c.renderTr(w)
		c.renderContent(w)
	case TbPlacementBottom:
		c.renderTr(w)
		c.renderContent(w)
		w.Write(strTR)
		c.tabBarFmt.render(strTDOp, w)
		c.tabBarImpl.Render(w)
	case TbPlacementLeft:
		c.renderTr(w)
		c.tabBarFmt.render(strTDOp, w)
		c.tabBarImpl.Render(w)
		c.renderContent(w)
	case TbPlacementRight:
		c.renderTr(w)
		c.renderContent(w)
		c.tabBarFmt.render(strTDOp, w)
		c.tabBarImpl.Render(w)
	}

	w.Write(strTableCl)
}

// renderContent renders the selected content component.
func (c *tabPanelImpl) renderContent(w Writer) {
	// Render only the selected content component
	if c.selected >= 0 {
		c2 := c.comps[c.selected]
		c.renderTd(c2, w)
		c2.Render(w)
	} else {
		w.Write(strTD)
	}
}
