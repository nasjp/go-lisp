package main

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

func (cs cells) getCAR(addr int) int {
	return cs[addr].car
}

func (cs cells) setCAR(addr int, x int) {
	cs[addr].car = x
}

func (cs cells) setCDR(addr int, x int) {
	cs[addr].cdr = x
}

func (cs cells) setTag(addr int, t tag) {
	cs[addr].tag = t
}

func (cs cells) setName(addr int, name string) {
	heap[addr].name = name
}

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

func main() {
	initcell()
}
