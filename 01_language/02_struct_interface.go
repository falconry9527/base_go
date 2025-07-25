package main

import "fmt"

type PayMethod interface {
}

type CreditCard struct {
	balance int
	limit   int
}

func (c *CreditCard) Pay(amout int) {
	if c.balance < amout {
		fmt.Println("余额不足")
		return
	}
	c.balance -= amout
}

func anyParam(param interface{}) {
	fmt.Println("param: ", param)
}

func main() {
	c := CreditCard{balance: 100, limit: 1000}
	c.Pay(20)

	var a PayMethod = &c
	fmt.Println("a.Pay: ", a)

	var b interface{} = &c
	fmt.Println("b: ", b)

	anyParam(c)
	anyParam(a)
	anyParam(b)

}
