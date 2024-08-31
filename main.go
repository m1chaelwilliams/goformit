package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"goformit/context"
	"goformit/logging"
	"goformit/models"
	"goformit/serialization"
)

func main() {
	inputFile := flag.String("i", "", "JSON input file")
	outputfile := flag.String("o", "output.json", "JSON output")
	verbose := flag.Bool("v", false, "display logging upon program cleanup")
	flag.Parse()

	if len(*inputFile) == 0 {
		log.Fatal("no JSON input file")
	}

	formJSON, err := serialization.NewFormJSON(*inputFile)
	if err != nil {
		log.Fatal(err)
	}

	ctx, err := context.NewAppContext(
		formJSON,
		func(promptJSON *serialization.PromptJSON) tea.Model {
			switch promptJSON.Type {
			case "input":
				return models.NewInputModelFromJSON(promptJSON)
			case "selection":
				return models.NewSelectionModelfromJSON(promptJSON)
			case "checkbox":
				return models.NewMultiSelectModelFromJSON(promptJSON)
			}
			return nil
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	program := tea.NewProgram(ctx.ActiveModel())
	_, err = program.Run()
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(*outputfile)
	if err != nil {
		log.Fatal(err)
	}

	contents, err := json.Marshal(ctx.FormResult())
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.Write(contents)
	if err != nil {
		log.Fatal(err)
	}

	if *verbose {
		logging.AppLogger.Dump(os.Stdout)
	}
}
