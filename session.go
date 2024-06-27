package main

import (
	"fmt"
	"time"

	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

type Session struct {
	Interval      int
	startTime     *time.Time
	signals       chan time.Time
	cancelSignals chan int
}

func (s *Session) startHook(e hook.Event) {
	if s.startTime != nil {
		return
	}

	fmt.Println("Start autoclicking")
	t := time.Now()
	s.startTime = &t
	s.signals <- *s.startTime
}

func (s *Session) endHook(e hook.Event) {
	if s.startTime == nil {
		fmt.Println("Autoclicking not running")
		return
	}

	fmt.Println("Stop autoclicking")

	diff := time.Now().Sub(*s.startTime)

	// Cancel here.
	s.cancelSignals <- 1
	s.startTime = nil

	fmt.Println("Clicked for", diff)
}

func (s *Session) Init(startCombinationStr, endCombinationStr string) {
	s.signals = make(chan time.Time)
	s.cancelSignals = make(chan int)

	startCombination := stringToKeyCombination(startCombinationStr)
	endCombination := stringToKeyCombination(endCombinationStr)

	if endCombinationStr == startCombinationStr {
		hook.Register(hook.KeyDown, startCombination, func(e hook.Event) {
			if s.startTime == nil {
				s.startHook(e)
			} else {
				s.endHook(e)
			}
		})
		return
	}

	hook.Register(hook.KeyDown, endCombination, s.endHook)
	hook.Register(hook.KeyDown, startCombination, s.startHook)
}

func (s *Session) Destroy() {
	hook.End()
}

func (s *Session) startClicker() {
	for {
		select {
		case <-s.cancelSignals:
			fmt.Println("End clicking")
			return
		default:
			robotgo.Click()
			time.Sleep(time.Millisecond * time.Duration(s.Interval))
		}
	}
}

func (s *Session) AutoclickTimer() {
	for {
		fmt.Println("Autoclick timer")
		<-s.signals
		s.startClicker()
	}
}
