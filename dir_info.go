package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type EntryType uint8

const (
	UNKNOWN EntryType = iota
	FILE
	DIRECTORY
	SYMLINK
	NOT_EXIST
	ERROR
	ERROR_FILE
	ERROR_DIRECTORY
	ERROR_SYMLINK
)

var entries = [...]string{
	"UNKNOWN",
	"FILE",
	"DIRECTORY",
	"SYMLINK",
	"NOT_EXIST",
}

func (t EntryType) String() string {
	return entries[t]
}

type EntryState uint8

const (
	NOT_CHECKED_YET EntryState = iota
	EQUALS
	DIFFERENT
	NEWEST_LEFT
	NEWEST_RIGHT
	MISSING_LEFT
	MISSING_RIGHT
)

var states = [...]string{
	"NOT_CHECKED_YET",
	"EQUALS",
	"DIFFERENT",
	"NEWEST_LEFT",
	"NEWEST_RIGHT",
	"MISSING_LEFT",
	"MISSING_RIGHT",
}

func (t EntryState) String() string {
	return states[t]
}

type EntryInfo struct {
	Hash             string
	Size             int64
	LastModification time.Time
	Type             EntryType
}

type GenericInfo struct {
	Name  string
	Left  EntryInfo
	Right EntryInfo
	State EntryState
}

type DirEntry struct {
	GenericInfo
	Info *DirInfo
}

//Quiza podemos usar una estructura diferente si sabemos que ninguno de los dos lados es un directorio y ahorrar 1 puntero
type FileEntry struct {
	GenericInfo
}

type DirInfo struct {
	GenericInfo
	LeftPath  string
	RightPath string

	Files []*DirEntry
}

func (d DirInfo) Print() {
	d.PrintTab(0)
}

func (d DirInfo) PrintTab(level int) {
	tabs := strings.Repeat("\t", level)
	fmt.Println(tabs + d.LeftPath + "\t\t\t" + d.RightPath)
	for i, s := range d.Files {

		fmt.Println(tabs, i, s.Name, s.State, s.Left.Hash, s.Left.Type, s.Right.Hash, s.Right.Type)

		if s.Info != nil && (s.Left.Type == DIRECTORY || s.Right.Type == DIRECTORY) {
			s.Info.PrintTab(level + 1)
		}
	}
}

func (d *DirInfo) EntryCount() int {
	return len(d.Files)
}

func (d *DirInfo) GetEntry(i int) *DirEntry {
	return d.Files[i]
}

func (d *DirEntry) GetInfo(isLeft bool) *EntryInfo {
	if isLeft {
		return &d.Left
	} else {
		return &d.Right
	}
}
func (d *DirEntry) ConvertToDirectory(isLeft bool, parentDirInfo *DirInfo) {
	info := d.GetInfo(isLeft)
	info.Type = DIRECTORY

	if d.Info == nil {
		dirInfo := NewDirInfo()
		dirInfo.Name = d.Name //TODO Se podria ahorrar repetir el name
		dirInfo.LeftPath = fmt.Sprintf("%s%c%s", parentDirInfo.LeftPath, os.PathSeparator, d.Name)
		dirInfo.RightPath = fmt.Sprintf("%s%c%s", parentDirInfo.RightPath, os.PathSeparator, d.Name)
		d.Info = dirInfo
	}
}

func (d *DirInfo) FindEntry(name string) (int, *DirEntry) {
	for i, entry := range d.Files {
		if entry.Name == name {
			return i, entry
		}
	}
	return -1, nil
}

func (d *DirInfo) AppendEntry(name string) *DirEntry {
	e := DirEntry{}
	e.Name = name
	e.State = NOT_CHECKED_YET
	e.Left.Type = UNKNOWN
	e.Right.Type = UNKNOWN
	e.Left.Hash = fmt.Sprintf("HASH1-%s", name)
	e.Right.Hash = fmt.Sprintf("HASH2-%s", name)
	d.Files = append(d.Files, &e)
	return &e
}

func (d *DirInfo) AppendFile(name string) *DirEntry {
	e := DirEntry{}
	e.Name = name
	e.State = NOT_CHECKED_YET
	e.Left.Type = FILE
	e.Right.Type = UNKNOWN
	e.Left.Hash = fmt.Sprintf("HASH1-%s", name)
	d.Files = append(d.Files, &e)
	return &e
}

func (d *DirInfo) AppendDirectory(name string) *DirEntry {
	e := DirEntry{}
	e.Name = name
	e.State = NOT_CHECKED_YET
	e.Left.Type = DIRECTORY
	e.Right.Type = UNKNOWN
	e.Left.Hash = fmt.Sprintf("HASH1-%s", name)
	d.Files = append(d.Files, &e)

	dirInfo := DirInfo{}
	dirInfo.Name = name //TODO Se podria ahorrar repetir el name
	dirInfo.LeftPath = fmt.Sprintf("%s%c%s", d.LeftPath, os.PathSeparator, name)
	dirInfo.RightPath = fmt.Sprintf("%s%c%s", d.RightPath, os.PathSeparator, name)

	e.Info = &dirInfo
	return &e
}

func NewDirInfo() *DirInfo {
	d := DirInfo{
		Files: make([]*DirEntry, 0),
	}
	return &d
}

func Prueba() {
	d1 := NewDirInfo()
	d1.LeftPath = "/Users/juan/a"
	d1.RightPath = "/Users/juan/b"

	e3 := d1.AppendDirectory("dir1")
	d1.AppendEntry("file1")
	d1.AppendEntry("file2")

	e1 := d1.GetEntry(1)
	e1.State = EQUALS
	e1.Left.Type = FILE
	e1.Right.Type = FILE

	_, e2 := d1.FindEntry("file2")
	e2.State = DIFFERENT
	e1.Left.Type = FILE
	e1.Right.Type = SYMLINK

	e3.State = DIFFERENT
	d2 := e3.Info
	d2.AppendFile("a")
	d2.AppendFile("b")

	d1.Print()
}
