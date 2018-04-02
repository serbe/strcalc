//
// Tiny calculator
//
// cal := modules.NewCal()
// fmt.Println(cal.MustCal("1x1+4/2x(1+2)"))
// result, err := cal.Cal("1x1+4/2x(1+2)")
// fmt.Println(cal.GetPostfixExpr("1x1+4/2x(1+2)"))
//

package main

import (
	"errors"
	"strconv"
	"strings"
)

func NewCal() *Cal {
	return &Cal{
		stack:   &Stack{},
		opStack: &Stack{},
	}
}

func toNumber(n string) int {
	iN, err := strconv.Atoi(n)
	if err != nil {
		panic(err)
	}
	return iN
}

func toString(iN int) string {
	return strconv.Itoa(iN)
}

var FnMap = map[string]func(string, string) string{
	"ADD": func(a, b string) string {
		return toString(toNumber(a) + toNumber(b))
	},
	"MIN": func(a, b string) string {
		return toString(toNumber(a) - toNumber(b))
	},
	"MUL": func(a, b string) string {
		return toString(toNumber(a) * toNumber(b))
	},
	// "DIV": func(a, b string) string {
	// 	if b == "0" {
	// 		panic(a + " can not divided by 0.")
	// 	}
	// 	return toString(toNumber(a) / toNumber(b))
	// },
	"REM": func(a, b string) string {
		return toString(toNumber(a) % toNumber(b))
	},
}

var OpMap = map[string]string{
	"+": "ADD",
	"-": "MIN",
	// "/": "DIV",
	"*": "MUL",
	"%": "REM",
}

type Cell struct {
	value    string
	typ      string
	op       string
	priority int // + - x / ( ) 0 0 1 1 2 2
}

type Stack struct {
	values []*Cell
}

func (s *Stack) Push(c *Cell) {
	s.values = append(s.values, c)
}

func (s *Stack) Pop() *Cell {
	if len(s.values) == 0 {
		return nil
	}
	top := s.values[len(s.values)-1]
	s.values = s.values[:len(s.values)-1]
	return top
}

func (s Stack) Top() *Cell {
	if len(s.values) == 0 {
		return nil
	}
	return s.values[len(s.values)-1]
}

type Cal struct {
	stack        *Stack
	opStack      *Stack
	_queue       []*Cell
	postfixQueue []*Cell
}

func (c Cal) isNumber(char string) bool {
	switch char {
	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		return true
	}
	return false
}

func (c *Cal) prepare(expr string) {
	splits := strings.Split(expr, "")
	count := len(splits)
	num := ""
	//group := false
	//subExpr := ""
	for i := 0; i < count; i++ {
		char := splits[i]
		if char == "" {
			continue
		}
		switch char {
		case "(", ")":
			num = ""
			c._queue = append(c._queue, &Cell{
				value:    char,
				typ:      "OP",
				priority: 2,
			})
		case "+", "-":
			num = ""
			c._queue = append(c._queue, &Cell{
				value:    char,
				typ:      "OP",
				op:       OpMap[char],
				priority: 0,
			})
		case "*", "%":
			num = ""
			c._queue = append(c._queue, &Cell{
				value:    char,
				typ:      "OP",
				op:       OpMap[char],
				priority: 1,
			})
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			num += char
			if (i+1 < count && !c.isNumber(splits[i+1])) || i == count-1 {
				c._queue = append(c._queue, &Cell{
					value:    num,
					typ:      "NUMBER",
					priority: 0,
				})
			}
		}
	}
}

func (c *Cal) postfixExpr() string {
	expr := ""
	for _, cell := range c.postfixQueue {
		expr += cell.value + " "
	}
	return expr
}

func (c *Cal) Reset() {
	c.opStack = &Stack{}
	c.stack = &Stack{}
	c._queue = []*Cell{}
	c.postfixQueue = []*Cell{}
}

func (c *Cal) postfix(expr string) string {
	c.prepare(expr)
	for _, cell := range c._queue {
		if cell.typ == "NUMBER" || cell.typ == "EXPR" {
			c.postfixQueue = append(c.postfixQueue, cell)
		} else if cell.typ == "OP" {
			if cell.value == "(" {
				c.opStack.Push(cell)
			} else if cell.value == ")" {
				for top := c.opStack.Pop(); top != nil && top.value != "("; {
					c.postfixQueue = append(c.postfixQueue, top)
					top = c.opStack.Pop()
				}
			} else {
				for top := c.opStack.Top(); top != nil && top.priority >= cell.priority && top.value != "("; {
					c.postfixQueue = append(c.postfixQueue, top)
					c.opStack.Pop() //remove
					top = c.opStack.Top()
				}
				c.opStack.Push(cell)
			}
		}
	}

	for top := c.opStack.Pop(); top != nil; {
		c.postfixQueue = append(c.postfixQueue, top)
		top = c.opStack.Pop()
	}

	return c.postfixExpr()
}

func (c Cal) GetPostfixExpr(expr string) string {
	expr = strings.Trim(expr, " ")
	if len(expr) == 0 {
		return ""
	}
	postfixExpr := c.postfix(expr)
	c.Reset()
	return postfixExpr
}

func (c Cal) Cal(expr string) (string, error) {
	var Err error
	var res *Cell
	func() {
		defer func() {
			if err := recover(); err != nil {
				Err = errors.New(err.(string))
			}
		}()

		expr = strings.Trim(expr, " ")
		if len(expr) == 0 {
			Err = errors.New("Given expr is empty.")
			return
		}
		c.postfix(expr)
		for _, cell := range c.postfixQueue {
			if cell.typ == "NUMBER" {
				c.stack.Push(cell)
			} else if cell.typ == "OP" {
				fn, ok := FnMap[cell.op]
				if !ok {
					panic("not support op " + cell.value)
				}
				b := c.stack.Pop()
				a := c.stack.Pop()
				if b == nil {
					panic("Invalid number B")
				}
				if a == nil {
					panic("Invalid number A")
				}
				c.stack.Push(&Cell{
					value: fn(a.value, b.value),
					typ:   "NUMBER",
				})
			}
		}
		res = c.stack.Pop()
		if res == nil {
			panic("Cal fail!")
		}
		c.Reset()
	}()
	result := ""
	if res != nil {
		result = res.value
	}
	return result, Err
}

func (c Cal) MustCal(expr string) string {
	res, err := c.Cal(expr)
	if err != nil {
		return ""
	}
	return res
}
