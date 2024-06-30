package main

import (
	"fmt"
	"strings"

	g "github.com/AllenDang/giu"
)

type Window struct {
	window   *g.MasterWindow
	KeyCombo []string
}

func main() {
	window := Window{
		window:   g.NewMasterWindow("Autoclicker", 400, 200, 0),
		KeyCombo: []string{"s", "ctrl", "alt"},
	}
	window.window.Run(window.loop)
}

func (w *Window) loop() {
	g.SingleWindow().Layout(
		g.Label(fmt.Sprintf("Combination: %s", strings.Join(w.KeyCombo, "+"))),
		g.Button("Register Hotkey").OnClick(func() {
			fmt.Println("Clicked")
		}),
	)
}
