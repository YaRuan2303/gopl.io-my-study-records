// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 136.

// The toposort program prints the nodes of a DAG in topological order.
package main

import (
	"fmt"
	"sort"
)

/*
e:\Go\src>toposort.exe
1:      intro to programming
2:      discrete math
3:      data structures
4:      algorithms
5:      linear algebra
6:      calculus
7:      formal languages
8:      computer organization
9:      compilers
10:     databases
11:     operating systems
12:     networks
13:     programming languages
*/

//!+table
// prereqs maps computer science courses to their prerequisites.
var prereqs = map[string][]string{ //集合，key为要学的课程，value为前置条件
	"algorithms": {"data structures"},
	"calculus":   {"linear algebra"},

	"compilers": {
		"data structures",
		"formal languages",
		"computer organization",
	},

	"data structures":       {"discrete math"},
	"databases":             {"data structures"},
	"discrete math":         {"intro to programming"},
	"formal languages":      {"discrete math"},
	"networks":              {"operating systems"},
	"operating systems":     {"data structures", "computer organization"},
	"programming languages": {"data structures", "computer organization"},
}

//!-table

//!+main
func main() {
	for i, course := range topoSort(prereqs) {
		fmt.Printf("%d:\t%s\n", i+1, course)
	}
}

/*
e:\Go\src>go run gopl.io-master\ch5\toposort
111111111 algorithms
111111111 data structures
111111111 discrete math
111111111 intro to programming
22222222222 [intro to programming]
22222222222 [intro to programming discrete math]
22222222222 [intro to programming discrete math data structures]
22222222222 [intro to programming discrete math data structures algorithms]
111111111 calculus
111111111 linear algebra
22222222222 [intro to programming discrete math data structures algorithms linear algebra]
22222222222 [intro to programming discrete math data structures algorithms linear algebra calculus]
*/

func topoSort(m map[string][]string) []string {
	var order []string
	seen := make(map[string]bool)     //去重
	var visitAll func(items []string) // 定义一个函数变量，因为visitAll是递归函数，如果不是递归函数，就可以用短变量来定义visitAll

	// 遍历每个key对应的value切片
	visitAll = func(items []string) { //匿名函数赋值，深度优先搜索
		for _, item := range items { //遍历每个key（item）
			fmt.Println("111111111", item)
			if !seen[item] {
				seen[item] = true
				visitAll(m[item]) // (递归))深度遍历每个key对应的value切片内容

				order = append(order, item) //这是啥意思？
				fmt.Println("22222222222", order)
			}
		}
	}

	var keys []string
	for key := range m { // 找出所有的key
		keys = append(keys, key)
	}

	sort.Strings(keys) //字符串排序，升序
	// fmt.Printf("%v\n", keys)
	// [algorithms calculus compilers data structures databases discrete math formal languages networks operating systems programming languages]
	visitAll(keys)
	return order
}

//!-main
