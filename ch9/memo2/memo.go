// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 275.

// Package memo provides a concurrency-safe memoization a function of
// type Func.  Concurrent requests are serialized by a Mutex.
package memo

import "sync"

// Func is the type of the function to memoize.
type Func func(string) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

func New(f Func) *Memo {
	return &Memo{f: f, cache: make(map[string]result)}
}

//!+

type Memo struct {
	f     Func
	mu    sync.Mutex // guards cache
	cache map[string]result
}

// Get is concurrency-safe.
func (memo *Memo) Get(key string) (value interface{}, err error) {
	memo.mu.Lock()
	res, ok := memo.cache[key]
	if !ok {
		res.value, res.err = memo.f(key) //主要是这里花时间和资源，在这里把锁释放出去？
		memo.cache[key] = res
	}
	memo.mu.Unlock()
	return res.value, res.err
}

// 使用总结：
// 被锁的区域颗粒度尽量细小，保证更多的并发，提高效率；
//!-
