package main

import (
	"fmt"
	"strings"
	"time"

	g "github.com/AllenDang/giu"
	hook "github.com/robotn/gohook"
)

type Window struct {
	window    *g.MasterWindow
	KeyCombo  []string
	tempCombo []string
	lastTap   *time.Time
	editMode  bool
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
		if !w.editMode {
			return
		}

		if w.lastTap == nil || w.lastTap.Sub(time.Now()).Milliseconds() > 1000 {
			w.tempCombo = make([]string, 0)
		}

		fmt.Println(string(e.Keychar), e.Rawcode)
	})

	go func() {
		s := hook.Start()
		<-hook.Process(s)
	}()
}

func (w *Window) loop() {
	buttonLabel := "Register Combination"

	if w.editMode {
		buttonLabel = "Save Combination"
	}

	g.SingleWindow().Layout(
		g.Label(fmt.Sprintf("Combination: %s", strings.Join(w.KeyCombo, "+"))),
		g.Button(buttonLabel).OnClick(func() {
			if w.editMode {
				w.editMode = false
				w.lastTap = nil
			} else {
				w.editMode = true
			}
		}),
	)
}
