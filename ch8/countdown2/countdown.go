// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 244.

// Countdown implements the countdown for a rocket launch.
package main

import (
	"fmt"
	"os"
	"time"
)

//!+

func main() {
	// ...create abort channel...

	//!-

	//!+abort
	abort := make(chan struct{})
	go func() {
		os.Stdin.Read(make([]byte, 1)) // read a single byte
		abort <- struct{}{}
	}()
	//!-abort

	//!+
	fmt.Println("Commencing countdown.  Press return to abort.")
	select {
	case <-time.After(10 * time.Second): //，阻塞，每隔10s会触发该通道接收操作;事件的超时操作
		// Do nothing.
	case <-abort:
		fmt.Println("Launch aborted!") //在10秒之类中断火箭发射
		return
	default:
	}
	launch() //火箭发射
}

//!-

func launch() {
	fmt.Println("Lift off!")
}
