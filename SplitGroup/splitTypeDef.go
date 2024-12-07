package main

type Set map[string]bool

func (set Set) exists(val string) bool {
	_, ok := set[val]
	if ok {
		return true
	} else {
		return false
	}
}

func (set Set) add(val string) {
	set[val] = true
}
