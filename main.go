package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"unicode"
)

const HEAPSIZE = 10000000
const NIL = 0

var heap = make(cells, 0, 10000000)

// TODO fix variable names
var ep int        //environment pointer
var hp int        //heap pointer
var sp int        //stack pointer
var fc = HEAPSIZE //free counter
var ap int        //arglist pointer

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

type subr func(argList int) int

// var subrPlus subr = func(argList int) int { return 0 }

func bindfunc(name string, tag tag, f subr) {
	sym := makesym(name)
	val := freshcell()
	heap.setTag(val, tag)
	heap.setSubr(val, f)
	heap.setCDR(val, 0)
	bindsym(sym, val)
}

// TODO Inpl
func initsubr() {
	m := map[string]subr{}

	for symname, f := range m {
		bindfunc(symname, SUBR, f)
	}

}

func bindsym(sym int, val int) {
	addr := assoc(sym, ep)
	if addr == 0 {
		assocsym(sym, val)
		return
	}
	heap.setCDR(addr, val)
}

type cell struct {
	tag  tag
	flag flag
	name string
	val  struct {
		num  int
		bind int
		subr subr
	}
	car int
	cdr int
}

type cells []cell

func (cs cells) getCAR(addr int) int           { return cs[addr].car }
func (cs cells) setCAR(addr int, x int)        { cs[addr].car = x }
func (cs cells) getCDR(addr int) int           { return cs[addr].cdr }
func (cs cells) setCDR(addr int, x int)        { cs[addr].cdr = x }
func (cs cells) getTag(addr int) tag           { return cs[addr].tag }
func (cs cells) setTag(addr int, t tag)        { cs[addr].tag = t }
func (cs cells) getName(addr int) string       { return cs[addr].name }
func (cs cells) setName(addr int, name string) { cs[addr].name = name }
func (cs cells) getNumber(addr int) int        { return cs[addr].val.num }
func (cs cells) setNumber(addr int, num int)   { cs[addr].val.num = num }
func (cs cells) getSymbol(addr int) int        { return cs[addr].val.num }
func (cs cells) setSubr(addr int, f subr)      { cs[addr].val.subr = f }
func (cs cells) isNumber(addr int) bool        { return cs[addr].tag == NUM }
func (cs cells) isSymbol(addr int) bool        { return cs[addr].tag == SYM }
func (cs cells) isList(addr int) bool          { return cs[addr].tag == LIS }
func (cs cells) isNIL(addr int) bool           { return cs[addr].tag == NIL }
func (cs cells) atomP(addr int) bool           { return cs.isNumber(addr) || cs.isSymbol(addr) }
func (cs cells) listP(addr int) bool           { return cs.isList(addr) || cs.isNIL(addr) }

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

func makenum(num int) int {
	addr := freshcell()
	heap.setTag(addr, NUM)
	heap.setNumber(addr, num)
	return addr
}

func freshcell() int {
	res := hp
	hp = heap.getCDR(hp)
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

func assoc(sym int, lis int) int {
	switch {
	case heap.isNIL(lis):
		return NIL
	case eqp(sym, heap.getCAR(heap.getCAR(lis))):
		return heap.getCAR(lis)
	}
	return assoc(sym, heap.getCDR(lis))
}

func eqp(addr1 int, addr2 int) bool {
	switch {
	case heap.isNumber(addr1) && heap.isNumber(addr2) && heap.getNumber(addr1) == heap.getNumber(addr2):
		return true
	case heap.isSymbol(addr1) && heap.isSymbol(addr2) && heap.getName(addr1) == heap.getName(addr2):
		return true
	}
	return false
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
		for c != EOL && c != SPACE && c != '(' && c != ')' {
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

func read() (int, error) {
	gettoken()
	switch stok.toktype {
	case NUMBER:
		num, err := strconv.Atoi(string(stok.buf[:len(stok.buf)-1]))
		if err != nil {
			return 0, err
		}
		return makenum(num), nil
	case SYMBOL:
		return makesym(string(stok.buf)), nil
	case QUOTE:
		addr, err := read()
		if err != nil {
			return 0, err
		}
		return cons(makesym("quote"), cons(addr, NIL)), nil
	case LPAREN:
		addr, err := readList()
		if err != nil {
			return 0, err
		}
		return addr, nil
	}
	return 0, errors.New("can't read")
}

func readList() (int, error) {
	gettoken()
	switch stok.toktype {
	case RPAREN:
		return NIL, nil
	case DOT:
		cdr, err := read()
		if err != nil {
			return 0, err
		}
		if heap.atomP(cdr) {
			gettoken()
		}
		return cdr, nil
	}
	stok.flag = BACK
	car, err := read()
	if err != nil {
		return 0, err
	}
	cdr, err := readList()
	if err != nil {
		return 0, err
	}
	return cons(car, cdr), nil
}

func print(addr int) {
	switch heap.getTag(addr) {
	case NUM:
		fmt.Printf("%d", heap.getNumber(addr))
	case SYM:
		fmt.Printf("%s", heap.getName(addr))
	case SUBR:
		fmt.Print("<subr>")
	case FSUBR:
		fmt.Print("<fsubr>")
	case FUNC:
		fmt.Print("<function>")
	case LIS:
		fmt.Print("(")
		printList(addr)
	}
}

func printList(addr int) {
	switch {
	case heap.isNIL(addr):
		fmt.Printf(")")
	case !heap.listP(heap.getCDR(addr)) && !heap.isNIL((heap.getCDR(addr))):
		fmt.Printf("%d . %d)", heap.getCAR(addr), heap.getCDR(addr))
	default:
		fmt.Print(heap.getCAR(addr))
		if !heap.isNIL(heap.getCDR(addr)) {
			fmt.Print(" ")
		}
		printList(heap.getCDR(addr))
	}
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}

func run() error {
	initcell()
	addr, err := read()
	if err != nil {
		return err
	}
	print(addr)
	// fmt.Printf("stok.buf: %v\n", stok.buf)
	// fmt.Printf("stok.toktype: %v\n", stok.toktype)
	return nil
}
