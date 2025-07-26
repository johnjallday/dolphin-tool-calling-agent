package gui

import (
  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/container"
  "fyne.io/fyne/v2/widget"
)

type AddAgentForm struct {
  widget.BaseWidget

  NameEntry  *widget.Entry
  ModelEntry *widget.Entry
  Tools      *widget.CheckGroup

  onSubmit func(name, model string, tools []string)
}

func NewAddAgentForm(
  toolpacks []string,
  onSubmit func(name, model string, tools []string),
) *AddAgentForm {
  f := &AddAgentForm{
    NameEntry:  widget.NewEntry(),
    ModelEntry: widget.NewEntry(),
    Tools:      widget.NewCheckGroup(toolpacks, nil),
    onSubmit:   onSubmit,
  }
  f.NameEntry.SetPlaceHolder("Agent name")
  f.ModelEntry.SetPlaceHolder("Model (eg “gpt-4”)")
  f.ExtendBaseWidget(f)
  return f
}

func (f *AddAgentForm) CreateRenderer() fyne.WidgetRenderer {
  title := widget.NewLabelWithStyle("Add New Agent",
    fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

  // Wrap the CheckGroup in a Scroll so it doesn't grow forever
  toolsScroll := container.NewVScroll(f.Tools)
  toolsScroll.SetMinSize(fyne.NewSize(200, 100))  // tweak as you like

  form := widget.NewForm(
    &widget.FormItem{Text: "Name", Widget: f.NameEntry},
    &widget.FormItem{Text: "Model", Widget: f.ModelEntry},
    &widget.FormItem{Text: "Toolpacks", Widget: toolsScroll},
  )

  btn := widget.NewButton("Create Agent", func() {
    f.onSubmit(f.NameEntry.Text, f.ModelEntry.Text, f.Tools.Selected)
  })

  box := container.NewVBox(
    title,
    form,
    btn,
  )
  return widget.NewSimpleRenderer(box)
}
