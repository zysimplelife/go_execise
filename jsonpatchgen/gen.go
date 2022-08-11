package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	jsonpatch "gomodules.xyz/jsonpatch/v2"
)

func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) < 2 {
		fmt.Println("jsonpatchgen [path base json] [path to new json]")
		return
	}

	oldJson := readJsonFromFile(argsWithoutProg[0])
	newJson := readJsonFromFile(argsWithoutProg[1])

	patch, e := jsonpatch.CreatePatch(oldJson, newJson)
	if e != nil {
		fmt.Printf("Error creating JSON patch:%v", e)
		return
	}
	for _, operation := range patch {
		fmt.Printf("%s\n", operation.Json())
	}
}

func readJsonFromFile(file string) (byteResult []byte) {
	fileContent, err := os.Open(file)

	if err != nil {
		log.Fatal(err)
		return
	}

	defer fileContent.Close()

	byteResult, _ = ioutil.ReadAll(fileContent)

	return
}
