package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-vgo/robotgo"
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
			time.Sleep(time.Millisecond * time.Duration(interval))
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
