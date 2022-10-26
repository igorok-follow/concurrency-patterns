package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
)

type File struct {
	Path    string
	Content string
}

func main() {
	channels := make([]<-chan *File, 100)

	for i := 0; i < 100; i++ {
		channels[i] = GetFile("./resources/dat" + strconv.Itoa(i))
	}

	files := fanIn(channels)
	for {
		fmt.Println(<-files)
	}
}

func GetFile(filePath string) <-chan *File {
	c := make(chan *File)

	go func() {
		content, err := ioutil.ReadFile(filePath)

		if err != nil {
			log.Fatal(err)
		}

		c <- &File{
			Path:    filePath,
			Content: string(content),
		}
	}()

	return c
}

func fanIn(inputs []<-chan *File) <-chan *File {
	c := make(chan *File)

	for i := range inputs {
		input := inputs[i]
		go func() { c <- <-input }()
	}

	return c
}
