package gopp

import (
	"fmt"
)

type literalSorter []string

func (l literalSorter) Len() int {
	return len(l)
}

func (l literalSorter) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l literalSorter) Less(i, j int) bool {
	if len(l[i]) > len(l[j]) {
		return true
	}
	if len(l[i]) < len(l[j]) {
		return false
	}
	return i < j
}

func printNode(node Node, indentCount int) {
	indent := func(tag string) {
		for i := 0; i < indentCount; i++ {
			fmt.Print(" ")
		}
		fmt.Println(tag)
	}
	switch node := node.(type) {
	case []Node:
		indent("[")
		for _, n := range node {
			printNode(n, indentCount+1)
		}
		indent("]")
	case AST:
		indent("[")
		for _, n := range node {
			printNode(n, indentCount+1)
		}
		indent("]")
	default:
		indent(fmt.Sprint(node))
	}
}