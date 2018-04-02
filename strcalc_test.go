package main

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMul(t *testing.T) {
	n := "91*1012+121*1463"
	re := regexp.MustCompile(`(\d+\*\d+)`)
	assert.Equal(t, true, re.MatchString(n))
	f := re.FindStringSubmatch(n)[1]
	assert.Equal(t, "91*1012", f)
	m := mul(f)
	assert.Equal(t, "92092", m)
	n = strings.Replace(n, f, m, 1)
	assert.Equal(t, "92092+121*1463", n)
}

func TestRem(t *testing.T) {
	n := "91*1012+121%1463"
	re := regexp.MustCompile(`(\d+\%\d+)`)
	assert.Equal(t, true, re.MatchString(n))
	f := re.FindStringSubmatch(n)[1]
	assert.Equal(t, "121%1463", f)
	m := rem(f)
	assert.Equal(t, "121", m)
	n = strings.Replace(n, f, m, 1)
	assert.Equal(t, "91*1012+121", n)
}

func TestAdd(t *testing.T) {
	n := "82+1077+192+80+10+161"
	re := regexp.MustCompile(`(\d+\+\d+)`)
	assert.Equal(t, true, re.MatchString(n))
	f := re.FindStringSubmatch(n)[1]
	assert.Equal(t, "82+1077", f)
	a := add(f)
	assert.Equal(t, "1159", a)
	n = strings.Replace(n, f, a, 1)
	assert.Equal(t, "1159+192+80+10+161", n)
}

func TestSub(t *testing.T) {
	n := "82+1077-192+80+10-161"
	re := regexp.MustCompile(`(\d+\-\d+)`)
	assert.Equal(t, true, re.MatchString(n))
	f := re.FindStringSubmatch(n)[1]
	assert.Equal(t, "1077-192", f)
	a := sub(f)
	assert.Equal(t, "885", a)
	n = strings.Replace(n, f, a, 1)
	assert.Equal(t, "82+885+80+10-161", n)
}

func TestBracketsRegexp(t *testing.T) {
	n := "(91*1012+121*1463+82+1077+192+80+10+161+134+2155+2+45+240+80+69*3128+200+4011+241*5072+270+8081+106+2449+251*6035+265*6636)%800"
	re := regexp.MustCompile(`\(([^\(].+?)\)`)
	assert.Equal(t, true, re.MatchString(n))
	f := re.FindStringSubmatch(n)[1]
	assert.Equal(t, "91*1012+121*1463+82+1077+192+80+10+161+134+2155+2+45+240+80+69*3128+200+4011+241*5072+270+8081+106+2449+251*6035+265*6636", f)
	n = strings.Replace(n, "("+f+")", calc(f), 1)
	assert.Equal(t, "4999999%800", n)
}

func TestBrackets(t *testing.T) {
	n := "(91*1012+121*1463+82+1077+192+80+10+161+134+2155+2+45+240+80+69*3128+200+4011+241*5072+270+8081+106+2449+251*6035+265*6636)%800"
	f := brackets(n)
	assert.Equal(t, "91*1012+121*1463+82+1077+192+80+10+161+134+2155+2+45+240+80+69*3128+200+4011+241*5072+270+8081+106+2449+251*6035+265*6636", f)
	n = strings.Replace(n, "("+f+")", calc(f), 1)
	assert.Equal(t, "4999999%800", n)
}

func TestFindMul(t *testing.T) {
	n := "(91*1012+121*1463+82+1077+192+80+10+161+134+2155+2+45+240+80+69*3128+200+4011+241*5072+270+8081+106+2449+251*6035+265*6636)%800"
	fm, pm := findMul(n)
	assert.Equal(t, "91*1012", fm)
	assert.Equal(t, 1, pm)
}

func TestIsNumeric(t *testing.T) {
	n := "(8+1)"
	assert.Equal(t, false, isNumeric(n[0]))
	assert.Equal(t, true, isNumeric(n[1]))
	assert.Equal(t, false, isNumeric(n[2]))
	assert.Equal(t, true, isNumeric(n[3]))
	assert.Equal(t, false, isNumeric(n[4]))
}

func TestGetNearNumbers(t *testing.T) {
	n := "(82+1077*192+80+10)"
	f, p := getNearNumbers(n, 8)
	assert.Equal(t, "1077*192", f)
	assert.Equal(t, 4, p)
}

func TestGetTwoNumber(t *testing.T) {
	n := "(91*1012+121*1463+82+1077+192+80+10+161+134+2155+2+45+240+80+69*3128+200+4011+241*5072+270+8081+106+2449+251*6035+265*6636)%800"
	left, right, p := getTwoNumber(n, 3)
	assert.Equal(t, "91", left)
	assert.Equal(t, "1012", right)
	assert.Equal(t, 1, p)
}

func TestRotate(t *testing.T) {
	n := []string{"0", "1", "2"}
	n = rotate(n, 0)
	assert.Equal(t, n[0], "0")
	n = rotate(n, 1)
	assert.Equal(t, n[0], "1")
	n = rotate(n, 3)
	assert.Equal(t, n[0], "1")
}

// BenchmarkCalc-4   	    1000	   1501164 ns/op	 6640568 B/op	    5498 allocs/op
func BenchmarkCalc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		calc("(91*1012+121*1463+82+1077+192+80+10+161+134+2155+2+45+240+80+69*3128+200+4011+241*5072+270+8081+106+2449+251*6035+265*6636)%800")
	}
}

func BenchmarkCal(b *testing.B) {
	cal := NewCal()
	for i := 0; i < b.N; i++ {
		cal.MustCal("(91*1012+121*1463+82+1077+192+80+10+161+134+2155+2+45+240+80+69*3128+200+4011+241*5072+270+8081+106+2449+251*6035+265*6636)%800")
	}
}

func BenchmarkBracketsRegexp(b *testing.B) {
	re := regexp.MustCompile(`\(([^\(].+?)\)`)
	n := "(91*1012+121*1463+82+1077+192+80+10+161+(134+2155+2+45+240+80)+69*3128+200+4011+241*5072+270+8081+106+2449+251*6035+265*6636)%800"
	for i := 0; i < b.N; i++ {
		_ = re.FindStringSubmatch(n)[1]
	}
}

func BenchmarkBrackets(b *testing.B) {
	n := "(91*1012+121*1463+82+1077+192+80+10+161+(134+2155+2+45+240+80)+69*3128+200+4011+241*5072+270+8081+106+2449+251*6035+265*6636)%800"
	for i := 0; i < b.N; i++ {
		_ = brackets(n)
	}
}

func BenchmarkContains(b *testing.B) {
	a := "0123456789"
	n := "(91*1012+121*1463+82+1077+192+80+10+161+134+2155+2+45+240+80+69*3128+200+4011+241*5072+270+8081+106+2449+251*6035+265*6636)%800"
	length := len(n)
	for i := 0; i < b.N; i++ {
		for r := 0; r < length; r++ {
			_ = strings.Contains(a, string(n[r]))
		}
	}
}

func BenchmarkCompareEqual(b *testing.B) {
	n := "(91*1012+121*1463+82+1077+192+80+10+161+134+2155+2+45+240+80+69*3128+200+4011+241*5072+270+8081+106+2449+251*6035+265*6636)%800"
	length := len(n)
	for i := 0; i < b.N; i++ {
		for r := 0; r < length; r++ {
			c := n[r]
			_ = c == '0' || c == '1' || c == '2' || c == '3' || c == '4' || c == '5' || c == '6' || c == '7' || c == '8' || c == '9'
		}
	}
}

func BenchmarkCompareNotEqual(b *testing.B) {
	n := "(91*1012+121*1463+82+1077+192+80+10+161+134+2155+2+45+240+80+69*3128+200+4011+241*5072+270+8081+106+2449+251*6035+265*6636)%800"
	length := len(n)
	for i := 0; i < b.N; i++ {
		for r := 0; r < length; r++ {
			c := n[r]
			_ = c != '+' && c != '-' && c != '*' && c != '%' && c != '(' && c != ')'
		}
	}
}

func BenchmarkCompareEqual2(b *testing.B) {
	n := "(91*1012+121*1463+82+1077+192+80+10+161+134+2155+2+45+240+80+69*3128+200+4011+241*5072+270+8081+106+2449+251*6035+265*6636)%800"
	length := len(n)
	for i := 0; i < b.N; i++ {
		for r := 0; r < length; r++ {
			c := n[r]
			_ = c == '+' || c == '-' || c == '*' || c == '%' || c == '(' || c == ')'
		}
	}
}
