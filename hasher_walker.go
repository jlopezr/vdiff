package main

import (
	"path/filepath"
	"time"
)

type HasherWalker struct {
	rootDirInfo *DirInfo
	done        chan bool
	cancel      chan bool
	msg         chan HashResult
}

type HashResult struct {
	dirInfo   *DirInfo
	name      string
	leftHash  string
	rightHash string
}

func NewHasherWalker() HasherWalker {
	done := make(chan bool, 1)
	cancel := make(chan bool, 1)
	msg := make(chan HashResult, 5)
	h := HasherWalker{
		done:   done,
		cancel: cancel,
		msg:    msg,
	}
	return h
}

func (h *HasherWalker) Walk(dirInfo *DirInfo) {
	h.rootDirInfo = dirInfo
	h.ProcessDirectory(dirInfo)
	h.done <- true
}

func (h *HasherWalker) ProcessDirectory(dirInfo *DirInfo) (string, string) {
	var directoryLeftHash string
	var directoryRightHash string

	//TODO si el status ya es DIFFERENT no hace falta calcular hash

	for _, f := range dirInfo.Files {
		time.Sleep(1 * time.Second)

		var leftHash string
		var rightHash string

		switch f.Left.Type {
		case DIRECTORY:
			leftHash, rightHash = h.ProcessDirectory(f.Info)
			break
		case FILE:
			leftHash, _ = hashFileMd5(filepath.Join(dirInfo.LeftPath,f.Name))
			break
		case SYMLINK:
		case ERROR:
		case ERROR_FILE:
		case ERROR_DIRECTORY:
		case ERROR_SYMLINK:
			//TODO Habra que ir pillando el hash de las cosas que ya estan calculadas, i.e symlink
			break
		}

		switch f.Right.Type {
		case DIRECTORY:
			if f.Left.Type != DIRECTORY {
				_, rightHash = h.ProcessDirectory(f.Info)
			}
			break
		case FILE:
			rightHash, _ = hashFileMd5(filepath.Join(dirInfo.RightPath,f.Name))
			break
		case SYMLINK:
		case ERROR:
		case ERROR_FILE:
		case ERROR_DIRECTORY:
		case ERROR_SYMLINK:
			//TODO Habra que ir pillando el hash de las cosas que ya estan calculadas, i.e symlink
			break
		}

		//TODO Si es un Updir los dos hash son "" y no se deberia enviar

		result := HashResult{
			dirInfo:   dirInfo,
			name:      f.Name,
			leftHash:  leftHash,
			rightHash: rightHash,
		}
		h.msg <- result

		directoryLeftHash = leftHash // TODO Calculate accumulated hash for directory
		directoryRightHash = rightHash // TODO Calculate accumulated hash for directory
	}
	return directoryLeftHash, directoryRightHash
}