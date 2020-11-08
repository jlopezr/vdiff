package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

func hashFileMd5(filePath string) (string, error) {
	//Initialize variable returnMD5String now in case an error has to be returned
	var returnMD5String string

	//Open the passed argument and check for any error
	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String, err
	}

	//Tell the program to call the following function when the current function returns
	defer file.Close()

	//Open a new hash interface to write to
	hash := md5.New()

	//Copy the file in the hash interface and check for any error
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}

	//Get the 16 bytes hash
	hashInBytes := hash.Sum(nil)[:16]

	//Convert the bytes to a string
	returnMD5String = hex.EncodeToString(hashInBytes)

	return returnMD5String, nil
}

func main() {
	excludePtr := flag.String("exclude", "", "regex that matches all files to be excluded")
	flag.Parse()

	var exclude bool
	var myr *regexp.Regexp

	if *excludePtr != "" {
		exclude = true
		myr, err := regexp.Compile(*excludePtr)
		if err != nil {
			panic(err)
		}
	} else {
		exclude = false
	}

	var files []string

	root := "."
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		//Avoid excluded files
		if exclude {
			if myr.MatchString(path) {
				return nil
			}
		}

		//Avoid directories
		if info.IsDir() {
			fmt.Println(path, "is a dir")
			return nil
		}

		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		hash, err := hashFileMd5(file)

		if err != nil {
			panic(err)
		}

		fmt.Println(file, hash)
	}
}
