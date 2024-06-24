package main

import (
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

var signals chan time.Time
var cancelSignals chan int

func startClicker(interval int) {
	fmt.Println(time.Duration(interval))
	for {
		select {
		case <-cancelSignals:
			fmt.Println("End clicking")
			return
		default:
			robotgo.Click()
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func autoclickerTimer(interval int) {
	for {
		<-signals
		startClicker(interval)
	}
}

func main() {
	intervalPtr := flag.Int("interval", 100, "The clicking interval in ms")
	flag.Parse()

	interval := *intervalPtr

	signals = make(chan time.Time)
	cancelSignals = make(chan int)

	fmt.Printf("Interval %d ms\n", interval)

	var wg sync.WaitGroup

	go func() {
		var startTime *time.Time = nil

		hook.Register(hook.KeyDown, []string{"q", "ctrl", "alt"}, func(e hook.Event) {
			if startTime == nil {
				fmt.Println("Autoclicking not running")
				return
			}

			fmt.Println("Stop autoclicking")

			diff := time.Now().Sub(*startTime)

			// Cancel here.
			cancelSignals <- 1
			startTime = nil

			fmt.Println("Clicked for", diff)
		})

		hook.Register(hook.KeyDown, []string{"s", "ctrl", "alt"}, func(e hook.Event) {
			if startTime != nil {
				return
			}

			fmt.Println("Start autoclicking")
			t := time.Now()
			startTime = &t
			signals <- *startTime
		})

		hook.Register(hook.MouseHold, []string{}, func(e hook.Event) {
		})

		hook.Register(hook.MouseUp, []string{}, func(e hook.Event) {
		})
	}()

	go autoclickerTimer(interval)

	s := hook.Start()
	<-hook.Process(s)

	wg.Wait()
}
