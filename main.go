package main

import (
	"flag"
	"fmt"
	"github.com/programmer10110/gifTyper/typer"
	"image/gif"
	"log"
	"os"
)

func main() {
	outputGif := flag.String("output", "out.gif", "Output gif file")
	text := flag.String("text", "", "Text to type")

	flag.Parse()

	if *text == "" {
		fmt.Println("Usage: gifTyper -text=\"Hello World!\" -output=\"out.gif\"")
		return
	}

	generator, err := typer.InitGenerator()
	if err != nil {
		log.Fatalln(err.Error())
	}
	err = generator.SetDelay(1)
	if err != nil {
		log.Fatalln(err.Error())
	}
	err = generator.SetFontSize(32)
	if err != nil {
		log.Fatalln(err.Error())
	}
	textGif, err := generator.GenerateGIF(*text)
	if err != nil {
		log.Fatalln(err.Error())
	}
	f, err := os.OpenFile(*outputGif, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer f.Close()
	if err := gif.EncodeAll(f, textGif); err != nil {
		log.Fatalln(err.Error())
	}
}
