package gui

import (
  "image/color"

  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/theme"
)

// whiteTextTheme wraps an existing theme and forces all text to white.
type whiteTextTheme struct {
  fyne.Theme
}

// NewWhiteTextTheme returns a Theme that renders all text (and placeholders)
// in white, delegating every other color lookup to the provided base theme.
func NewWhiteTextTheme(base fyne.Theme) fyne.Theme {
  return &whiteTextTheme{Theme: base}
}


// greyedTextTheme wraps an existing theme and forces all text to a very light grey.
type greyedTextTheme struct {
  fyne.Theme
}


func NewGreyedTextTheme(base fyne.Theme) fyne.Theme {
  return &greyedTextTheme{Theme: base}
}



func (g *greyedTextTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
  switch name {
  case theme.ColorNameForeground, theme.ColorNamePlaceHolder:
    // slightly greyed‐out white
    //return color.NRGBA{R: 0xEE, G: 0xEE, B: 0xEE, A: 0xFF}
    return color.NRGBA{R: 0xDA, G: 0xDA, B: 0xDA, A: 0xFF}
  default:
    return g.Theme.Color(name, variant)
  }
}
