package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"time"
)

func createDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func main() {

	start := time.Now()

	var src string
	var dest string
	flag.StringVar(&src, "s", "src", "源目录")
	flag.StringVar(&dest, "d", "dest", "目标目录")
	flag.Parse()

	files, err := ioutil.ReadDir(src)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		ext := filepath.Ext(file.Name())

		if ext != ".rar" && ext != ".zip" {
			continue
		}
		re, _ := regexp.Compile("[0-9]{3}")
		bytes := re.Find([]byte(file.Name()))
		rename := string(bytes)
		folder := path.Join(dest, rename)
		createDirIfNotExist(folder)
		var cmd *exec.Cmd
		if ext == ".rar" {
			cmd = exec.Command("unrar", "e", path.Join(src, file.Name()), folder)
		} else {
			cmd = exec.Command("ditto", "-x", "-k", path.Join(src, file.Name()), folder)
		}
		err := cmd.Run()
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		childFiles, err1 := ioutil.ReadDir(folder)
		if err1 != nil {
			panic(err1)
		}

		for _, childFile := range childFiles {
			if childFile.IsDir() {
				err = filepath.Walk(path.Join(folder, childFile.Name()), func(p string, info os.FileInfo, err error) error {
					extName := filepath.Ext(p)
					newName := path.Join(folder, info.Name())
					if extName == ".jpg" || extName == ".jpeg" || extName == ".png" {
						os.Rename(p, newName)
					}
					return nil
				})
				os.RemoveAll(path.Join(folder, childFile.Name()))
			}

		}
	}
	t := time.Now()
	elapsed := t.Sub(start)
	println(elapsed.Milliseconds())
}
