// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 251.

// The du4 command computes the disk usage of the files in a directory.
package main

// The du4 variant includes cancellation:
// it terminates quickly when the user hits return.

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

//功能：
//关键点：
//1. 1.1.并发获取目录下的文件总数和总大小；1.2.并发获取系统资源时要做控制，保证同一时间只有一定数量的情况，防止系统资源耗尽】
//2. 主线程保证在所有辅助协程运行完毕在做类似关闭通道的操作，也就是同步机制
//3. 当有中断广播信号时，要保证所有协程退出，防着其一直阻塞造成泄漏； select多路监听机制

//!+1
var done = make(chan struct{}) //目的是：中断信号，让当前所有的线程都退出，该通道必须是全局变量，充当广播的效果

func cancelled() bool { //一次性的接口
	select {
	case <-done: //这句啥意思？哪种情况？通道关闭后就会触发这句；
		return true
	default:
		return false
	}
}

//!-1

func main() {
	// Determine the initial directories.
	roots := os.Args[1:]
	if len(roots) == 0 {
		roots = []string{"."}
	}

	//!+2
	// Cancel traversal when input is detected.
	go func() { // 外部中断信号？
		os.Stdin.Read(make([]byte, 1)) // read a single byte //阻塞等待输入
		close(done)
	}()
	//!-2

	// Traverse each root of the file tree in parallel.
	fileSizes := make(chan int64)
	var n sync.WaitGroup
	for _, root := range roots {
		n.Add(1)
		go walkDir(root, &n, fileSizes)
	}
	go func() {
		n.Wait() //目的是要关闭通道，要等所有线程运行完才能关闭，同步效果；如果不这样，会引起线程里的发送操作恐慌，因为通道已关闭
		close(fileSizes)
	}()

	// Print the results periodically.
	tick := time.Tick(500 * time.Millisecond)
	var nfiles, nbytes int64
loop:
	//!+3
	for {
		select { //是按
		case <-done: //表示等待外界的终止信号，如果通道关闭，则触发这句？？？看这意思好像是这样的；
			// Drain fileSizes to allow existing goroutines to finish.
			for range fileSizes { //目的是不让go协程们发送阻塞着，这边接收，让他们发送完，然后退出程序，不造成泄露，让协程优雅的退出
				// Do nothing. 这种情况是线程们正在阻塞等待发送；一直阻塞着，需要让他们退出，不造成泄漏
				//如果通道么有关闭，则这里一直阻塞着，程序无法退出；所以当所有线程退出后，这里的程序就不会阻塞，会退出
			}
			return //程序直接返回了
		case size, ok := <-fileSizes:
			// ...
			//!-3
			if !ok {
				break loop // fileSizes was closed
			}
			nfiles++
			nbytes += size
		case <-tick: //定时打印作用，定时器，定时触发，无缓冲通道
			printDiskUsage(nfiles, nbytes)
		}
	}
	printDiskUsage(nfiles, nbytes) // final totals
} //!-main

func printDiskUsage(nfiles, nbytes int64) {
	fmt.Printf("%d files  %.1f GB\n", nfiles, float64(nbytes)/1e9)
}

// walkDir recursively walks the file tree rooted at dir
// and sends the size of each found file on fileSizes.
//!+4
func walkDir(dir string, n *sync.WaitGroup, fileSizes chan<- int64) {
	defer n.Done()
	if cancelled() { //表示外部发出了终止信号，终止一切活动；为什么会需要这一步，因为表示发出了中断信号，那么线程们就不要在继续运行下去了，直接都立刻停止。起的这个作用
		return //线程结束
	}
	for _, entry := range dirents(dir) { //if entry为nil,则for循环终止
		// ...
		//!-4
		if entry.IsDir() {
			n.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			go walkDir(subdir, n, fileSizes)
		} else {
			fileSizes <- entry.Size()
		}
		//!+4
	}
}

//!-4

var sema = make(chan struct{}, 20) // concurrency-limiting counting semaphore
// 目的是防止过度申请系统资源，保证同一时间只有20个申请
// dirents returns the entries of directory dir.
//!+5
func dirents(dir string) []os.FileInfo {
	select { //多路复用，多路监听，一次性；
	case sema <- struct{}{}: // acquire token
	case <-done: //这里也检查系统是否退出的消息，防止继续运行，浪费时间
		return nil // cancelled //这种事再sema阻塞时，然后系统又发终止信号时，运行这个？？？目的是不让程序阻塞在这里；
		//当程序发出终止信号时，程序中所有的阻塞点都要停止，退出。
	}
	//问题:如果都不阻塞，会先运行谁？sema还是done?

	defer func() { <-sema }() // release token

	// ...read directory...
	//!-5

	f, err := os.Open(dir)
	if err != nil { //表示打开失败
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
		return nil
	}
	defer f.Close()

	entries, err := f.Readdir(0) // 0 => no limit; read all entries
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
		// Don't return: Readdir may return partial results.
	}
	return entries
}
