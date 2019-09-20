// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 156.

// Package geometry defines simple types for plane geometry.
//!+point
package geometry   //忘了和package main的区别？？。。

import "math"

type Point struct{ X, Y float64 }   //Point是一种类型，可以简称为类，Point类

// traditional function
func Distance(p, q Point) float64 {
	return math.Hypot(q.X-p.X, q.Y-p.Y)  //math是包，Hypot是math包里的函数；
}

// same thing, but as a method of the Point type
//思路解析：表示Point类型的方法Distance，
//方法功能:传入q，计算并返回出p和q之间的距离值
func (p Point) Distance(q Point) float64 {
	return math.Hypot(q.X-p.X, q.Y-p.Y)
}

//!-point

//!+path

// A Path is a journey connecting the points with straight lines.
type Path []Point  //声明一个类型为Path, 该Path的实际类型为[]Point切片类型，而Point为一个结构体；

// Distance returns the distance traveled along the path.
//思路解析：表示一个Path类型的方法，方法名为Distance
//方法功能：
func (path Path) Distance() float64 {
	sum := 0.0
	for i := range path {
		if i > 0 { //表示从第二个point开始, 计算出当前point和它前一个point之间的距离，然后累加算出来的所有距离的总和sum；
			sum += path[i-1].Distance(path[i]) //path[i-1]表示一个Point类型的变量point，这里Distance是point变量的方法
		}
	}
	return sum
}


p := Point{1, 2}
q := Point{4, 6}
fmt.Println(Distance(p, q)) // "5", function call
fmt.Println(p.Distance(q)) // "5", method call  p叫做Distance方法的接收器参数

/*可以看到， 上面的两个函数调用都是Distance， 但是却没有发生冲突。 第一个Distance的调用
实际上用的是包级别的函数geometry.Distance， 而第二个则是使用刚刚声明的Point， 调用的
是Point类下声明的Point.Distance方法*/
//!-path
