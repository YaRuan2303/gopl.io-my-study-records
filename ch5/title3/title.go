// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 153.

// Title3 prints the title of an HTML document specified by a URL.
package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/net/html"
)

//const表示声明一个常量；
const day = 24 * time.Hour //Hour为Duration类型的一个常量，day也为Duration类型；
fmt.Println(day.Seconds())  //???这是啥意思

func (c Celsius) String() string {
	return fmt.Sprintf("%g°C", c)
}

// Copied from gopl.io/ch5/outline2.
func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}
//基于以上原因， 安全的做法是有选择性的recover。 换句话说， 只恢复应该被恢复的panic异
//常， 此外， 这些异常所占的比例应该尽可能的低。 为了标识某个panic是否应该被恢复， 我们
//可以将panic value设置成特殊类型。 在recover时对panic value进行检查， 如果发现panic
//value是特殊类型， 就将这个panic作为errror处理， 如果不是， 则按照正常的panic进行处理
//（ 在下面的例子中， 我们会看到这种方式） 。

//下面的例子是title函数的变形， 如果HTML页面包含多个 <title> ， 该函数会给调用者返回一
//个错误（ error） 。 在soleTitle内部处理时， 如果检测到有多个 <title> ， 会调用panic， 阻止
//函数继续递归， 并将特殊类型bailout作为panic的参数。


//!+
// soleTitle returns the text of the first non-empty title element
// in doc, and an error if there was not exactly one.
func soleTitle(doc *html.Node) (title string, err error) {
	type bailout struct{}

	//该defer体 用来捕获soleTitle可能产生的panic;
	defer func() {
		switch p := recover(); p { //recover()捕获程序恐慌异常
		case nil:
			// no panic
		case bailout{}:  //这是啥意思？？？表示可接受的panic，按错误处理来做，不让程序崩溃退出；
			// "expected" panic
			err = fmt.Errorf("multiple title elements")
		default:
			panic(p) // unexpected panic; carry on panicking   //这是啥意思？？？ 正常panic,程序退出？
		}
	}()

	//这个forEachNode函数是什么机制？实现什么功能？流程？？
	//forEachNode里面直接调用匿名函数，该匿名函数继承了变量title，可以操作它title
	// Bail out of recursion if we find more than one non-empty title.
	//这里的执行流程是：
	//1、先执行forEachNode内容，然后会调用匿名函数func
	forEachNode(doc, func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" &&
			n.FirstChild != nil {
			if title != "" {
				panic(bailout{}) // multiple title elements 调用panic函数产生恐慌,程序要退出，如果没有捕获处理，则整个程序将退出；
			}
			title = n.FirstChild.Data
		}
	}, nil)

	if title == "" {
		return "", fmt.Errorf("no title element")
	}
	return title, nil
}

//在上例中， deferred函数调用recover， 并检查panic value。 当panic value是bailout{}类型时，
//deferred函数生成一个error返回给调用者。 当panic value是其他non-nil值时， 表示发生了未知
//的panic异常， deferred函数将调用panic函数并将当前的panic value作为参数传入； 此时， 等
//同于recover没有做任何操作。 （ 请注意： 在例子中， 对可预期的错误采用了panic， 这违反了
//之前的建议， 我们在此只是想向读者演示这种机制（指panic(bailout{})??）。 ）
//有些情况下， 我们无法恢复。 某些致命错误会导致Go在运行时终止程序， 如内存不足
//!-
func title(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	// Check Content-Type is HTML (e.g., "text/html; charset=utf-8").
	ct := resp.Header.Get("Content-Type")
	if ct != "text/html" && !strings.HasPrefix(ct, "text/html;") {
		resp.Body.Close()
		return fmt.Errorf("%s has type %s, not text/html", url, ct)
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close() //Body用完就释放资源；
	if err != nil {
		return fmt.Errorf("parsing %s as HTML: %v", url, err)
	}
	title, err := soleTitle(doc)
	if err != nil {
		return err
	}
	fmt.Println(title)
	return nil
}

func main() {
	for _, arg := range os.Args[1:] {
		if err := title(arg); err != nil {
			fmt.Fprintf(os.Stderr, "title: %v\n", err)
		}
	}
}
