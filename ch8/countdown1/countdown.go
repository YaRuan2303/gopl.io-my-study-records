// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 244.

// Countdown implements the countdown for a rocket launch.
package main

import (
	"fmt"
	"time"
)

//!+
func main() {
	fmt.Println("Commencing countdown.")
	tick := time.Tick(1 * time.Second)                //返回一个时间类型的接收单向通道
	for countdown := 10; countdown > 0; countdown-- { //结果是每隔一秒打印一次数字
		fmt.Println(countdown)
		<-tick //这是啥意？
	}
	launch()
}

//!-

func launch() {
	fmt.Println("Lift off!")
}
