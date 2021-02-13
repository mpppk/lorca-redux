package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	_ "github.com/mpppk/lorca-fsa/examples/cra/statik"
	"github.com/rakyll/statik/fs"

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

func readDirRequestHandler(action *fsa.Action, dispatch fsa.Dispatch) error {
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

func newServer(port int) (*http.Server, error) {
	statikFS, err := fs.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize html fs: %w", err)
	}
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(statikFS))
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}, nil
}

func main() {
	devMode := false
	if len(os.Args) > 1 && os.Args[1] == "dev" {
		devMode = true
	}

	server, err := newServer(3000)
	panicIfErr(err)

	go func() {
		panicIfErr(server.ListenAndServe())
	}()

	handlers := fsa.NewHandlers()
	handlers.Handle("APP/CLICK_READ_DIR_BUTTON", fsa.HandlerFunc(readDirRequestHandler))

	config := &fsa.LorcaConfig{
		AppName:          "lorca-cra-sample",
		Url:              "http://localhost:3000",
		Width:            720,
		Height:           480,
		EnableExtensions: devMode,
		Handlers:         handlers,
	}

	ui, err := fsa.Start(config)
	panicIfErr(err)

	defer func() {
		panicIfErr(ui.Close())
	}()
	fsa.Wait(ui)
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
