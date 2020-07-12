package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"

	fsa "github.com/mpppk/lorca-fsa"

	"github.com/zserge/lorca"
)

type ReadDirPayload struct {
	Dir   string   `json:"dir"`
	Files []string `json:"files"`
}

func newHandlers(ui lorca.UI) *fsa.Handlers {
	handlers := fsa.NewLorcaHandlers(ui)
	newReadDirAction := func(dir string, files []string) *fsa.Action {
		return &fsa.Action{
			Type: "SERVER/READ_DIR",
			Payload: ReadDirPayload{
				Dir:   dir,
				Files: files,
			},
			Error: false,
			Meta:  nil,
		}
	}

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
	ui, err := lorca.New("", "", 720, 480)
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()

	handlers := newHandlers(ui)
	ui.Bind("dispatchToServer", handlers.Dispatch)

	ui.Load("http://localhost:3000")

	// Wait until the interrupt signal arrives or browser window is closed
	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)
	select {
	case <-sigc:
	case <-ui.Done():
	}

	log.Println("exiting...")
}
