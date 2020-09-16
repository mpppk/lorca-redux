package lorcafsa

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/mpppk/lorca"
)

type LorcaConfig struct {
	AppName            string
	Url                string
	UserDir            string
	Width              int
	Height             int
	EnableExtensions   bool
	CustomArgs         []string
	Handlers           *Handlers
	DispatchMethodName string
	Logger             *log.Logger
}

func Start(config *LorcaConfig) (lorca.UI, error) {
	if config.Logger != nil {
		config.Handlers.SetLogger(config.Logger)
	}

	configDir := config.UserDir
	if config.UserDir == "" {
		cd, err := getConfigDir(config.AppName)
		if err != nil {
			return nil, err
		}
		configDir = cd
	}
	if config.EnableExtensions {
		enableChromeExtensions()
	}

	ui, err := lorca.New("", configDir, config.Width, config.Height)
	if err != nil {
		return nil, fmt.Errorf("failed to create new lorca UI: %w", err)
	}

	dispatcher := NewLorcaDispatcher(ui)
	if config.Logger != nil {
		dispatcher.SetLogger(config.Logger)
	}
	config.Handlers.SetDispatcher(dispatcher)

	dispatchMethodName := "dispatchToServer"
	if config.DispatchMethodName != "" {
		dispatchMethodName = config.DispatchMethodName
	}

	dispatch := func(action *Action) {
		if err := config.Handlers.Dispatch(action); err != nil && config.Logger != nil {
			config.Logger.Printf("warn: failed to dispatch action: error(%v): action(%#v)\n", err, action)
		}
	}

	if err := ui.Bind(dispatchMethodName, dispatch); err != nil {
		if err := ui.Close(); err != nil {
			panic(err)
		}
		return nil, fmt.Errorf("failed to bind dispatch method to dispatchToServer: %w", err)
	}

	if err := ui.Load(config.Url); err != nil {
		if err := ui.Close(); err != nil {
			panic(err)
		}
		return nil, fmt.Errorf("failed to load URL from %s: %w", config.Url, err)
	}
	return ui, nil
}

func Wait(ui lorca.UI) {
	// Wait until the interrupt signal arrives or browser window is closed
	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)
	select {
	case <-sigc:
	case <-ui.Done():
	}
}

func enableChromeExtensions() {
	var defaultArgs []string
	for _, arg := range lorca.DefaultChromeArgs {
		if arg != "--disable-extensions" {
			defaultArgs = append(defaultArgs, arg)
		}
	}
	lorca.DefaultChromeArgs = defaultArgs
}

func getConfigDir(appName string) (string, error) {
	configRootDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user config dir path: %w", err)
	}
	return filepath.Join(configRootDir, appName, "Chrome"), nil
}
