// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 247.

//!+main

// The du1 command computes the disk usage of the files in a directory.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

//功能：根据指定的目录，计算出该目录下的文件个数和文件总大小
//实现方案和流程：
//两个线程，一个无缓冲通道；
//一个线程专门来找出目录下的文件，并把文件大小发送到通道里，找完所有文件后，关闭通道
//另一个线程为主线程，阻塞等待从通道里获取文件大小，进行累加计算，

func main() {
	// Determine the initial directories.
	flag.Parse()
	roots := flag.Args() //返回字符串切片，表示可以指定多个目录
	if len(roots) == 0 {
		roots = []string{"."}
	}

	// Traverse the file tree.
	fileSizes := make(chan int64)
	go func() {
		for _, root := range roots {
			walkDir(root, fileSizes)
		}
		close(fileSizes)
	}()

	// Print the results.
	var nfiles, nbytes int64
	for size := range fileSizes {
		nfiles++
		nbytes += size
	}
	printDiskUsage(nfiles, nbytes)
}

func printDiskUsage(nfiles, nbytes int64) {
	fmt.Printf("%d files  %.1f GB\n", nfiles, float64(nbytes)/1e9)
}

//!-main

//!+walkDir

// walkDir recursively递归 walks the file tree rooted at dir
// and sends the size of each found file on fileSizes.
func walkDir(dir string, fileSizes chan<- int64) { //表示该通道只能用于发送
	for _, entry := range dirents(dir) {
		if entry.IsDir() { // 是目录
			subdir := filepath.Join(dir, entry.Name()) //这是干嘛？
			walkDir(subdir, fileSizes)                 //递归
		} else { // 是文件
			fileSizes <- entry.Size()
		}
	}
}

// dirents returns the entries of directory dir.
func dirents(dir string) []os.FileInfo {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du1: %v\n", err)
		return nil
	}
	return entries
}

//!-walkDir

// The du1 variant uses two goroutines and
// prints the total after every file is found.
