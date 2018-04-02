package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Kind int

const (
	Print Kind = iota + 1
	Lparen
	Rparen
	Plus
	Minus
	Multi
	Divi
	Assign
	VarName
	IntNum
	EOF
	Others
)

type Token struct {
	Kind  Kind
	Value string
}

type Calc struct {
	tokens    []*Token
	variables map[string]int
	stack     []int
	generator *TokenGenerator
}

func newCalc() *Calc {
	c := new(Calc)
	c.tokens = make([]*Token, 0)
	c.variables = make(map[string]int)
	c.stack = make([]int, 0, 10)
	c.generator = nil
	return c
}

func (c *Calc) setNewLine(line string) {
	c.generator = newTokenGenerator(line)
}

func (c *Calc) push(val int) {
	c.stack = append(c.stack, val)
}

func (c *Calc) pop() (int, bool) {
	length := len(c.stack)
	if length > 0 {
		val := c.stack[length-1]
		c.stack = c.stack[:length-1]
		return val, true
	}
	return 0, false
}

func (c *Calc) statement() {
	if tk, ok := c.generator.readNext(); ok {
		switch tk.Kind {
		case VarName:
			c.generator.next()
			varName := tk.Value
			if tk, ok := c.generator.next(); ok && tk.Kind == Assign {
				c.expression()
				if c.variables[varName], ok = c.pop(); !ok {
					fmt.Printf("variables not found:%s\n", varName)
				}
			} else {
				fmt.Println("token should be '='")
			}
		case Print:
			c.generator.next()
			c.expression()
			if val, ok := c.pop(); ok {
				fmt.Printf("Answer: %d\n", val)
			}
		default:
			c.expression()
		}
	}
}

func (c *Calc) expression() {
	c.term()
	for tk, ok := c.generator.readNext(); ok && (tk.Kind == Plus || tk.Kind == Minus); tk, ok = c.generator.readNext() {
		c.generator.next()
		op := tk.Kind
		c.term()
		c.operate(op)
	}
}

func (c *Calc) term() {
	c.factor()

	for tk, ok := c.generator.readNext(); ok && (tk.Kind == Multi || tk.Kind == Divi); tk, ok = c.generator.readNext() {
		c.generator.next()
		op := tk.Kind
		c.factor()
		c.operate(op)
	}
}

func (c *Calc) factor() {
	if tk, ok := c.generator.next(); ok {
		switch tk.Kind {
		case VarName:
			if val, ok := c.variables[tk.Value]; ok {
				c.push(val)
			} else {
				fmt.Printf("variables not found:%s\n", tk.Value)
			}
		case IntNum:
			if num, err := strconv.Atoi(tk.Value); err == nil {
				c.push(num)
			} else {
				fmt.Printf("value cannot convert to int:%s\n", tk.Value)
			}
		case Lparen:
			c.expression()
			if tk, ok := c.generator.next(); tk.Kind != Rparen && ok {
				fmt.Println("Syntax Error: missing ')' ")
			}
		}
	} else {
		fmt.Println("missing token")
	}

}

func (c *Calc) operate(op Kind) bool {
	var d1, d2 int
	var ok bool
	if d2, ok = c.pop(); ok {
	} else {
		fmt.Println("Invalid argument")
		return false
	}

	if d1, ok = c.pop(); ok {
	} else {
		fmt.Println("Invalid argument")
		return false
	}

	if op == Divi && d2 == 0 {
		fmt.Println("Zero Divide Error")
		return false
	}

	switch op {
	case Plus:
		c.push(d1 + d2)
	case Minus:
		c.push(d1 - d2)
	case Multi:
		c.push(d1 * d2)
	case Divi:
		c.push(d1 / d2)
	default:
		fmt.Println("Invalid operator")
		return false
	}
	return true
}

type TokenGenerator struct {
	source       []rune
	maxIndex     int
	currentIndex int
}

func newTokenGenerator(text string) *TokenGenerator {
	tg := new(TokenGenerator)
	tg.source = []rune(text)
	tg.maxIndex = len(text)
	tg.currentIndex = 0
	return tg
}

func (tg *TokenGenerator) tokenGenerator(index int) (*Token, bool, int) {
	//skip spaces
	for index < tg.maxIndex && tg.source[index] == ' ' {
		index++
	}
	if index >= tg.maxIndex {
		return nil, false, index
	}
	var token = &Token{Kind: Others, Value: ""}
	ch := tg.source[index]
	if isDigit(ch) {
		token.Kind = IntNum
		for isDigit(ch) {
			token.Value = token.Value + string(ch)
			index++
			if index >= tg.maxIndex {
				break
			} else {
				ch = tg.source[index]
			}
		}
	} else if num, ok := toLower(ch); ok {
		token.Value = string(num)
		token.Kind = VarName
		index++
	} else {
		switch ch {
		case '(':
			token.Kind = Lparen
		case ')':
			token.Kind = Rparen
		case '+':
			token.Kind = Plus
		case '-':
			token.Kind = Minus
		case '*':
			token.Kind = Multi
		case '/':
			token.Kind = Divi
		case '=':
			token.Kind = Assign
		case '?':
			token.Kind = Print
		}
		index++
	}
	return token, true, index
}

func (tg *TokenGenerator) next() (*Token, bool) {
	tk, ok, index := tg.tokenGenerator(tg.currentIndex)
	if ok {
		tg.currentIndex = index
	}
	return tk, ok
}

func (tg *TokenGenerator) readNext() (*Token, bool) {
	tk, ok, _ := tg.tokenGenerator(tg.currentIndex)
	return tk, ok
}

func (tg *TokenGenerator) readFirstToken() (*Token, bool) {
	tk, ok, _ := tg.tokenGenerator(0)
	return tk, ok
}

// Utils
func isDigit(r rune) bool {
	if _, err := strconv.Atoi(string(r)); err == nil {
		return true
	}
	return false
}

func toLower(r rune) (rune, bool) {
	if 'a' <= r && r <= 'z' {
		return r, true
	} else if 'A' <= r && r <= 'Z' {
		return []rune(strings.ToLower(string(r)))[0], false
	} else {
		return '0', false
	}
}

func main() {
	var err error
	calc := newCalc()
	reader := bufio.NewReaderSize(os.Stdin, 4096)
	for line := ""; err == nil; line, err = reader.ReadString('\n') {
		calc.setNewLine(line)
		calc.statement()
		fmt.Print("input>")
	}
	if err != io.EOF {
		panic(err)
	} else {
		fmt.Println("\nbyebye!")
		os.Exit(0)
	}
}
