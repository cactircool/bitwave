package models

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type GenericPane struct {
	model    tea.Model
	Viewport viewport.Model
	style    lipgloss.Style
	weight   float32 // layout weight
	gap int
}

func NewPane(model tea.Model, weight float32, style lipgloss.Style) Pane {
	return &GenericPane{
		model:    model,
		Viewport: viewport.New(0, 0),
		style:    style,
		weight:   weight,
		gap: 0,
	}
}

func NewPaneWithGap(model tea.Model, weight float32, style lipgloss.Style, gap int) Pane {
	return &GenericPane{
		model:    model,
		Viewport: viewport.New(0, 0),
		style:    style,
		weight:   weight,
		gap: gap,
	}
}

func (p *GenericPane) Style() lipgloss.Style {
	return p.style
}

func (p *GenericPane) SetStyle(style lipgloss.Style) {
	p.style = style
}

func (p *GenericPane) Model() tea.Model {
	return p.Model()
}

func (p *GenericPane) Gap() int {
	return p.gap
}

func (p *GenericPane) Weight() float32 {
	return p.weight
}

func (p *GenericPane) Resize(outerW, outerH int) {
	chromeW, chromeH := chromeSize(p.Style())

	contentW := outerW - chromeW
	contentH := outerH - chromeH

	if contentW < 0 {
		contentW = 0
	}
	if contentH < 0 {
		contentH = 0
	}

	p.Viewport.Width = contentW
	p.Viewport.Height = contentH
}

func (p *GenericPane) Update(w, h int, msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	if v, ok := p.Model().(*NestedModel); ok {
		v.recalculate(w, h, p.Style())
	}
	p.model, cmd = p.Model().Update(msg)
	p.Viewport.SetContent(p.Model().View())
	return cmd
}

func (p GenericPane) View() string {
	return p.Style().Render(p.Viewport.View())
}

func chromeSize(style lipgloss.Style) (w, h int) {
	bw, bh := style.GetHorizontalBorderSize(), style.GetVerticalBorderSize()
	pw, ph := style.GetHorizontalPadding(), style.GetVerticalPadding()
	return bw + pw, bh + ph
}
