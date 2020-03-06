package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type tokenkind int

const (
	_ tokenkind = iota
	TK_RESERVED
	TK_NUM
	TK_EOL
)

type token struct {
	kind tokenkind
	next *token
	num  int
	str  string
}

func newToken(kind tokenkind, cur *token, str string) *token {
	t := &token{
		kind: kind,
		str:  str,
	}
	cur.next = t
	return t
}

// func isNumber(ch rune) bool {
//   n, err := strconv.Atoi(string(ch))

// }

func takeNum(str string) (int, int) {
	var numRunes []rune
	var steps int
	for _, ch := range str {
		if _, err := strconv.Atoi(string(ch)); err != nil {
			break
		}
		steps++
		numRunes = append(numRunes, ch)
	}

	n, err := strconv.Atoi(string(numRunes))
	if err != nil {
		return 0, 0
	}

	return n, steps
}

func tokenize(line string) *token {
	head := &token{}
	cur := head

	chars := []rune(line)

	for i := 0; i < len(chars); {
		if chars[i] == ' ' {
			i++
			continue
		}

		if chars[i] == '+' {
			cur = newToken(TK_RESERVED, cur, string(chars[i]))
			i++
			continue
		}

		if chars[i] == '-' {
			cur = newToken(TK_RESERVED, cur, string(chars[i]))

			i++
			continue
		}

		n, s := takeNum(line[i:])
		cur = newToken(TK_NUM, cur, line[i:i+s])
		cur.num = n
		i += s
	}
	return head.next
}

func gen(t *token) string {
	var txt string
	cur := t
	var count int

	var pre string
	for {
		if cur == nil {
			break
		}

		switch cur.kind {
		case TK_RESERVED:
			pre = cur.str
		case TK_NUM:
			switch pre {
			case "+":
				count += cur.num
			case "-":
				count -= cur.num
			}
		}

		txt += cur.str + " "
		cur = cur.next
	}

	txt += " = " + strconv.Itoa(count)
	return txt
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}

func run() error {
	sc := bufio.NewScanner(os.Stdin)
	if fmt.Print("> "); sc.Scan() {
		line := tokenize(sc.Text())
		out := gen(line)
		fmt.Println(out)
	}
	return nil
}
