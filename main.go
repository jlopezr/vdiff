package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

/*
FUNCIONES QUE FALTAN:
- si hay un error de I/O o permisos que muestre la linea y el problema pero continue
- que puedas poner multiples expressiones para excluir ficheros
- que puedas decir que no siga los links simbolicos

https://golang.org/doc/code.html

https://appliedgo.net/tui/
https://github.com/nsf/termbox-go
https://github.com/gdamore/tcell
https://github.com/jroimartin/gocui
https://github.com/briandowns/spinner
https://github.com/logrusorgru/aurora

https://ieftimov.com/post/golang-datastructures-trees/
https://reinkrul.nl/blog/go/golang/merkle/tree/2020/05/21/golang-merkle-tree.html

COLORES:
unknown		gris
same		negro
orphan		lila
older		gris
newer		rojo
different	rojo

||[ ] filename    | size | modified |   ||[ ] filename    | size | modified ||

Inspiracion:
https://naarakstudio.com/direqual/index.html

*/

func checkFile(fileName string, hash1 string, path string) {
	hash2, err := hashFileMd5(path)
	if err != nil {
		panic(err)
	}

	if hash1 != hash2 {
		fmt.Println("[DIFFERENCES ]", fileName)
	} else {
		fmt.Println(fileName, "OK !")
	}
}

func check(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	textScanner := bufio.NewScanner(file)

	if textScanner.Scan() == false {
		panic("[Error] File to check is empty :(")
	}
	header := textScanner.Text()

	if strings.HasPrefix(header, "Options") {
		fmt.Println("HEADER:", header)
	}

	//TODO get exclude from header https://gobyexample.com/command-line-subcommands
	exclude := ".git"

	//replace all . -> \.
	exclude = strings.ReplaceAll(exclude, ".", `\.`)

	//replace all * -> .*
	exclude = strings.ReplaceAll(exclude, "*", `.*`)

	var excludeRegex *regexp.Regexp

	if exclude != "" {
		var err error
		excludeRegex, err = regexp.Compile(exclude)
		if err != nil {
			panic(err)
		}
	}

	scanner := NewFileScanner(textScanner)
	scanner.Scan()

	root := "."
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		//Avoid excluded files
		if excludeRegex != nil {
			if excludeRegex.MatchString(path) {
				return nil
			}
		}

		//Avoid directories
		if info.IsDir() {
			fmt.Println(path, "is a dir")
			return nil
		}

		//LEO UN ARCHIVO DEL FICHERO
		//caso 1. es el mismo archivo => compruebo el hash
		//caso 2. el nombre fichero < nombre disco  =>
		//			archivo que falta en disco =>
		//			itero del fichero hasta encontrar el mismo (y compruebo su hash) o posterior
		//caso 3. el nombre fichero > nombre disco  =>
		//			archivo que sobra en disco =>
		//          itero del disco (next) pero he de mantener el ultimo leido del fichero

		//TODO Como se si ya esta en EOF? => HasMore
		//if scanner.Text() != scanner.EOF {
		fileName := scanner.FileName()
		hash := scanner.Hash()

		/*
			fmt.Println("==============================")
			fmt.Println("FICHERO:", fileName)
			fmt.Println("DISCO  :", path)
			fmt.Println("------------------------------")
		*/

		if fileName == path {
			// CASE 1. The same file in both file and disk
			checkFile(fileName, hash, path)
			scanner.Scan() // Advances to next file
			return nil
		} else if fileName < path {
			// CASE 2. Missing file in disk
			for fileName <= path {
				if fileName < path {
					fmt.Println("[MISSING FILE]", fileName)
				} else {
					checkFile(fileName, hash, path)
				}
				scanner.Scan() //TODO check eof
				fileName = scanner.FileName()
			}
		} else {
			// CASE 3. Extra file in disk
			fmt.Println("[EXTRA   FILE]", path)
		}
		//}

		return nil
	})
	if err != nil {
		panic(err)
	}
	if scanner.HasMore() {
		fmt.Println("[MISSING FILE]", scanner.FileName())
		for scanner.Scan() {
			fmt.Println("[MISSING FILE]", scanner.FileName())
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func generate(exclude string) {
	// print header for later checking
	if exclude != "" {
		fmt.Println("Options: -exclude", exclude)
	}

	//replace all . -> \.
	exclude = strings.ReplaceAll(exclude, ".", `\.`)

	//replace all * -> .*
	exclude = strings.ReplaceAll(exclude, "*", `.*`)

	var excludeRegex *regexp.Regexp

	if exclude != "" {
		var err error
		excludeRegex, err = regexp.Compile(exclude)
		if err != nil {
			panic(err)
		}
	}

	var files []string

	root := "."
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		//Avoid excluded files
		if excludeRegex != nil {
			if excludeRegex.MatchString(path) {
				return nil
			}
		}

		//Avoid directories
		if info.IsDir() {
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

func Xmain() {

	excludePtr := flag.String("exclude", "", "regex that matches all files to be excluded")
	outputPtr := flag.String("output", "", "file where to store the results")
	flag.StringVar(outputPtr, "o", "", "file where to store the results")
	checkPtr := flag.String("check", "", "file to check")
	flag.Parse()

	if *checkPtr != "" {
		check(*checkPtr)
	} else {
		generate(*excludePtr)
	}

	//Prueba()
}
