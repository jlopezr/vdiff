package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

/*
FUNCIONES QUE FALTAN:
X si hay un error de I/O o permisos que muestre la linea y el problema pero continue
X que puedas poner multiples expressiones para excluir ficheros
X que puedas decir que no siga los links simbolicos
- apañar el check para que funcione con links simbolicos y multiples exclusiones

Modules structure
https://golang.org/doc/code.html
Do not defer on close
https://www.joeshaw.org/dont-defer-close-on-writable-files/
Other TUI libs
https://appliedgo.net/tui/
https://github.com/nsf/termbox-go
https://github.com/gdamore/tcell
https://github.com/jroimartin/gocui
https://github.com/briandowns/spinner
https://github.com/logrusorgru/aurora
Data structures
https://ieftimov.com/post/golang-datastructures-trees/
https://reinkrul.nl/blog/go/golang/merkle/tree/2020/05/21/golang-merkle-tree.html
Optimize binary size
https://blog.filippo.io/shrink-your-go-binaries-with-this-one-weird-trick/
Enums
http://golang-basic.blogspot.com/2014/07/step-by-step-guide-to-declaring-enums.html

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

type FileAndInfo struct {
	file string
	info os.FileInfo
}

func generate(exclusions ArrayFlags) {
	// print header for later checking
	var regex *FileExclusions
	if exclusions != nil {
		fmt.Println("Options: -exclude", exclusions.String())
		regex = CreateFileExclusions(exclusions)
	}

	var files []FileAndInfo

	root := "."
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		//Avoid excluded files
		if regex != nil {
			if regex.MatchString(path) {
				return nil
			}
		}

		//Avoid directories
		if info.IsDir() {
			return nil
		}

		files = append(files, FileAndInfo{
			file: path,
			info: info,
		})
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, f := range files {

		var hash string
		var err error

		//Treat symlinks
		if f.info.Mode()&os.ModeSymlink == os.ModeSymlink {
			hash, err = os.Readlink(f.file)
			hash = fmt.Sprintf("-> %s", hash)
		} else {
			hash, err = hashFileMd5(f.file)
		}

		if err != nil {
			hash = fmt.Sprintf("*ERROR* %s", err.Error())
			//panic(err)
		}

		fmt.Println(f.file, hash)
	}
}

type FileExclusions struct {
	regex []regexp.Regexp
}

func (f *FileExclusions) MatchString(txt string) bool {
	for _, r := range (*f).regex {
		if r.MatchString(txt) {
			return true
		}
	}

	return false
}

func CreateFileExclusions(a ArrayFlags) *FileExclusions {
	result := &FileExclusions {}

	for _, f := range a {
		//replace all . -> \.
		tmp := strings.ReplaceAll(f, ".", `\.`)
		//replace all * -> .*
		tmp = strings.ReplaceAll(tmp, "*", `.*`)

		regex, err := regexp.Compile(tmp)
		if err != nil {
				panic(err)
			}

		result.regex = append(result.regex, *regex)
	}

	return result
}

type ArrayFlags []string

func (a *ArrayFlags) String() string {
	return strings.Join(*a,",")
}

func (a *ArrayFlags) Set(value string) error {
	*a = append(*a, value)
	return nil
}

func OLDmain() {
	/*
	outputPtr := flag.String("output", "", "file where to store the results")
	flag.StringVar(outputPtr, "o", "", "file where to store the results")
	checkPtr := flag.String("check", "", "file to check")
	var exclusions ArrayFlags
	flag.Var(&exclusions, "e", "regex that matches all files to be excluded. Can be set multiple times.")
	flag.Parse()

	if *checkPtr != "" {
		check(*checkPtr)
	} else {
		generate(exclusions)
	}
*/
	Prueba()
}
