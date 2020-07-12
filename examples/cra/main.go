package main

import (
	"fmt"
	"io/ioutil"
	"os"

	fsa "github.com/mpppk/lorca-fsa"
)

type ReadDirPayload struct {
	Dir   string   `json:"dir"`
	Files []string `json:"files"`
}

func newReadDirAction(dir string, files []string) *fsa.Action {
	return &fsa.Action{
		Type: "SERVER/READ_DIR",
		Payload: ReadDirPayload{
			Dir:   dir,
			Files: files,
		},
	}
}

func newHandlers() *fsa.Handlers {
	handlers := fsa.NewHandlers()

	readDirRequestHandler := func(action *fsa.Action, dispatch fsa.Dispatch) error {
		dir := action.Payload.(string)
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}

		var fileNames []string
		for _, file := range files {
			fileNames = append(fileNames, file.Name())
		}
		return dispatch(newReadDirAction(dir, fileNames))
	}
	handlers.Handle("APP/CLICK_READ_DIR_BUTTON", fsa.HandlerFunc(readDirRequestHandler))
	return handlers
}

func main() {
	devMode := false
	if len(os.Args) > 1 && os.Args[1] == "dev" {
		devMode = true
	}

	handlers := newHandlers()

	config := &fsa.LorcaConfig{
		AppName:          "lorca-cra-sample",
		Url:              "http://localhost:3000",
		Width:            720,
		Height:           480,
		EnableExtensions: devMode,
		Handlers:         handlers,
	}

	ui, err := fsa.Start(config)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := ui.Close(); err != nil {
			panic(err)
		}
	}()
	fsa.Wait(ui)
	fmt.Println("wait finish")
}
