// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 123.

// Outline prints the outline of an HTML document tree.
package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

//!+
func main() {
	for _, url := range os.Args[1:] {
		if !strings.HasPrefix(url, "https://") {
			url = "https://" + url
		}
		//fmt.Println(url)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
			os.Exit(1)
		}
		doc, err := html.Parse(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "outline: %v\n", err)
			os.Exit(1)
		}
		outline(nil, doc)
	}
}

func outline(stack []string, n *html.Node) {
	if n.Type == html.ElementNode {
		stack = append(stack, n.Data) // push tag
		fmt.Println(stack)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		outline(stack, c)
		fmt.Println("111111", stack)

	}
}

//!-
/* E:\Go\src\gopl.io-master\ch5\outline>go run main.go www.baidu.com
[html]
[html head]
111111 [html head]
[html head script]
111111 [html head script]

111111 [html head]  //这是什么意思？？
111111 [html head]
111111 [html]
111111 [html]
[html body]
111111 [html body]
[html body noscript]
111111 [html body noscript]
111111 [html body]
111111 [html body]
111111 [html]
111111 [] */

/*
[html]
[html head]
[html head script]
[html body]
[html body noscript] */
