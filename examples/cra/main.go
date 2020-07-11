package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/zserge/lorca"
)

func main() {
	ui, err := lorca.New("http://localhost:3000", "", 720, 480)
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()

	a := func() {
		fmt.Println("called")
	}

	ui.Bind("dispatchToServer", a)

	// Wait until the interrupt signal arrives or browser window is closed
	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)
	select {
	case <-sigc:
	case <-ui.Done():
	}

	log.Println("exiting...")
}
