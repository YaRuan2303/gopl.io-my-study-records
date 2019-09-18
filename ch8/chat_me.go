// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 254.
//!+

// Chat is a server that lets clients chat with each other.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

//!+broadcaster
type client chan<- string // an outgoing message channel


type server chan

/* var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages
)
 */

 
//var	entering = make(chan string) 
//var	leaving  = make(chan client)
var	messages = make(chan string) // all incoming client messages


func broadcaster() {
	clients := make(map[client]bool) // all connected clients //这里的client是(发送)通道类型
	for { //死循环,相当于C语言里的while(1)
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				cli <- msg  //消息发到每个client通道里，每个clientWriter会读到，然后发送到各个client（当前client不发，其他client都发）
			}

		case cli := <-entering: //entering里装的通道cli，表示和server连接上的client 
			clients[cli] = true //打标机，表示活跃的client

		case cli := <-leaving:  //leaving里装的通道cli，表示与server断开的连接
			delete(clients, cli) //将断开连接的client从活跃字典里移去
			close(cli) //关闭该通道
		}
	}
}

//!-broadcaster
/* 
You are 127.0.0.1:64208       $ ./netcat3
127.0.0.1:64211 has arrived    You are 127.0.0.1:64211
Hi!
127.0.0.1:64208: Hi!
127.0.0.1:64208: Hi!
Hi yourself.
127.0.0.1:64211: Hi yourself. 
 */

//!+handleConn  //这是长连接还是短连接？client发完消息就关闭连接了？
func handleConn(conn net.Conn) { // conn表示一个连接
	ch := make(chan string) // outgoing client messages 这个通道干嘛用的？
	go clientWriter(conn, ch) //s向当前连接的client发送消息，比如“You are 127.0.0.1:64208”；阻塞等待

	who := conn.RemoteAddr().String() //远程client的地址
	ch <- "You are " + who //将收到的client信息发送到ch通道里，clientWriter里会收到这个消息？这是为啥？？？

	messages <- who + " has arrived"  //组装通知消息：""某个客户端到来" 广播到其他客户端,将当前连接登录的消息广播给其他客户端；
	entering <- ch  //将ch通道（相当于数据了，通道类型的数据）转入entering通道，表示将活跃的client加入活跃通道？
	//以上是连接初始化阶段的信息；下面是连接阻塞等待，连接中的信息；



	//这块会阻塞接收消息，除非对方关闭套接字，这块就不阻塞了，往下走？？
	input := bufio.NewScanner(conn)  //NewScanner功能？读取行内容，像这种工具要了解掌握使用；
	for input.Scan() { //循环读取每一行，这读的内容是谁的？client的 //这是的模式是阻塞模式？
		messages <- who + ": " + input.Text()  //client发来的消息，这个消息会被广播出去，包括广播给该client自己；
	}
	// NOTE: ignoring potential errors from input.Err()

	//到这里表示client关闭套接字了？？
	leaving <- ch //关闭当前断开连接的client的通道代表；
	messages <- who + " has left"  //消息广播
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {  //将消息发送给指定client???
	for msg := range ch {  //接收通道，阻塞等待接收东西，然后发送给conn客户端
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors  //以行格式写入内容，对应bufio.NewScanner的读取
	}
}

//!-handleConn

//功能：？

//!+main
func main() {
	//server
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err) //程序退出
	}

	go broadcaster()  //线程，专门用来发送消息给各个客户端

	for { 
		conn, err := listener.Accept() //阻塞等待新连接
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)  //消息处理线程，来一个消息，就开一个线程去处理该消息
	}
}

//!-main
