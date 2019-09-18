// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 221.
//!+

// Netcat1 is a read-only TCP client.
package main

import (
	"io"
	"log"
	"net"
	"os"
)

/* var port = flag.String("port", "8000", "sever's port")

func main() {
	flag.Parse()
	addr := "localhost:" + *port
	go getTime(&addr)

	//for sleep 1s

}

func getTime(addr *string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err) //程序退出
	}
	defer conn.Close()

	mustCopy(os.Stdout, conn)

}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
*/

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	mustCopy(os.Stdout, conn)
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}

//!-

// Netcat1 is a read-only TCP client.
/* package main

import (
        "io"
        "log"
        "net"
        "os"
        "strings"
        "time"
)

func main() {
        for _, v := range os.Args[1:] {
                keyValue := strings.Split(v, "=")
                go connTcp(keyValue[1])
        }
        for {
                time.Sleep(1 * time.Second)  //目的是不让主程序退出
        }
}

func connTcp(uri string) {
        conn, err := net.Dial("tcp", uri)
        if err != nil {
                log.Fatal(err)
        }
        defer conn.Close()
        mustCopy(os.Stdout, conn)

}

func mustCopy(dst io.Writer, src io.Reader) {
        if _, err := io.Copy(dst, src); err != nil {
                log.Fatal(err)
        }
} */
