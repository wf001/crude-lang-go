package crudego

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

/*
Token
*/
type TokKind int

const (
	TK_KIND_RESERVED TokKind = iota + 1
	TK_KIND_NUM
	TK_KIND_EOF
)

type Token struct {
	Kind TokKind
	Next *Token
	Val  string
}

/*
Node
*/
type NodeKind int

const (
	ND_KIND_ADD NodeKind = iota + 1
	ND_KIND_SUB
	ND_KIND_MUL
	ND_KIND_DIV
	ND_KIND_NUM
)

type Node struct {
	Kind NodeKind
	Lhs  *Node
	Rhs  *Node
	Val  int
}

/*
Token Func
*/
func printToken(tok *Token) {
	if DEBUG {
		fmt.Println("==========================")
		for ; tok != nil; tok = tok.Next {
			Info("tok %p\n", tok)
			Info("%+v\n", tok)
		}
		fmt.Println("==========================")
	}
}
func Info(s string, v interface{}) {
	if DEBUG {
		_, file, line, _ := runtime.Caller(1)
		reg := "[/]"
		files := regexp.MustCompile(reg).Split(file, -1)
		fmt.Printf("%s %d| ", files[len(files)-1], line)
		fmt.Printf(s, v)
	}
}
func sep() {
	fmt.Println("--------------------------------")
}

func newToken(kind TokKind, cur *Token, val string) *Token {
	Info("kind: %d\n", kind)
	Info("cur: %+v\n", cur)
	Info("val: %s\n", val)
	tok := new(Token)
	tok.Kind = kind
	tok.Val = val
	cur.Next = tok
	return tok
}
func TokenizeHandler() *Token {
	flag.Parse()
	arg := flag.Arg(0)
	// space trim
	arg = strings.Replace(arg, " ", "", -1)
	// gen num arr
	reg := "[+-]"
	arg_arr := regexp.MustCompile(reg).Split(arg, -1)
	cur_len := len(arg_arr[0])

	head := new(Token)
	head.Next = nil
	head.Kind = -1
	cur := head

	cur = newToken(TK_KIND_NUM, cur, arg_arr[0])

	//tokenize
	for _, s := range arg_arr[1:] {
		op := string(arg[cur_len])
		if op == "+" || op == "-" {
			cur = newToken(TK_KIND_RESERVED, cur, string(arg[cur_len]))
			cur = newToken(TK_KIND_NUM, cur, s)
			cur_len += len(s) + 1
		}
	}
	cur = newToken(TK_KIND_EOF, cur, "")
	printToken(head)

	return head
}

func isNumber(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return i
}

/*
Node Func
*/
func newNode(kind NodeKind, lhs *Node, rhs *Node) *Node {
	Info("newNode %d\n", kind)
	node := new(Node)
	node.Kind = kind
	node.Lhs = lhs
	node.Rhs = rhs
	return node
}
func newNodeNum(val int) *Node {
	Info("newNodeNum %d\n", val)
	node := new(Node)
	Info("%p\n", node)
	node.Kind = ND_KIND_NUM
	node.Val = val
	return node
}

func printNode(node *Node) {
	if DEBUG {
		if node.Kind == ND_KIND_NUM {
			Info("node %p\n", node)
			Info("%+v\n", node)
			return
		}
		Info("node %p\n", node)
		Info("%+v\n", node)
		printNode(node.Lhs)
		printNode(node.Rhs)
	}
}
func Expr(tok *Token) *Node {
	var m_node *Node
	Info("%s\n", "expr")
	Info("%p\n", tok)
	tok, node := mul(tok)
	for {
		Info("%s\n", "expr for")
		Info("%s\n", tok.Val)
		if tok.Val == "+" {
			tok = tok.Next
			tok, m_node = mul(tok)
			node = newNode(ND_KIND_ADD, node, m_node)
		} else if tok.Val == "-" {
			tok = tok.Next
			tok, m_node = mul(tok)
			node = newNode(ND_KIND_SUB, node, m_node)
		} else {
			Info("%s\n", "=================")
			printNode(node)
			Info("%s\n", "=================")
			return node
		}
	}
}
func mul(tok *Token) (*Token, *Node) {
    var p_node *Node
	Info("%s\n", "mul")
	Info("%p\n", tok)
	tok, node := primary(tok)
	for {
		Info("%s\n", "mul for")
		Info("%s\n", tok.Val)
		if tok.Val == "*" {
			tok = tok.Next
			tok, p_node = primary(tok)
			node = newNode(ND_KIND_MUL, node, p_node)
		} else if tok.Val == "/" {
			tok = tok.Next
			tok, p_node = primary(tok)
			node = newNode(ND_KIND_DIV, node, p_node)
		} else {
			return tok, node
		}
	}
}
func primary(tok *Token) (*Token, *Node) {
	Info("%s\n", "pri")
	Info("%p\n", tok)
	Info("%s\n", tok.Val)
	if tok.Val == "(" {
		tok = tok.Next
		node := Expr(tok)
		if tok.Val != ")" {
			panic("error")
		}
		tok = tok.Next
		return tok, node
	}
	i := isNumber(tok.Val)
	tok = tok.Next
	return tok, newNodeNum(i)
}
