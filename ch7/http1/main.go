// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 191.

// Http1 is a rudimentary e-commerce server.
package main

import (
	"fmt"
	"log"
	"net/http"
)

//!+main

func main() {
	db := database{"shoes": 50, "socks": 5}
	log.Fatal(http.ListenAndServe("localhost:8000", db)) //database类型的变量db 满足Handler接口，最后会调用db的ServeHTTP方法；
}

type dollars float32

//定义dollars类型的方法String()，满足fmt.Stringer接口；
func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

type database map[string]dollars

//定义database类型的方法ServeHTTP(w http.ResponseWriter, req *http.Request)，满足http.Handler接口；
func (db database) ServeHTTP(w http.ResponseWriter, req *http.Request) { //http.ResponseWriter为一个接口类型；
	for item, price := range db {
		fmt.Fprintf(w, "%s: %s\n", item, price)  //这里的price，会调用dollars类型的String方法；
		//http.ResponseWriter接口满足io.Writer接口；
	}
}

//!-main

/*
//!+handler
package http

type Handler interface {
	ServeHTTP(w ResponseWriter, r *Request)
}

func ListenAndServe(address string, h Handler) error
//!-handler
*/
