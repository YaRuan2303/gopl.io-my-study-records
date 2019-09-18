// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 261.
//!+

// Package bank provides a concurrency-safe bank with one account.
package bank

var deposits = make(chan int) // send amount to deposit
var balances = make(chan int) // receive balance
var msg = make(char bool)


func Deposit(amount int) { deposits <- amount }
func Balance() int       { return <-balances }

func WithDraw(amout int) (msg chan bool, err) {
	balance := Balance()
	if amout > balance {
		return msg <- false, err := fmt.Error("the balance is not enough!")
	}
	//balances <- balance -= amout
	
	return msg <- true, nil
}

func teller() {
	var balance int // balance is confined to teller goroutine
	for {
		select {
		case amount := <-deposits:
			balance += amount
		case balances <- balance:

		case drawMsg := <- msg:
			if drawMsg {
				fmt.println
			}
		}
	}
}

func init() {
	go teller() // start the monitor goroutine
}

//!-
