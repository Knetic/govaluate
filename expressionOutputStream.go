/*
The MIT License (MIT)

Copyright (c) 2014-2016 George Lester

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package govaluate

import (
	"bytes"
)

/*
	Holds a series of "transactions" which represent each token as it is output by an outputter (such as ToSQLQuery()).
	Some outputs (such as SQL) require a function call or non-c-like syntax to represent an expression.
	To accomplish this, this struct keeps track of each translated token as it is output, and can return and rollback those transactions.
*/
type expressionOutputStream struct {
	transactions []string
}

func (this *expressionOutputStream) add(transaction string) {
	this.transactions = append(this.transactions, transaction)
}

func (this *expressionOutputStream) rollback() string {

	index := len(this.transactions) - 1
	ret := this.transactions[index]

	this.transactions = this.transactions[:index]
	return ret
}

func (this *expressionOutputStream) createString(delimiter string) string {

	var retBuffer bytes.Buffer
	var transaction string

	penultimate := len(this.transactions) - 1

	for i := 0; i < penultimate; i++ {

		transaction = this.transactions[i]

		retBuffer.WriteString(transaction)
		retBuffer.WriteString(delimiter)
	}
	retBuffer.WriteString(this.transactions[penultimate])

	return retBuffer.String()
}
