package layout

// FocusStack manages which model is currently capturing input
// Simple stack-based approach: only the top model receives input
type FocusStack struct {
	stack []*GenericLayout
}

var globalFocusStack = &FocusStack{}

func (f *FocusStack) Push(layout *GenericLayout) {
	f.stack = append(f.stack, layout)
}

func (f *FocusStack) Pop() *GenericLayout {
	if len(f.stack) <= 1 {
		return nil // Don't pop the root
	}
	last := f.stack[len(f.stack)-1]
	f.stack = f.stack[:len(f.stack)-1]
	return last
}

func (f *FocusStack) Current() *GenericLayout {
	if len(f.stack) == 0 {
		return nil
	}
	return f.stack[len(f.stack)-1]
}

func (f *FocusStack) IsCurrent(layout *GenericLayout) bool {
	return f.Current() == layout
}

func (f *FocusStack) Depth() int {
	return len(f.stack)
}

// Global functions for convenience
func pushFocus(layout *GenericLayout) {
	globalFocusStack.Push(layout)
}

func popFocus() *GenericLayout {
	return globalFocusStack.Pop()
}

func currentFocus() *GenericLayout {
	return globalFocusStack.Current()
}

func isFocusCurrent(layout *GenericLayout) bool {
	return globalFocusStack.IsCurrent(layout)
}
