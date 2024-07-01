package main

import (
	"fmt"
	"strings"

	g "github.com/AllenDang/giu"
	hook "github.com/robotn/gohook"
)

type Window struct {
	window    *g.MasterWindow
	KeyCombo  []string
	tempCombo []string
}

func main() {
	window := Window{
		window:    g.NewMasterWindow("Autoclicker", 400, 200, 0),
		KeyCombo:  []string{"s", "ctrl", "alt"},
		tempCombo: []string{},
	}
	window.Init()
	window.window.Run(window.loop)
}

func (w *Window) Init() {
	hook.Register(hook.KeyDown, []string{}, func(e hook.Event) {
		fmt.Println(e)
	})

	go func() {
		s := hook.Start()
		<-hook.Process(s)
	}()
}

func (w *Window) loop() {
	g.SingleWindow().Layout(
		g.Label(fmt.Sprintf("Combination: %s", strings.Join(w.KeyCombo, "+"))),
		g.Button("Register Hotkey").OnClick(func() {
			fmt.Println("Clicked")
		}),
	)
}
