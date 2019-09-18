// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 240.

// Crawl1 crawls web links starting with the command-line arguments.
//
// This version quickly exhausts available file descriptors
// due to excessive concurrent calls to links.Extract.
//
// Also, it never terminates because the worklist is never closed.
package main

import (
	"fmt"
	"log"
	"os"

	"gopl.io-master/ch5/links"
)

//!+crawl
func crawl(url string) []string {
	fmt.Println(url)
	list, err := links.Extract(url) //爬取新链接, 提取
	if err != nil {
		log.Print(err) //会退出？？
	}
	return list
}

//!-crawl

//!+main
//深度优先算法爬虫

func main() {
	worklist := make(chan []string) //无缓冲通道，切片类型

	// Start with the command-line arguments.
	go func() { worklist <- os.Args[1:] }()

	// Crawl the web concurrently.
	seen := make(map[string]bool)

	for list := range worklist { //阻塞获取，worklist不能close
		for _, link := range list { //循环获取切片中的字符串
			if !seen[link] {
				seen[link] = true
				go func(link string) { //开一个协程去爬虫，把爬取结果发生给通道（阻塞等待发送）
					worklist <- crawl(link) //根据给定的link去爬虫新的链接， 如果返回为空怎么办？不阻塞了？
				}(link)
			}
		}
	}
}

//!-main

/*
//!+output
$ go build gopl.io/ch8/crawl1
$ ./crawl1 http://gopl.io/
http://gopl.io/
https://golang.org/help/

https://golang.org/doc/
https://golang.org/blog/
...
2015/07/15 18:22:12 Get ...: dial tcp: lookup blog.golang.org: no such host
2015/07/15 18:22:12 Get ...: dial tcp 23.21.222.120:443: socket:
                                                        too many open files
...
//!-output
*/
