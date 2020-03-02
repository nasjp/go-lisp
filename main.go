package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"
)

const HEAPSIZE = 10000000
const NIL = 0

var heap = make(cells, 0, 10000000)

// TODO fix variable names
var ep int //environment pointer
var hp int //heap pointer
var sp int //stack pointer
var fc int //free counter
var ap int //arglist pointer

type tag int

const (
	_ tag = iota
	EMP
	NUM
	SYM
	LIS
	SUBR
	FSUBR
	FUNC
)

type flag int

const (
	FRE tag = iota
	USE
)

type cell struct {
	tag  tag
	flag flag
	name string
	val  struct {
		num  int
		bind int
		subr func() int
	}
	car int
	cdr int
}

type cells []cell

func (cs cells) getCAR(addr int) int           { return cs[addr].car }
func (cs cells) setCAR(addr int, x int)        { cs[addr].car = x }
func (cs cells) setCDR(addr int, x int)        { cs[addr].cdr = x }
func (cs cells) setTag(addr int, t tag)        { cs[addr].tag = t }
func (cs cells) setName(addr int, name string) { heap[addr].name = name }

type toktype int

const (
	_ toktype = iota
	LPAREN
	RPAREN
	QUOTE
	DOT
	NUMBER
	SYMBOL
	OTHER
)

type backtrack int

const (
	_ backtrack = iota
	GO
	BACK
)

type token struct {
	ch      rune
	flag    backtrack
	toktype toktype
	buf     []rune
}

const BUFSIZE = 256

var stok = token{flag: GO, toktype: OTHER}

const (
	EOL        = '\n'
	TAB        = '\t'
	SPACE      = ' '
	NUL   rune = 0
)

func initcell() {
	for addr := 0; addr < HEAPSIZE; addr++ {
		heap = append(heap, cell{cdr: addr + 1})
	}
	ep = makesym("nil")
	assocsym(makesym("nil"), NIL)
	assocsym(makesym("t"), makesym("t"))
}

func makesym(name string) int {
	addr := freshcell()
	heap.setTag(addr, SYM)
	heap.setName(addr, name)
	return addr
}

func freshcell() int {
	res := hp
	hp = heap.getCAR(hp)
	heap.setCDR(res, 0)
	fc--
	return res
}

func cons(car int, cdr int) int {
	addr := freshcell()
	heap.setTag(addr, LIS)
	heap.setCAR(addr, car)
	heap.setCDR(addr, cdr)
	return addr
}

func assocsym(sym int, val int) {
	ep = cons(cons(sym, val), ep)
}

func gengetchar(txt string) func() rune {
	txt += "\n"
	i := 0
	return func() rune {
		i++
		return rune(txt[i-1])
	}
}

func getchar() rune {
	reader := bufio.NewReader(os.Stdin)
	// FIXME add err handling
	input, _ := reader.ReadString('\n')
	return []rune(input)[0]
}

func numbertoken(buf []rune) bool {
	if (buf[0] == '+') || (buf[0] == '-') {
		buf = buf[1:]
		if buf[0] == NUL {
			return false
		}
	}

	for _, c := range buf {
		if c == NUL {
			break
		}
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

func isalpha(c rune) bool {
	if c < 'A' || c > 'z' {
		return false
	}
	if c > 'Z' && c < 'a' {
		return false
	}
	return true
}

func issymch(c rune) bool {
	switch c {
	case '!':
	case '?':
	case '+':
	case '-':
	case '*':
	case '/':
	case '=':
	case '<':
	case '>':
		return true
	}

	return false
}

func symboltoken(buf []rune) bool {
	if unicode.IsDigit(buf[0]) {
		return false
	}

	for _, c := range buf {
		if c == NUL {
			break
		}
		if unicode.IsDigit(c) || isalpha(c) || issymch(c) {
			continue
		}
		return false
	}
	return true
}

func gettoken() {
	if stok.flag == BACK {
		stok.flag = GO
		return
	}

	if stok.ch == ')' {
		stok.toktype = RPAREN
		stok.ch = NUL
		return
	}

	if stok.ch == '(' {
		stok.toktype = LPAREN
		stok.ch = NUL
		return
	}

	// sc := bufio.NewScanner(os.Stdin)
	// sc.Scan()

	var pos int
	// getchar := gengetchar(sc.Text())

	c := getchar()
	for c == SPACE || c == EOL || c == TAB {
		c = getchar()
	}

	switch c {
	case '(':
		stok.toktype = LPAREN
	case ')':
		stok.toktype = RPAREN
	case '\'':
		stok.toktype = QUOTE
	case '.':
		stok.toktype = DOT
	default:
		pos++
		for c != EOL && pos < BUFSIZE && c != SPACE && c != '(' && c != ')' {
			pos++
			stok.buf = append(stok.buf, c)
			c = getchar()
		}

		stok.buf = append(stok.buf, NUL)
		stok.ch = c

		if numbertoken(stok.buf) {
			stok.toktype = NUMBER
			break
		}
		if symboltoken(stok.buf) {
			stok.toktype = SYMBOL
			break
		}
		stok.toktype = OTHER
	}
}

func main() {
	initcell()
	gettoken()
	fmt.Printf("stok.buf: %v\n", stok.buf)
	fmt.Printf("stok.toktype: %v\n", stok.toktype)
}
