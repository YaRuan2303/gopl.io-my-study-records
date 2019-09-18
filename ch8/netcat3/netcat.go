// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 227.

// Netcat is a simple read/write client for TCP servers.
package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

//!+
func main() { //客户端
	conn, err := net.Dial("tcp", "localhost:8000") //书里说的标准输入关闭，怎么关闭？？
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan struct{})
	go func() { //消息接收线程
		io.Copy(os.Stdout, conn) // NOTE: ignoring errors  //go程序会阻塞在这一行，等待远端发来的消息，直到对方连接关闭**，才运行下面代码段；
		log.Println("done")      //因为有通道阻塞，这句log才会打印出来  //2019/09/03 10:18:24 done
		done <- struct{}{}       // signal the main goroutine，发送协程结束的消息给主协程；
		// fmt.Println("444444444")
	}()

	//主线程负责发消息给client
	//mustCopy(conn, os.Stdin)    //程序会阻塞等待终端的输入，如果对方conn关闭，这里终端输入时，程序会报错退出，提示远程连接关闭了mustCopy实现；
	// io.Copy(conn, os.Stdin) //阻塞等待终端的输入；这句是我修改的，应该用它，不然下面的资源都没释放，还造成协程泄漏
	if _, err := io.Copy(conn, os.Stdin); err != nil { //这里是阻塞等待接收终端的输入
		log.Println(err) //打印出失败原因
	}
	fmt.Println("222222222222") //后面这些都不运行，因为是程序报错（对方连接关闭，上面一句报错，程序退出）
	conn.Close()
	fmt.Println("33333") //为什么这句不运行
	<-done               // wait for background goroutine to finish  //等待辅协程退出，达到同步效果，按顺序退出；
	fmt.Println("55555555")

}

//!-

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err) //程序退出
	}
}
