package main

import (
	"flag"
	"goformit/internal/logging"
	"goformit/pkg/goformit"
	"log"
	"os"
)

func main() {
	inputFile := flag.String("i", "", "JSON input file")
	outputfile := flag.String("o", "output.json", "JSON output")
	verbose := flag.Bool("v", false, "display logging upon program cleanup")
	flag.Parse()

	if len(*inputFile) == 0 {
		log.Fatal("no JSON input file")
	}

	form, err := goformit.NewFormFromJSON(*inputFile)
	if err != nil {
		log.Fatal(err)
	}

	err = form.Run()
	if err != nil {
		log.Fatal(err)
	}

	formResult, err := form.Result()
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(*outputfile)
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.Write([]byte(formResult))
	if err != nil {
		log.Fatal(err)
	}

	if *verbose {
		logging.AppLogger.Dump(os.Stdout)
	}
}
