package main

import (
	"flag"
	"fmt"
	"sync"

	hook "github.com/robotn/gohook"
	"github.com/wisepythagoras/autoclicker/core"
)

func main() {
	var startCombinationStr string
	var endCombinationStr string
	intervalPtr := flag.Int("interval", 100, "The clicking interval in ms")
	flag.StringVar(&startCombinationStr, "start", "s+ctrl+alt", "The start key combination")
	flag.StringVar(&endCombinationStr, "end", "q+ctrl+alt", "The end key combination")
	flag.Parse()

	interval := *intervalPtr

	fmt.Printf("Interval %d ms\n", interval)

	session := core.Session{Interval: interval}
	var wg sync.WaitGroup

	go session.Init(startCombinationStr, endCombinationStr)
	go session.AutoclickTimer()

	s := hook.Start()
	<-hook.Process(s)

	wg.Wait()
}
