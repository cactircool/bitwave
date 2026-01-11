package layout

import (
	"strings"

	"github.com/cactircool/bitwave/bindings"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type GenericLayout struct {
	direction Direction
	children  []LayoutChild
	width     int
	height    int
	focused   int
}

type Direction int

const (
	Horizontal Direction = iota
	Vertical
)

type SizeMode int

const (
	Weighted SizeMode = iota
	Static
)

type LayoutChild struct {
	model        SizedModel
	sizeMode     SizeMode
	weight       float64
	size         int
	gap          int
	baseStyle    lipgloss.Style // Original style
	currentStyle lipgloss.Style // Current style (may include focus styling)
}

func NewLayout(direction Direction) *GenericLayout {
	return &GenericLayout{
		direction: direction,
		focused:   -1, // No focus initially
	}
}

func (l *GenericLayout) Add(model SizedModel, weight float64, style lipgloss.Style, gap int) {
	l.children = append(l.children, LayoutChild{
		model:        model,
		sizeMode:     Weighted,
		weight:       weight,
		gap:          gap,
		baseStyle:    style,
		currentStyle: style,
	})
}

func (l *GenericLayout) AddStatic(model SizedModel, size int, style lipgloss.Style, gap int) {
	l.children = append(l.children, LayoutChild{
		model:        model,
		sizeMode:     Static,
		size:         size,
		gap:          gap,
		baseStyle:    style,
		currentStyle: style,
	})
}

func (l *GenericLayout) SetSize(width, height int) {
	oldWidth, oldHeight := l.width, l.height
	l.width = width
	l.height = height
	if oldWidth != width || oldHeight != height {
		l.layoutChildren()
	}
}

func (l *GenericLayout) GetFocusState() FocusState {
	// A layout is focusable if it has any focusable children
	for _, child := range l.children {
		if child.model.GetFocusState() != NotFocusable {
			return Focusable
		}
	}
	return NotFocusable
}

func (l *GenericLayout) OnFocus(baseStyle lipgloss.Style) (lipgloss.Style, tea.Cmd) {
	// When layout gains focus, DON'T focus children yet
	// Children will be focused when we're pushed onto the focus stack
	return baseStyle.Border(lipgloss.ThickBorder()), nil
}

// func (l *GenericLayout) OnBlur() {
// 	// Blur the currently focused child
// 	if l.focused >= 0 && l.focused < len(l.children) {
// 		l.children[l.focused].model.OnBlur()
// 		l.children[l.focused].currentStyle = l.children[l.focused].baseStyle
// 	}
// }

// func (l *GenericLayout) OnBlur() {
// 	// Blur the currently focused child
// 	if l.focused >= 0 && l.focused < len(l.children) {
// 		l.children[l.focused].model.OnBlur()
// 		l.children[l.focused].currentStyle = l.children[l.focused].baseStyle
// 		// Re-layout since style changed
// 		l.layoutChildren()
// 	}
// }

func (l *GenericLayout) OnBlur() {
	// Blur the currently focused child
	if l.focused >= 0 && l.focused < len(l.children) {
		l.children[l.focused].model.OnBlur()
		oldStyle := l.children[l.focused].currentStyle
		l.children[l.focused].currentStyle = l.children[l.focused].baseStyle

		// Re-layout if frame size changed
		if oldStyle.GetHorizontalFrameSize() != l.children[l.focused].baseStyle.GetHorizontalFrameSize() ||
		   oldStyle.GetVerticalFrameSize() != l.children[l.focused].baseStyle.GetVerticalFrameSize() {
			l.layoutChildren()
		}
	}
}

// focusFirst focuses the first focusable child
func (l *GenericLayout) focusFirst() tea.Cmd {
	for i := range l.children {
		if l.children[i].model.GetFocusState() != NotFocusable {
			return l.focusChild(i)
		}
	}
	return nil
}

// focusChild focuses a specific child by index
// func (l *GenericLayout) focusChild(index int) tea.Cmd {
// 	if index < 0 || index >= len(l.children) {
// 		return nil
// 	}

// 	// Blur previous if different from new target
// 	if l.focused != index && l.focused >= 0 && l.focused < len(l.children) {
// 		l.children[l.focused].model.OnBlur()
// 		l.children[l.focused].currentStyle = l.children[l.focused].baseStyle
// 	}

// 	// Focus new
// 	l.focused = index
// 	style, cmd := l.children[index].model.OnFocus(l.children[index].baseStyle)
// 	l.children[index].currentStyle = style
// 	return cmd
// }

// func (l *GenericLayout) focusChild(index int) tea.Cmd {
// 	if index < 0 || index >= len(l.children) {
// 		return nil
// 	}

// 	// Blur previous if different from new target
// 	if l.focused != index && l.focused >= 0 && l.focused < len(l.children) {
// 		l.children[l.focused].model.OnBlur()
// 		l.children[l.focused].currentStyle = l.children[l.focused].baseStyle
// 	}

// 	// Focus new
// 	l.focused = index
// 	style, cmd := l.children[index].model.OnFocus(l.children[index].baseStyle)
// 	l.children[index].currentStyle = style

// 	// Re-layout since style may have changed frame sizes
// 	l.layoutChildren()

// 	return cmd
// }

func (l *GenericLayout) focusChild(index int) tea.Cmd {
	if index < 0 || index >= len(l.children) {
		return nil
	}

	// Blur previous if different from new target
	if l.focused != index && l.focused >= 0 && l.focused < len(l.children) {
		l.children[l.focused].model.OnBlur()
		l.children[l.focused].currentStyle = l.children[l.focused].baseStyle
	}

	// Focus new
	l.focused = index
	oldStyle := l.children[index].currentStyle
	style, cmd := l.children[index].model.OnFocus(l.children[index].baseStyle)
	l.children[index].currentStyle = style

	// Re-layout if frame size changed
	if oldStyle.GetHorizontalFrameSize() != style.GetHorizontalFrameSize() ||
	   oldStyle.GetVerticalFrameSize() != style.GetVerticalFrameSize() {
		l.layoutChildren()
	}

	return cmd
}

// cycleForward moves focus to the next focusable child
func (l *GenericLayout) cycleForward() tea.Cmd {
	if len(l.children) == 0 {
		return nil
	}

	start := l.focused
	if start < 0 {
		start = 0
	}

	next := start
	for {
		next = (next + 1) % len(l.children)
		if l.children[next].model.GetFocusState() != NotFocusable {
			return l.focusChild(next)
		}
		if next == start {
			break // Wrapped around
		}
	}
	return nil
}

// cycleBackward moves focus to the previous focusable child
func (l *GenericLayout) cycleBackward() tea.Cmd {
	if len(l.children) == 0 {
		return nil
	}

	start := l.focused
	if start < 0 {
		start = 0
	}

	prev := start
	for {
		prev = (prev - 1 + len(l.children)) % len(l.children)
		if l.children[prev].model.GetFocusState() != NotFocusable {
			return l.focusChild(prev)
		}
		if prev == start {
			break // Wrapped around
		}
	}
	return nil
}

// func (l *GenericLayout) layoutChildren() {
// 	if len(l.children) == 0 {
// 		return
// 	}

// 	// Calculate what we need to subtract from available space
// 	totalGap := 0
// 	totalStaticSize := 0
// 	totalWeight := 0.0

// 	for i := range l.children {
// 		child := &l.children[i]

// 		// Gaps between children (not after last child)
// 		if i < len(l.children)-1 {
// 			totalGap += child.gap
// 		}

// 		if child.sizeMode == Static {
// 			// Static children: their size PLUS their chrome
// 			totalStaticSize += child.size
// 			if l.direction == Horizontal {
// 				totalStaticSize += child.currentStyle.GetHorizontalFrameSize()
// 			} else {
// 				totalStaticSize += child.currentStyle.GetVerticalFrameSize()
// 			}
// 		} else {
// 			// Weighted children: just accumulate weight
// 			totalWeight += child.weight
// 		}
// 	}

// 	// Available space is total minus gaps and static children (including their chrome)
// 	var totalSpace int
// 	if l.direction == Horizontal {
// 		totalSpace = l.width
// 	} else {
// 		totalSpace = l.height
// 	}

// 	availableForWeighted := totalSpace - totalGap - totalStaticSize
// 	fmt.Fprintf(os.Stderr, "avilableForWeighted = %d = %d - %d - %d\n", availableForWeighted, totalSpace, totalGap, totalStaticSize)
// 	if availableForWeighted < 0 {
// 		availableForWeighted = 0
// 	}

// 	// Now we need to subtract the chrome for weighted children too
// 	weightedChrome := 0
// 	for _, child := range l.children {
// 		if child.sizeMode == Weighted {
// 			if l.direction == Horizontal {
// 				weightedChrome += child.currentStyle.GetHorizontalFrameSize()
// 			} else {
// 				weightedChrome += child.currentStyle.GetVerticalFrameSize()
// 			}
// 		}
// 	}

// 	availableForWeighted -= weightedChrome
// 	if availableForWeighted < 0 {
// 		availableForWeighted = 0
// 	}

// 	// Assign sizes
// 	for i := range l.children {
// 		child := &l.children[i]

// 		var innerWidth, innerHeight int

// 		if child.sizeMode == Static {
// 			// Static: use their fixed size as the inner dimension
// 			if l.direction == Horizontal {
// 				innerWidth = child.size
// 				innerHeight = l.height - child.currentStyle.GetVerticalFrameSize()
// 			} else {
// 				innerWidth = l.width - child.currentStyle.GetHorizontalFrameSize()
// 				innerHeight = child.size
// 			}
// 		} else {
// 			// Weighted: distribute available space by weight
// 			var innerSize int
// 			if totalWeight > 0 {
// 				ratio := child.weight / totalWeight
// 				innerSize = int(float64(availableForWeighted) * ratio)
// 			}

// 			if l.direction == Horizontal {
// 				innerWidth = innerSize
// 				innerHeight = l.height - child.currentStyle.GetVerticalFrameSize()
// 			} else {
// 				innerWidth = l.width - child.currentStyle.GetHorizontalFrameSize()
// 				innerHeight = innerSize
// 			}
// 		}

// 		if innerWidth < 0 {
// 			innerWidth = 0
// 		}
// 		if innerHeight < 0 {
// 			innerHeight = 0
// 		}

// 		child.model.SetSize(innerWidth, innerHeight)
// 	}
// }

func (l *GenericLayout) layoutChildren() {
	if len(l.children) == 0 {
		return
	}

	// Calculate total gaps
	totalGap := 0
	for i := range l.children {
		if i < len(l.children)-1 {
			totalGap += l.children[i].gap
		}
	}

	// Calculate space used by static children (including their chrome)
	totalStaticSize := 0
	for i := range l.children {
		child := &l.children[i]
		if child.sizeMode == Static {
			totalStaticSize += child.size
			if l.direction == Horizontal {
				totalStaticSize += child.currentStyle.GetHorizontalFrameSize()
			} else {
				totalStaticSize += child.currentStyle.GetVerticalFrameSize()
			}
		}
	}

	// Calculate total chrome for weighted children
	weightedChrome := 0
	totalWeight := 0.0
	for _, child := range l.children {
		if child.sizeMode == Weighted {
			totalWeight += child.weight
			if l.direction == Horizontal {
				weightedChrome += child.currentStyle.GetHorizontalFrameSize()
			} else {
				weightedChrome += child.currentStyle.GetVerticalFrameSize()
			}
		}
	}

	// Available space for weighted children's CONTENT
	var totalSpace int
	if l.direction == Horizontal {
		totalSpace = l.width
	} else {
		totalSpace = l.height
	}

	availableForWeighted := totalSpace - totalGap - totalStaticSize - weightedChrome
	if availableForWeighted < 0 {
		availableForWeighted = 0
	}


	// Assign sizes to children
	for i := range l.children {
		child := &l.children[i]

		var innerWidth, innerHeight int

		if child.sizeMode == Static {
			// Static: use their fixed size as the inner dimension
			if l.direction == Horizontal {
				innerWidth = child.size
				innerHeight = l.height - child.currentStyle.GetVerticalFrameSize()
			} else {
				innerWidth = l.width - child.currentStyle.GetHorizontalFrameSize()
				innerHeight = child.size
			}
		} else {
			// Weighted: distribute available space by weight
			var innerSize int
			if totalWeight > 0 {
				ratio := child.weight / totalWeight
				innerSize = int(float64(availableForWeighted) * ratio)
			}

			if l.direction == Horizontal {
				innerWidth = innerSize
				innerHeight = l.height - child.currentStyle.GetVerticalFrameSize()
			} else {
				innerWidth = l.width - child.currentStyle.GetHorizontalFrameSize()
				innerHeight = innerSize
			}
		}

		// Ensure non-negative dimensions
		if innerWidth < 0 {
			innerWidth = 0
		}
		if innerHeight < 0 {
			innerHeight = 0
		}

		child.model.SetSize(innerWidth, innerHeight)
	}
}

func (l *GenericLayout) Init() tea.Cmd {
	cmds := make([]tea.Cmd, 0, len(l.children))
	for _, child := range l.children {
		if cmd := child.model.Init(); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	return tea.Batch(cmds...)
}

func (l *GenericLayout) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Only handle navigation if we're the current focus
	if !isFocusCurrent(l) {
		// Forward to focused child
		if l.focused >= 0 && l.focused < len(l.children) {
			model, cmd := l.children[l.focused].model.Update(msg)
			l.children[l.focused].model = model.(SizedModel)
			return l, cmd
		}
		return l, nil
	}

	// We're current - handle our keys
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case bindings.QuitProgram:
			return l, tea.Quit

		case bindings.CycleFocusForward:
			return l, l.cycleForward()

		case bindings.CycleFocusBackward:
			return l, l.cycleBackward()

		case bindings.CycleEnter:
			// Enter the focused child if it's interactive or a layout
			if l.focused >= 0 && l.focused < len(l.children) {
				child := l.children[l.focused].model
				focusState := child.GetFocusState()

				// If it's a layout, dive into it
				if childLayout, ok := child.(*GenericLayout); ok && focusState != NotFocusable {
					pushFocus(childLayout)
					return l, childLayout.focusFirst()
				}

				// If it's interactive, let it handle enter
				if focusState == Interactive {
					model, cmd := child.Update(msg)
					l.children[l.focused].model = model.(SizedModel)
					return l, cmd
				}
			}
			return l, nil

		case bindings.CycleEscape:
			// Pop focus back to parent
			popped := popFocus()
			if popped != nil {
				// Blur the currently focused child before popping out
				if popped.focused >= 0 && popped.focused < len(popped.children) {
					popped.children[popped.focused].model.OnBlur()
					popped.children[popped.focused].currentStyle = popped.children[popped.focused].baseStyle
				}
				popped.layoutChildren()
			}
			return l, nil
		}
	}

	// Forward other messages to focused child
	if l.focused >= 0 && l.focused < len(l.children) {
		model, cmd := l.children[l.focused].model.Update(msg)
		l.children[l.focused].model = model.(SizedModel)
		return l, cmd
	}

	return l, nil
}

func (l *GenericLayout) View() string {
	if len(l.children) == 0 {
		return ""
	}

	views := make([]string, 0, len(l.children)*2)

	for i, child := range l.children {
		content := child.model.View()
		views = append(views, child.currentStyle.Render(content))

		if i < len(l.children)-1 && child.gap > 0 {
			if l.direction == Horizontal {
				views = append(views, strings.Repeat(" ", child.gap))
			} else {
				views = append(views, strings.Repeat("\n", child.gap))
			}
		}
	}

	if l.direction == Horizontal {
		return lipgloss.JoinHorizontal(lipgloss.Top, views...)
	}
	return lipgloss.JoinVertical(lipgloss.Left, views...)
}
