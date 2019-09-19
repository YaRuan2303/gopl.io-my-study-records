// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 151.

// Defer2 demonstrates a deferred call to runtime.Stack during a panic.
package main

import (
	"fmt"
	"os"
	"runtime"
)

//!+
func main() {
	defer printStack()
	f(3)
}

func printStack() {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	os.Stdout.Write(buf[:n])
}

//!-

func f(x int) {
	fmt.Printf("f(%d)\n", x+0/x) // panics if x == 0
	defer fmt.Printf("defer %d\n", x)
	f(x - 1)
}

/*
//!+printstack
goroutine 1 [running]:
main.printStack()
	src/gopl.io/ch5/defer2/defer.go:20
main.f(0)
	src/gopl.io/ch5/defer2/defer.go:27
main.f(1)
	src/gopl.io/ch5/defer2/defer.go:29
main.f(2)
	src/gopl.io/ch5/defer2/defer.go:29
main.f(3)
	src/gopl.io/ch5/defer2/defer.go:29
main.main()
	src/gopl.io/ch5/defer2/defer.go:15
//!-printstack
*/

//f(3)
//f(2)
//f(1)
//defer 1
//panic: runtime error: integer divide by zero
//defer 2
//defer 3
//
//goroutine 1 [running]:
//main.printStack()
//F:/Go/src/gopl.io-master/ch5/defer2/defer.go:23 +0x62
//goroutine 1 [running]:
//panic(0x4b33a0, 0x56b970)
//main.f(0x0)
//D:/Go/src/runtime/panic.go:679 +0x1c0
//F:/Go/src/gopl.io-master/ch5/defer2/defer.go:30 +0x1c5
//main.f(0x0)
//main.f(0x1)
//F:/Go/src/gopl.io-master/ch5/defer2/defer.go:30 +0x1c5
//main.f(0x1)
//F:/Go/src/gopl.io-master/ch5/defer2/defer.go:32 +0x194
//F:/Go/src/gopl.io-master/ch5/defer2/defer.go:32 +0x194
//main.f(0x2)
//main.f(0x2)
//F:/Go/src/gopl.io-master/ch5/defer2/defer.go:32 +0x194
//F:/Go/src/gopl.io-master/ch5/defer2/defer.go:32 +0x194
//main.f(0x3)
//main.f(0x3)
//F:/Go/src/gopl.io-master/ch5/defer2/defer.go:32 +0x194
//F:/Go/src/gopl.io-master/ch5/defer2/defer.go:32 +0x194
//main.main()
//main.main()
//F:/Go/src/gopl.io-master/ch5/defer2/defer.go:18 +0x57
//F:/Go/src/gopl.io-master/ch5/defer2/defer.go:18 +0x57
//
//Process finished with exit code 2
//
//将panic机制类比其他语言异常机制的读者可能会惊讶， runtime.Stack为何能输出已经被释放
//函数的信息？ 在Go的panic机制中， 延迟函数的调用在释放堆栈信息之前key..!!

//
//通常来说， 不应该对panic异常做任何处理， 但有时， 也许我们可以从异常中恢复， 至少我们
//可以在程序崩溃前， 做一些操作。 举个例子， 当web服务器遇到不可预料的严重问题时， 在崩
//溃前应该将所有的连接关闭； 如果不做任何处理， 会使得客户端一直处于等待状态。 如果web
//服务器还在开发阶段， 服务器甚至可以将异常信息反馈到客户端， 帮助调试。

//如果在deferred函数中调用了内置函数recover， 并且定义该defer语句的函数发生了panic异
//常， recover会使程序从panic中恢复， 并返回panic value。 导致panic异常的函数不会继续运
//行， 但能正常返回。 在未发生panic时调用recover， recover会返回nil。

//让我们以语言解析器为例， 说明recover的使用场景(key..!!)。 考虑到语言解析器的复杂性， 即使某个
//语言解析器目前工作正常， 也无法肯定它没有漏洞。 因此， 当某个异常出现时， 我们不会选
//择让解析器崩溃， 而是会将panic异常当作普通的解析错误， 并附加额外信息提醒用户报告此
//错误


func Parse(input string) (s *Syntax, err error) {
	defer func() {
		if p := recover(); p != nil { //recover捕获谁的异常？
			err = fmt.Errorf("internal error: %v", p)
		}
	}()
	// ...parser...
}


//不加区分的恢复所有的panic异常， 不是可取的做法； 因为在panic之后， 无法保证包级变量的
//状态仍然和我们预期一致。 比如， 对数据结构的一次重要更新没有被完整完成、 文件或者网
//络连接没有被关闭、 获得的锁没有被释放。 此外， 如果写日志时产生的panic被不加区分的恢
//复， 可能会导致漏洞被忽略。
//虽然把对panic的处理都集中在一个包下， 有助于简化对复杂和不可以预料问题的处理， 但作
//为被广泛遵守的规范， 你不应该试图去恢复其他包引起的panic。 公有的API应该将函数的运
//行失败作为error返回， 而不是panic。 同样的， 你也不应该恢复一个由他人开发的函数引起的
//panic， 比如说调用者传入的回调函数， 因为你无法确保这样做是安全的。
//有时我们很难完全遵循规范， 举个例子， net/http包中提供了一个web服务器， 将收到的请求
//分发给用户提供的处理函数。 很显然， 我们不能因为某个处理函数引发的panic异常， 杀掉整
//个进程； web服务器遇到处理函数导致的panic时会调用recover， 输出堆栈信息， 继续运行。
//gopl
//Recover捕获异常 205这样的做法在实践中很便捷， 但也会引起资源泄漏， 或是因为recover操作， 导致其他问题。
//基于以上原因， 安全的做法是有选择性的recover。 换句话说， 只恢复应该被恢复的panic异
//常， 此外， 这些异常所占的比例应该尽可能的低。 为了标识某个panic是否应该被恢复， 我们
//可以将panic value设置成特殊类型。 在recover时对panic value进行检查， 如果发现panic
//value是特殊类型， 就将这个panic作为errror处理， 如果不是， 则按照正常的panic进行处理
//（ 在下面的例子中， 我们会看到这种方式） 。
//下面的例子是title函数的变形， 如果HTML页面包含多个 <title> ， 该函数会给调用者返回一
//个错误（ error） 。 在soleTitle内部处理时， 如果检测到有多个 <title> ， 会调用panic， 阻止
//函数继续递归， 并将特殊类型bailout作为panic的参数。