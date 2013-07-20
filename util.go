package gopp

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
