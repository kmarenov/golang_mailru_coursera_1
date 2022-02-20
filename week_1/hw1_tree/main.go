package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	lev := 0
	completed := make(map[int]bool)
	return tree(out, path, printFiles, lev, completed)
}

func tree(out io.Writer, path string, printFiles bool, level int, completed map[int]bool) error {
	level++
	completed[level] = false
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	items, _ := f.Readdir(-1)

	sort.Slice(items, func(i, j int) bool {
		return items[i].Name() < items[j].Name()
	})

	count := 0
	if printFiles {
		count = len(items)
	} else {
		for _, v := range items {
			if v.IsDir() {
				count++
			}
		}
	}

	n := 0
	for _, v := range items {
		if v.IsDir() || printFiles {
			str := ""
			for j := 1; j < level; j++ {
				if completed[j] {
					str += "\t"
				} else {
					str += "│\t"
				}
			}

			if n == count-1 {
				str += "└───"
				completed[level] = true
			} else {
				str += "├───"
			}

			str += v.Name()

			if !v.IsDir() && printFiles {
				if v.Size() > 0 {
					str += " (" + strconv.Itoa(int(v.Size())) + "b)"
				} else {
					str += " (empty)"
				}
			}

			str += "\n"
			_, err := out.Write([]byte(str))
			if err != nil {
				return err
			}

			n++
		}

		if v.IsDir() {
			err = tree(out, fmt.Sprintf("%v%c%v", path, os.PathSeparator, v.Name()), printFiles, level, completed)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
