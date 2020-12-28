package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
)

type DirectoryWalker struct {
	fileExclusions   *FileExclusions
	rootDirInfo      *DirInfo
	//currentDirInfo   *DirInfo
	totalFiles       int
	totalDirectories int
}

func (d *DirectoryWalker) SetExclusions(exclusions ArrayFlags) {
	if exclusions != nil {
		d.fileExclusions = CreateFileExclusions(exclusions)
	}
}

func (d *DirectoryWalker) CreateRootDirInfo(directory string) bool {
	cwd, err := os.Getwd()
	if err != nil {
		panic("cannot get current directory")
	}

	if d.rootDirInfo == nil {
		d.rootDirInfo = NewDirInfo()
		d.rootDirInfo.Name = "."

		//TODO Check if it is directory
		d.rootDirInfo.Left.Type = DIRECTORY
		d.rootDirInfo.LeftPath = filepath.Join(cwd, directory)
		return true
	} else {
		if d.rootDirInfo.Right.Type != UNKNOWN {
			panic("Calling 3 times CreateRootDir")
		}
		//TODO Check if it is directory
		d.rootDirInfo.Right.Type = DIRECTORY
		d.rootDirInfo.RightPath = filepath.Join(cwd, directory)
		return false
	}
}

func (d *DirectoryWalker) ProcessDirectory(currentDirInfo *DirInfo, directory string, isLeft bool) {
	file, err := os.Open(directory)
	if err != nil {
		log.Fatalf("failed opening directory: %s", err)
	}
	defer file.Close()

	list, _ := file.Readdirnames(0) // 0 to read all files and folders
	for _, name := range list {
		fileName := filepath.Join(file.Name(), name)

		if d.fileExclusions.MatchString(fileName) {
			continue
		}

		info, err := os.Lstat(fileName)
		entry := currentDirInfo.FindOrAppendEntry(name)
		details := entry.GetInfo(isLeft)
		if err != nil {
			// Error
			details.Type = ERROR
			details.Hash = fmt.Sprintf("*ERROR* %s", err.Error())
		} else {
			if info.IsDir() {
				// Directory
				entry.ConvertToDirectory(isLeft, currentDirInfo)
				details.Hash = ""
				details.LastModification = info.ModTime()
				d.ProcessDirectory(entry.Info, path.Join(directory, name), isLeft)
			} else if info.Mode()&os.ModeSymlink == os.ModeSymlink {
				// Symlink
				link, err := os.Readlink(fileName)
				if err != nil {
					details.Type = ERROR_SYMLINK
					details.Hash = fmt.Sprintf("*ERROR* %s", err.Error())
				} else {
					details.Type = SYMLINK
					details.Hash = fmt.Sprintf("->%s", link)
					details.Size = info.Size()
					details.LastModification = info.ModTime()
				}
			} else {
				// File
				details.Type = FILE
				details.Hash = fmt.Sprintf("HASH-%s", fileName)
				details.Size = info.Size()
				details.LastModification = info.ModTime()
				//TODO Hash
			}
		}
	}
}

func (d *DirectoryWalker) Walk(directory string) {
	isLeft := d.CreateRootDirInfo(directory)
	d.ProcessDirectory(d.rootDirInfo, directory, isLeft)
}

func (d *DirectoryWalker) ProcessHashes(currentDirInfo *DirInfo) {
	allEqual := true
	for _,f := range currentDirInfo.Files {
		if f.State == NOT_CHECKED_YET {
			if f.Left.Type != f.Right.Type {
				f.State = DIFFERENT
				allEqual = false
			} else {
				switch f.Left.Type {
				case FILE:
					if f.Left.Size != f.Right.Size {
						f.State = DIFFERENT
						allEqual = false
					} else {
						//TODO Calculate Hash
						f.State = EQUALS
					}
				case SYMLINK:
					if f.Left.Hash != f.Right.Hash {
						f.State = DIFFERENT
						allEqual = false
					} else {
						f.State = EQUALS
					}
				case DIRECTORY:
					d.ProcessHashes(f.Info)
					f.State = f.Info.State
					if f.State != EQUALS {
						allEqual = false
					}
				}
			}
		}
	}
	if allEqual {
		currentDirInfo.State = EQUALS
	} else {
		currentDirInfo.State = DIFFERENT
	}
}

func Prueba2() *DirInfo {
	w := DirectoryWalker{}
	flags := ArrayFlags{}
	flags.Set(".git")
	w.SetExclusions(flags)
	w.Walk("./a")
	w.Walk("./b")
	w.ProcessHashes(w.rootDirInfo)
	w.rootDirInfo.Print()

	return w.rootDirInfo
}
