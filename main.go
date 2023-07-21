package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/momaek/formattag/align"
)

var version string

func main() {
	var (
		file           string
		showVersion    bool
		writeToConsole bool
	)

	flag.StringVar(&file, "file", "", "input data")
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.BoolVar(&writeToConsole, "C", false, "Write result to console")
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if showVersion {
		fmt.Println("Version:", version)
		return
	}

	if len(file) > 0 {
		align.Init(file)
	}

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		reader := bufio.NewReader(os.Stdin)
		align.Init(reader)
		writeToConsole = true
	}

	b, err := align.Do()
	if err != nil {
		log.Fatal("align failed ", err)
	}

	if writeToConsole {
		fmt.Println(string(b))
		return
	}

	ioutil.WriteFile(file, b, 0)
}
