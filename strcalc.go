package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func ps2s(ps []string, src string) string {
	src = strings.Replace(src, "ps[", "", -1)
	src = strings.Replace(src, "]", "", -1)
	i, _ := strconv.Atoi(src)
	return ps[i]
}

func mul(src string) string {
	values := strings.Split(src, "*")
	v1, _ := strconv.Atoi(values[0])
	v2, _ := strconv.Atoi(values[1])
	return strconv.Itoa(v1 * v2)
}

func rem(src string) string {
	values := strings.Split(src, "%")
	v1, _ := strconv.Atoi(values[0])
	v2, _ := strconv.Atoi(values[1])
	return strconv.Itoa(v1 % v2)
}

func add(src string) string {
	values := strings.Split(src, "+")
	v1, _ := strconv.Atoi(values[0])
	v2, _ := strconv.Atoi(values[1])
	return strconv.Itoa(v1 + v2)
}

func sub(src string) string {
	values := strings.Split(src, "-")
	v1, _ := strconv.Atoi(values[0])
	v2, _ := strconv.Atoi(values[1])
	return strconv.Itoa(v1 - v2)
}

func brackets(src string) string {
	length := len(src)
	if length < 3 {
		return src
	}
	start := -1
	end := length + 1
	noclose := false
	for i := 0; i < length; i++ {
		if src[i] == '(' {
			start = i
			noclose = true
		}
		if noclose && src[i] == ')' {
			end = i
			noclose = false
		}
	}
	if start > -1 && !noclose {
		return src[start+1 : end]
	}
	return src
}

func findMul(src string) (string, int) {
	return find(src, '*')
}

func findRem(src string) (string, int) {
	return find(src, '%')
}

func findAdd(src string) (string, int) {
	return find(src, '+')
}

func findSub(src string) (string, int) {
	return find(src, '-')
}

func find(src string, b byte) (string, int) {
	length := len(src)
	if length < 3 {
		return "", -1
	}
	pos := 0
	for i := 1; i < length-1; i++ {
		if src[i] == b {
			pos = i
			break
		}
	}
	if pos > 0 {
		return getNearNumbers(src, pos)
	}
	return "", -1
}

func getNearNumbers(src string, pos int) (string, int) {
	leftValue, rightValue, leftPos := getTwoNumber(src, pos)
	return leftValue + string(src[pos]) + rightValue, leftPos
}

func getTwoNumber(src string, pos int) (string, string, int) {
	findLeft := false
	findRight := false
	leftValue := ""
	rightValue := ""
	leftPos := -1
	length := len(src)
	for i := 1; i < length; i++ {
		if pos-i < 0 {
			findLeft = true
		}
		if pos+i > length-1 {
			findRight = true
		}
		if !findLeft {
			if isNumeric(src[pos-i]) {
				leftValue = string(src[pos-i]) + leftValue
				leftPos = pos - i
			} else {
				findLeft = true
			}
		}
		if !findRight {
			if isNumeric(src[pos+i]) {
				rightValue = rightValue + string(src[pos+i])
			} else {
				findRight = true
			}
		}
		if findLeft && findRight {
			break
		}
	}
	return leftValue, rightValue, leftPos
}

func isNumeric(b byte) bool {
	return !(b == '+' || b == '-' || b == '*' || b == '%' || b == '(' || b == ')')
}

func calc(src string) string {
	f := brackets(src)
	if f != src {
		src = strings.Replace(src, "("+f+")", calc(f), 1)
		src = calc(src)
	}

	fm, pm := findMul(src)
	if pm != -1 {
		src = strings.Replace(src, fm, mul(fm), 1)
		src = calc(src)
	}
	fa, pa := findAdd(src)
	fs, ps := findSub(src)
	if pa != -1 && ps != -1 {
		if pa < ps {
			src = strings.Replace(src, fa, add(fa), 1)
			src = calc(src)
		} else {
			src = strings.Replace(src, fs, sub(fs), 1)
			src = calc(src)
		}
	}
	if pa != -1 {
		src = strings.Replace(src, fa, add(fa), 1)
		src = calc(src)
	}
	if ps != -1 {
		src = strings.Replace(src, fs, sub(fs), 1)
		src = calc(src)
	}
	fr, pr := findRem(src)
	if pr != -1 {
		src = strings.Replace(src, fr, rem(fr), 1)
		src = calc(src)
	}
	return src
}

func rotate(src []string, r int) []string {
	r = r % len(src)
	for i := 0; i < r; i++ {
		src = append(src[1:], src[0])
	}
	return src
}

func main() {
	asS := "28,37,58,177,67,83,134,122,183,139,104,200,63,129,131,138,201,2,120,5,9,78,28,178,238,228,187,5,9,78,89,138,185,22,181,193,106,192,149,62,140,252,21,202,125,94,139,190,145,80,114,31,182,52,156,180,210,205,53,190,7,112,18,144,48,110,39,170,79,16,19,190,151,94,45,119,28,22,236,149,56,108,133,66,70,147,195,178,32,181,66,186,103,226,220,163,172,220,221,212,237,23,60,118,83,151,51,159,89,201,219,118,83,151,51,13,211,146,106,120,51,211,108,160,16,226,234,212,28,84,105,119,28,26,57,61,220,26,97,60,248,199,206,139,59,35,10,190,128,162,70,175,139,252,193,202,21,32,148,211,249,50,216,121,129,127,209,58,181,39,182,202,138,127,66,223,197,56,102,194,182,74,248,194,182,74,206,194,182,74,203,194,182,74,51,185,93,3,123,185,82,212,95,200,232,51,242,189,125,170,36,187,84,222,153,187,35,137,97,177,72,56,155,177,4,173,242,138,94,160,176,94,155,195,235,47,89,241,103,144,217,201,95,208,40,75,43,66,82,144,29,40,76,17,123,138,197,131,109,45,55,157,204,96,44,148,86,173,68,185,170,212,237,7,189,188,166,68,38,171,255,199,129,163,172,27,213,163,172,215,220,119,28,69,59,58,74,78,19,219,127,191,12,212,237,7,189,137,74,254,242,59,157,101,222,66,82,123,234,150,95,179,155,160,16,230,197,122,183,139,107,47,90,72,227,86,102,106,150,122,183,137,190,47,206,51,67,142,44,194,14,160,16,230,197,52,187,144,187,158,58,172,207,109,233,125,171,47,91,139,78,151,80,140,233,147,135,210,114,142,44,194,14,144,217,233,5,159,65,110,167,159,203,20,231,138,197,131,109,47,89,23,174,128,199,199,41,144,217,233,5,198,50,167,23,144,217,204,254,35,196,26,166,198,50,167,23,198,50,243,147,159,203,20,231,5,2,137,13,144,217,204,254,104,155,53,214,92,222,146,65,88,99,149,188,80,211,172,25,80,211,13,152,168,235,93,162,80,211,4,187,163,172,175,210,80,211,186,39,80,211,13,152,35,200,187,49,5,189,153,156,46,218,73,162,125,212,203,161,212,200,246,24,13,80,147,35,80,211,19,46,80,211,19,46,80,85,157,168,194,182,74,7,80,211,19,46,37,187,116,199,119,28,138,104,13,95,231,14,141,0,7,2,139,59,35,10,192,116,142,153,159,89,195,153,163,172,217,103,168,128,29,75,139,59,28,168,167,114,196,153,34,240,231,232,13,127,111,108,95,85,50,218,139,59,223,115,80,211,186,39,149,202,0,60,122,183,243,68,202,78,227,33,82,208,118,172,167,99,193,192,167,99,64,242,77,36,18,162,194,182,74,168,142,44,194,14,88,99,35,242,35,188,227,48,91,194,42,51,35,231,85,153,172,104,238,157,93,190,142,240,104,196,150,119,35,231,63,94,35,188,239,119,35,188,244,15,34,240,231,232,34,240,231,232,159,65,110,167,167,114,166,192,35,188,131,123,35,188,245,137,35,188,245,70,159,89,31,155,128,199,198,79,160,16,230,197,123,30,236,222,109,74,75,61,185,22,173,130,185,22,173,161,182,52,22,58,103,254,94,190,222,124,22,133,167,99,72,119,212,237,63,29,159,65,142,92,198,199,69,204,115,127,68,66,193,194,69,36,101,50,1,2,195,80,140,212,125,212,207,121,139,59,2,223,36,67,231,3,124,158,88,112,128,199,121,141,142,44,137,222,67,78,143,182,217,61,15,26,119"
	as := strings.Split(asS, ",")
	psS := "80,80,8080,80,3128,3128,9090,3128,53281,53281,8081,1234,3130,3129,3129,3130,8080,8080,8080,3128,3128,3128,3128,3129,8888,2000,8080,3128,3128,80,8080,3128,53281,3128,80,80,80,80,80,8080,80,80,80,80,80,3128,3128,3128,3128,8080,8080,80,80,80,3128,3128,80,80,80,8080,3128,80,8080,1080,80,80,80,80,3128,3128,80,3128,80,8888,80,80,8080,3128,80,8080,80,80,8080,8088,8080,8080,8080,3128,808,80,14826,10098,80,54566,54566,8080,3128,3128,80,3128,8080,3128,80,3128,80,3128,80,3128,3128,3128,3128,3128,9999,31288,80,3128,8080,8080,3128,80,8080,80,80,80,80,80,80,8080,3128,3128,3128,80,80,3128,80,8080,8080,8080,8118,3128,80,80,80,3128,80,80,8118,3128,3128,8080,80,80,3128,8080,3129,3128,80,8118,80,80,80,3128,80,80,80,80,80,3130,3129,80,8080,80,80,80,3128,8118,8080,1212,80,10010,10010,8080,8080,8080,8000,8888,8080,8080,8080,3128,80,8081,3128,8888,8080,3128,3128,80,8080,3128"
	ps := strings.Split(psS, ",")
	n := "(91*1012+121*1463+82+1077+192+ps[1]+10+161+134+2155+2+45+240+ps[89]+69*ps[133]+200+4011+241*5072+270+ps[191]+106+2449+251*6035+265*6636)%800"

	// n := "(91*1012+121*1463+82+1077+192+80+10+161+134+2155+2+45+240+80+69*3128+200+4011+241*5072+270+8081+106+2449+251*6035+265*6636)%800"
	rePS := regexp.MustCompile(`(ps\[\d+\])`)
	// //	reAS := regexp.MustCompile(`(as\[\d+\])`)

	fPS := rePS.FindAllString(n, -1)
	// //	fAS := reAS.FindAllString(n, -1)

	for _, item := range fPS {
		n = strings.Replace(n, item, ps2s(ps, item), 1)
	}

	//	for _, item := range fAS {
	//		fmt.Println(item)
	//	}

	//	fmt.Println(len(as))
	//	fmt.Println(len(ps))
	fmt.Println(n)
	fmt.Println(calc(n))
	res, _ := strconv.Atoi(calc(n))
	as = rotate(as, res)
	fmt.Println(as[0])
}
