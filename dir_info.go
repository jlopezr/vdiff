package main

import "fmt"

type GenericInfo interface {
	Hash() string
	EntryType() string
	Print()
}

const (
	FILE      int = 0
	DIRECTORY int = 1
)

const (
	NOT_CHECKED_YET int = -1
	EQUALS          int = 0
	DIFFERENT       int = 1
	MISSING_RIGHT   int = 2
	EXTRA_RIGHT     int = 3
)

type EntryInfo struct {
	name  string
	left  *EntryInfo
	right *EntryInfo
	state int
}

type DirInfo struct {
	leftPath  string
	rightPath string
	children  []*EntryInfo
}

func (d DirInfo) EntryType() int {
	return DIRECTORY
}

func (d DirInfo) Print() {
	fmt.Println(d.leftPath + "\t\t\t" + d.rightPath)
	for i, s := range d.children {
		fmt.Println(i, s.name, s.state)
	}
}

func (d *DirInfo) AppendEntry(name string) {
	e := EntryInfo{
		name:  name,
		state: NOT_CHECKED_YET,
	}

	d.children = append(d.children, &e)
}

type FileInfo struct {
}

func (d FileInfo) EntryType() int {
	return FILE
}

func (d FileInfo) Print() {
	fmt.Println("*TODO*")
}

func Prueba() {
	d1 := DirInfo{
		leftPath:  "/Users/juan/a",
		rightPath: "/Users/juan/b",
		children:  make([]*EntryInfo, 0),
	}

	d1.AppendEntry("juan")
	d1.AppendEntry("toni")

	d1.Print()
}
