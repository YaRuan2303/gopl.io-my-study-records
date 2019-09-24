// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 165.

// Package intset provides a set of integers based on a bit vector.
package intset

import (
	"bytes"
	"fmt"
)
//使用到位数组的代码，一般出于两个考虑：
// 1. 降低存储空间。
// 2. 加快查找效率(能迅速判断某个地元素是否在一个集合中)。


//!+intset

// An IntSet is a set of small non-negative integers.
// Its zero value represents the empty set.
//bit数组类型的含义：
//1、uint64类型的切片，切片每个成员代表一个字，该切片代表一个字集合（一个字代表一个uint64整数）
type IntSet struct {
	words []uint64
}

// Has reports whether the set contains the non-negative value x.
func (s *IntSet) Has(x int) bool {
	word, bit := x/64, uint(x%64)  //x/64表示该字x在字集合中的位置(下标值)；x%64表示x在当前集合元素uint64内的bit位（第几bit位）
	return word < len(s.words) && s.words[word]&(1<<bit) != 0  //判断第bit+1位是否为1；（(1<<bit)表示1左移bit位，然后与s.words[word]相与）
	//len(s.words)表示切片长度；
}

//该方法的主要作用是将二进制数组中表示该整数的位置为1。
//首先我们得找到该整数位于 char 数组的第几个元组中，这里利用该整数除以8即可（代码中除以8用右移三位实现），
//例如整数25位于25/8 = 3 余 1，表明该整数是用char 数组的第四个元素的第二位表示。（右边低位开始算，第二位）
//那么在该元素的第几位可以利用该整数的后三位表示（0~7刚好可以表示8个位置），即 25 & 0x7 = 1，则代表25在该元素的第二位。
//将相应位置1，可以先将整数1左移相应位数，然后与二进制数组进行或操作即可。
//————————————————
//版权声明：本文为CSDN博主「杨柳_」的原创文章，遵循 CC 4.0 BY-SA 版权协议，转载请附上原文出处链接及本声明。
//原文链接：https://blog.csdn.net/qq_37375427/article/details/79797359
// Add adds the non-negative value x to the set.
func (s *IntSet) Add(x int) {
	word, bit := x/64, uint(x%64)
	for word >= len(s.words) { //表示word下标志超出s.words集合的长度，表示s.words集合大小不够记录x了，这里用for来保证扩容，表示扩容来记录下x元素；
		s.words = append(s.words, 0)  //这是什么意思？扩容？
	}
	s.words[word] |= 1 << bit  //将第bit+1位 置1；（先将1左移bit位，然后将其结果与s.words[word]元素相或）
}

// UnionWith sets s to the union of s and t.
func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword  //合并s与t集合内容
		} else { //表示s.words集合空间不足，
			s.words = append(s.words, tword) //对s.words切片扩容
		}
	}
}

//!-intset

//!+string

// String returns the set as a string of the form "{1 2 3}".
func (s *IntSet) String() string {
	var buf bytes.Buffer  //声明一个对象；
	buf.WriteByte('{')  //调用该对象的方法来操作该对象内容；
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < 64; j++ {
			if word&(1<<uint(j)) != 0 { //表示word的第j+1位被置位1
				if buf.Len() > len("{") {
					buf.WriteByte(' ')  //这是啥意思？
				}
				fmt.Fprintf(&buf, "%d", 64*i+j)  //64*i+j表示被置位1的那个字（元素）
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}

//!-string
