// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 194.

// Http3 is an e-commerce server that registers the /list and /price
// endpoints by calling (*http.ServeMux).Handle.
package main

import (
	"fmt"
	"log"
	"net/http"
)

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

//!+main

/*  Handler接口原型；
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}*/

func main() {
	db := database{"shoes": 50, "socks": 5}
	//mux变量满足Handler接口；因为mux类型实现了ServeHTTP(ResponseWriter, *Request)方法；
	mux := http.NewServeMux()
	//http.HandlerFunc(db.list)表示对db.list函数强制转换HandlerFunc函数类型；
	//并且HandlerFun函数类型 满足Handler接口； 巧妙的用法，这是编程范式吗？？
	//db.list表示方法值，HandlerFunc类型的值，所以db.list满足Handler接口；
	mux.Handle("/list", http.HandlerFunc(db.list))  //注册Handle处理器1；
	mux.Handle("/price", http.HandlerFunc(db.price)) //注册Handle处理器2；

	log.Fatal(http.ListenAndServe("localhost:8000", mux)) //mux实体实现了Handler接口；
}

//总结： mux，db.list，Handler之间的关联，实现、调用流程？？ 2019年9月27日16:35:50
//1. HandlerFunc类型满足Handler接口， 将db.list强制转换成HandlerFunc类型；
//2. mux类型满足Handler接口,mux相当于处理器集合，将db.list（HandlerFunc类型）函数注册到mux里；
//3. mux里的多个handler函数，通过字符串来唯一识别，比如"/list"、"/price"
//4. 将mux做为到ListenAndServe函数里Handler接口参数的值，因为mux满足Handler接口；
//5. 当server收到client请求后，ListenAndServe会调用Handler接口值，也就是mux实体；
//6. 然后进一步调用mux实现的对应Handler接口的方法ServeHTTP（他就是handler的实际实现过程）；
//7. mux.ServeHTTP方法的实现内容是：会根据http路径找到注册器集合里对应的handler接口值（HandlerFunc）;
//8. 然后调用对应的handler接口值（HandlerFunc函数类型）的ServeHTTP方法；
//9. HandlerFunc.ServeHTTP的内部实现就是调用 HandlerFunc类型的接收器f，也就是db.list、db.price函数值的实现；

type database map[string]dollars

func (db database) list(w http.ResponseWriter, req *http.Request) {
	for item, price := range db {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

func (db database) price(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	price, ok := db[item]
	if !ok {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
		return
	}
	fmt.Fprintf(w, "%s\n", price)
}

//!-main

/*
//!+handlerfunc
package http

type HandlerFunc func(w ResponseWriter, r *Request)

func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
	f(w, r)
}
//!-handlerfunc
*/
