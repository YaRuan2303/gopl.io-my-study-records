// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

package eval

import (
	"fmt"
	"math"
	"testing"
)

//!+Eval
func TestEval(t *testing.T) {
	tests := []struct {  //初始化切片内容；
		expr string
		env  Env
		want string
	}{
		{"sqrt(A / pi)", Env{"A": 87616, "pi": math.Pi}, "167"},

		{"pow(x, 3) + pow(y, 3)", Env{"x": 12, "y": 1}, "1729"},
		{"pow(x, 3) + pow(y, 3)", Env{"x": 9, "y": 10}, "1729"},
		{"5 / 9 * (F - 32)", Env{"F": -40}, "-40"},
		{"5 / 9 * (F - 32)", Env{"F": 32}, "0"},
		{"5 / 9 * (F - 32)", Env{"F": 212}, "100"},
		//!-Eval
		// additional tests that don't appear in the book
		{"-1 + -x", Env{"x": 1}, "-2"},
		{"-1 - x", Env{"x": 1}, "-2"},
		//!+Eval
	}
	var prevExpr string
	for _, test := range tests {
		// Print expr only when it changes.
		if test.expr != prevExpr {
			fmt.Printf("\n%s\n", test.expr)
			prevExpr = test.expr
		}

		//Parse功能：根据提供的字符串形式的表达式，识别出该表达式属于哪种类型的表达式，比如call类型，binary类型等等；
		expr, err := Parse(test.expr)  //解析字符串形式的表达式，比如"sqrt(A / pi)"，返回表达式接口（那当前是谁实现了该接口）；
		//Parse解析字符串表达式后，返回处理该表达式的实体对象；以接口值expr形式返回；可以%T打印该接口值对应的实体类型；
		fmt.Printf("11111111111test.expr is %s;  exprt interface 's type  is %T\n", test.expr, expr)
		//!output: 11111111111exprt interface 's type  is eval.call、eval.binary、eval.binary

		if err != nil {
			t.Error(err) // parse error
			continue
		}
		//调用实体映射的接口方法expr.Eval(test.env)；返回值为表达式计算出的结果值
		got := fmt.Sprintf("%.6g", expr.Eval(test.env)) //test.env存着要计算的值

		fmt.Printf("\t%v => %s\n", test.env, got)
		//map[A:87616 pi:3.141592653589793] => 167   注意打印内容：%v格式 打印出test.env变量的类型
		//got为float64类型，打印出他的%s格式，就会调got实现的string方法，如果没有则打印fmt的string方法，
		if got != test.want {
			t.Errorf("%s.Eval() in %v = %q, want %q\n",
				test.expr, test.env, got, test.want)
			fmt.Printf("222222222222222got = %q\n", got)   //222222222222222got = "167"  %q格式答打印出的内容；
		}
	}
}

//!-Eval

/*
//!+output
sqrt(A / pi)
	map[A:87616 pi:3.141592653589793] => 167

pow(x, 3) + pow(y, 3)
	map[x:12 y:1] => 1729
	map[x:9 y:10] => 1729

5 / 9 * (F - 32)
	map[F:-40] => -40
	map[F:32] => 0
	map[F:212] => 100
//!-output

// Additional outputs that don't appear in the book.

-1 - x
	map[x:1] => -2

-1 + -x
	map[x:1] => -2
*/

func TestErrors(t *testing.T) {
	for _, test := range []struct{ expr, wantErr string }{
		{"x % 2", "unexpected '%'"},
		{"math.Pi", "unexpected '.'"},
		{"!true", "unexpected '!'"},
		{`"hello"`, "unexpected '\"'"},
		{"log(10)", `unknown function "log"`},
		{"sqrt(1, 2)", "call to sqrt has 2 args, want 1"},
	} {
		expr, err := Parse(test.expr)
		if err == nil {
			vars := make(map[Var]bool)
			err = expr.Check(vars)
			if err == nil {
				t.Errorf("unexpected success: %s", test.expr)
				continue
			}
		}
		fmt.Printf("%-20s%v\n", test.expr, err) // (for book)
		if err.Error() != test.wantErr {
			t.Errorf("got error %s, want %s", err, test.wantErr)
		}
	}
}

/*
//!+errors
x % 2               unexpected '%'
math.Pi             unexpected '.'
!true               unexpected '!'
"hello"             unexpected '"'

log(10)             unknown function "log"
sqrt(1, 2)          call to sqrt has 2 args, want 1
//!-errors
*/
