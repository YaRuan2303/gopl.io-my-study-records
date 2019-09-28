// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 198.

// Package eval provides an expression evaluator.
package eval

import (
	"fmt"
	"math"
)

//!+env

type Env map[Var]float64  //Var存什么内容？？  Env 为map类型；

//!-env

//!+Eval1

func (v Var) Eval(env Env) float64 {   //Var类型满足Expr接口； Eval是什么意思？？
	return env[v]    //表示返回字符串v表达式的值？
}

func (l literal) Eval(_ Env) float64 {   //literal类型满足Expr接口；
	return float64(l)  //float64(l)类型转换；
}

//!-Eval1

//!+Eval2

func (u unary) Eval(env Env) float64 {     //unary类型满足Expr接口；
	switch u.op {
	case '+':
		return +u.x.Eval(env)  //TODO: x为接口，x.Eval(env)最后是调的谁的方法？？ unary的？这不成了递归调用？
		//x.Eval(env) TODO: 这又是调用谁的方法啊？？ 在parse函数里有分析指定出该接口x的实体；
	case '-':
		return -u.x.Eval(env)
	}
	panic(fmt.Sprintf("unsupported unary operator: %q", u.op))
}

func (b binary) Eval(env Env) float64 {   //binary类型满足Expr接口；
	switch b.op {
	case '+':
		return b.x.Eval(env) + b.y.Eval(env)
	case '-':
		return b.x.Eval(env) - b.y.Eval(env)
	case '*':
		return b.x.Eval(env) * b.y.Eval(env)
	case '/':
		return b.x.Eval(env) / b.y.Eval(env)
	}
	panic(fmt.Sprintf("unsupported binary operator: %q", b.op))
}

func (c call) Eval(env Env) float64 {    //call类型满足Expr接口；
	switch c.fn {
	case "pow":
		return math.Pow(c.args[0].Eval(env), c.args[1].Eval(env))
	case "sin":
		return math.Sin(c.args[0].Eval(env))
	case "sqrt":
		return math.Sqrt(c.args[0].Eval(env))  //args[0]为一个Expr接口，那这里是谁实现该接口？
		//TODO:  c.args[0].Eval(env) 这句是什么意思啊？？？
	}
	panic(fmt.Sprintf("unsupported function call: %s", c.fn))
}

//!-Eval2
