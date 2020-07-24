package main

import (
	go_v8 "github.com/pefish/go-v8"
	"log"
)

func main() {
	v8Isolate := go_v8.NewV8Isolate()
	defer v8Isolate.Close()
	err := v8Isolate.RunInNewContext(`console.log("haha")`)
	if err != nil {
		log.Fatal(err)
	}
}
