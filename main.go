package main

import (
	"flag"
	"fmt"
	"sync"
	"time"

	hook "github.com/robotn/gohook"
)

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

		hook.Register(hook.KeyHold, []string{}, func(e hook.Event) {
			fmt.Println(e)
		})

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
