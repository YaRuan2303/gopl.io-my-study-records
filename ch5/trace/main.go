// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 146.

// The trace program uses defer to add entry/exit diagnostics to a function.
package main

import (
	"fmt"
	"log"
	"time"
)
//调试复杂程序时， defer机制也常被用于记录何时进入和退出函数。 下例中的
//bigSlowOperation函数， 直接调用trace记录函数的被调情况。 bigSlowOperation被调时，
//trace会返回一个函数值， 该函数值会在bigSlowOperation退出时被调用。 通过这种方式， 我
//们可以只通过一条语句控制函数的入口和所有的出口？？， 甚至可以记录函数的运行时间， 如例
//子中的start。 需要注意一点： 不要忘记defer语句后的圆括号(key...!!!)， 否则本该在进入时执行的操作
//会在退出时执行， 而本该在退出时执行的， 永远不会被执行。
//
//我们知道， defer语句中的函数会在return语句更新返回值变量后再执行， 又因为在函数中定义
//的匿名函数可以访问该函数包括返回值变量在内的所有变量， 所以， 对匿名函数采用defer机
//制， 可以使其观察函数的返回值？？啥意思？。
//!+main



//如果先func,在trace，如下结果：这说明什么？
//2019/09/18 18:17:55 enter bigSlowOperation
//double(4) = 8
//double result is  8 ？为什么会先运行
//2019/09/18 18:17:55 exit bigSlowOperation (28.0016ms) 退出前运行

//如果先trace,再func，如下结果：这说明什么？
//2019/09/18 18:21:07 enter bigSlowOperation 按顺序执行
//double(4) = 8  //程序退出前一刻执行
//2019/09/18 18:21:07 exit bigSlowOperation (20.0011ms) //程序退出前一刻执行
//double result is  8

//如下double函数的执行流程（思路流程整理清楚，key）：
//1.defer里的匿名函数是先执行还是最后执行？？？

func double(x int) (result int) {
	//defer trace("bigSlowOperation")()
	defer func() { fmt.Printf("double(%d) = %d\n", x,result) }()   //分析这句，key...result=8why

	return x + x
	//执行流程：
	//先x + x = 8；在return 之前执行defer体
}

//double(4) = 8
//double result is  8


//如下bigSlowOperation函数的执行流程（思路流程整理清楚，key）：
//1. 先执行trace内部的流程，然后返回一个函数值，该值属于defer体内容，在big函数退出前一刻被执行；
//2. 然后执行time.sleep睡10s，睡完后big程序要准备退出了，这时执行defer体里的函数值程序；

func bigSlowOperation() { //功能：记录该函数被调用的初始时间和结束时间，还有该函数运行时间；
	defer trace("bigSlowOperation")() // don't forget the extra parentheses //这是什么用法啊？？
	// ...lots of work...
	time.Sleep(10 * time.Second) // simulate slow operation by sleeping
}

func trace(msg string) func() {
	start := time.Now()
	log.Printf("enter %s", msg)
	return func() { log.Printf("exit %s (%s)", msg, time.Since(start)) }
}

//!-main

func main() {
	bigSlowOperation()
	fmt.Println("double result is ", double(4))

}

/*
!+output
$ go build gopl.io/ch5/trace
$ ./trace
2015/11/18 09:53:26 enter bigSlowOperation
2015/11/18 09:53:36 exit bigSlowOperation (10.000589217s)
!-output
*/

