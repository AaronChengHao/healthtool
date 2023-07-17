package main

import (
	"fmt"
	"github.com/flopp/go-findfont"
	"io/ioutil"

	"github.com/golang/freetype/truetype"
)

func main() {
	fontPath, err := findfont.Find("arial.ttf")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Found 'arial.ttf' in '%s'\n", fontPath)

	// load the font with the freetype library
	fontData, err := ioutil.ReadFile(fontPath)
	if err != nil {
		panic(err)
	}
	_, err = truetype.Parse(fontData)
	if err != nil {
		panic(err)
	}
	// use the font...
}
