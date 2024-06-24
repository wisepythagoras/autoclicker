package main

import (
	"flag"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

var signals chan time.Time
var cancelSignals chan int

func startClicker(interval int) {
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

func stringToKeyCombination(str string) []string {
	return strings.Split(str, "+")
}

func main() {
	var startCombinationStr string
	var endCombinationStr string
	intervalPtr := flag.Int("interval", 100, "The clicking interval in ms")
	flag.StringVar(&startCombinationStr, "start", "s+ctrl+alt", "The start key combination")
	flag.StringVar(&endCombinationStr, "end", "q+ctrl+alt", "The end key combination")
	flag.Parse()

	interval := *intervalPtr
	startCombination := stringToKeyCombination(startCombinationStr)
	endCombination := stringToKeyCombination(endCombinationStr)

	signals = make(chan time.Time)
	cancelSignals = make(chan int)

	fmt.Printf("Interval %d ms\n", interval)

	var wg sync.WaitGroup

	go func() {
		var startTime *time.Time = nil

		startHook := func(e hook.Event) {
			if startTime != nil {
				return
			}

			fmt.Println("Start autoclicking")
			t := time.Now()
			startTime = &t
			signals <- *startTime
		}

		endHook := func(e hook.Event) {
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
		}

		if endCombinationStr == startCombinationStr {
			hook.Register(hook.KeyDown, startCombination, func(e hook.Event) {
				if startTime == nil {
					startHook(e)
				} else {
					endHook(e)
				}
			})
			return
		}

		hook.Register(hook.KeyDown, endCombination, endHook)
		hook.Register(hook.KeyDown, startCombination, startHook)

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
