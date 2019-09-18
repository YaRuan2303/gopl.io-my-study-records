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

func triple(x int) (result int) {
	defer func() { result += x }()
	return double(x)
} f



//调试复杂程序时， defer机制也常被用于记录何时进入和退出函数。 下例中的
//bigSlowOperation函数， 直接调用trace记录函数的被调情况。 bigSlowOperation被调时，
//trace会返回一个函数值， 该函数值会在bigSlowOperation退出时被调用（key）。 通过这种方式， 我
//们可以只通过一条语句控制函数的入口和所有的出口？？， 甚至可以记录函数的运行时间， 如例
//子中的start。 需要注意一点： 不要忘记defer语句后的圆括号(key...!!!)， 否则本该在进入时执行的操作
//会在退出时执行， 而本该在退出时执行的， 永远不会被执行。
//
//我们知道， defer语句中的函数会在return语句更新返回值变量后再执行， 又因为在函数中定义
//的匿名函数可以访问该函数包括返回值变量在内的所有变量， 所以， 对匿名函数采用defer机
//制， 可以使其观察函数的返回值？？啥意思？就是double函数的例子意思。
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
//1.defer里的匿名函数是先执行还是最后执行？？？ 看执行结果是最后执行；
//2. 所以意思是：关于defer和函数结合使用的注意点是：
// 2.1. 如果defer体的函数的返回值为函数值，则函数正常执行，返回的函数值延迟到最后执行
// 2.2. 如果defer体的函数的返回值不为函数值，则该函数延迟到最后执行；

//执行结果流程如下：
//init: result is  0
//11111111111111111 0
//22222222222222 8
//double(4) = 8
func double(x int) (result int) {
	//defer trace("bigSlowOperation")()
	fmt.Println("init: result is ", result)  //0
	defer func() {
		fmt.Println("22222222222222", result)
		fmt.Printf("double(%d) = %d\n", x,result)
	}()   //分析这句，key...result=8why
	fmt.Println("11111111111111111", result)
	return x + x
	//执行流程：
	//先x + x = 8；在return 之前执行defer体
}

//double(4) = 8
//double result is  8


//如下bigSlowOperation函数的执行流程（思路流程整理清楚，key）：
//1. 先执行trace内部的流程，然后返回一个函数值，该值属于defer体内容，在big函数退出前一刻被执行；
//2. 然后执行time.sleep睡10s，睡完后big程序要准备退出了，这时执行defer体里的函数值程序；

//执行结果流程如下：
//33333333333333  这是为啥？
//2019/09/18 19:36:38 enter bigSlowOperation
//2019/09/18 19:36:48 exit bigSlowOperation (10.0075724s)
//44444444444444  这是为啥？

//第二遍执行，又是如下这个结果，感觉这个结果是对的，但是为啥每次运行的结果还不唯一？？
//2019/09/18 19:40:11 enter bigSlowOperation
//33333333333333
//44444444444444
//2019/09/18 19:40:21 exit bigSlowOperation (10.0065723s)

////第三遍执行结果又不一样。。。。这是啥逻辑？
//33333333333333
//2019/09/18 19:41:31 enter bigSlowOperation
//44444444444444
//2019/09/18 19:41:41 exit bigSlowOperation (10.0065724s)
func bigSlowOperation() { //功能：记录该函数被调用的初始时间和结束时间，还有该函数运行时间；
	defer trace("bigSlowOperation")() // don't forget the extra parentheses //这是什么用法啊？？
	// ...lots of work...
	fmt.Println("33333333333333")
	time.Sleep(10 * time.Second) // simulate slow operation by sleeping
	fmt.Println("44444444444444")
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
	fmt.Println(triple(4)) // "12"

}

/*
!+output
$ go build gopl.io/ch5/trace
$ ./trace
2015/11/18 09:53:26 enter bigSlowOperation
2015/11/18 09:53:36 exit bigSlowOperation (10.000589217s)
!-output
*/

