package lorcafsa

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sync"

	"github.com/mpppk/lorca"
)

var emptyLogger = log.New(ioutil.Discard, "", 0)

// Action represents Flux Standard Action
type Action struct {
	Type    Type        `json:"type"`
	Payload interface{} `json:"payload"`
	Error   bool        `json:"error"`
	Meta    interface{} `json:"meta"`
}

// Type represents type property of Flux Standard Action
type Type string

// Dispatch is function for dispatch action
type Dispatch func(action *Action) error

// Dispatcher dispatch action
type Dispatcher interface {
	Dispatch(action *Action) error
}

// Handler handle action
type Handler interface {
	Do(action *Action, dispatch Dispatch) error
}

// LorcaDispatcher dispatch action via ui.Bind()
type LorcaDispatcher struct {
	ui     lorca.UI
	logger *log.Logger
}

// NewLorcaDispatcher generate LorcaDispatcher
func NewLorcaDispatcher(ui lorca.UI) *LorcaDispatcher {
	return &LorcaDispatcher{ui: ui, logger: emptyLogger}
}

// Dispatch dispatch action via ui.Eval()
func (l *LorcaDispatcher) Dispatch(action *Action) error {
	l.logger.Println("debug: action will be dispatched from server:", formatAction(action))
	if err := l.callHandleServerAction(action); err != nil {
		return fmt.Errorf("failed to dispatch action: %s: %w", formatAction(action), err)
	}
	return nil
}

// SetLogger sets logger
func (l *LorcaDispatcher) SetLogger(logger *log.Logger) {
	l.logger = logger
}

func (l *LorcaDispatcher) callHandleServerAction(action *Action) error {
	actionJson, err := json.Marshal(action)
	if err != nil {
		return fmt.Errorf("failed to marshal action: %s: %w", formatAction(action), err)
	}
	js := fmt.Sprintf("handleServerAction(JSON.parse(`%s`))", actionJson)
	l.logger.Println("debug: action will be dispatched to front as:", string(actionJson))
	l.ui.Eval(js)
	return nil
}

// The HandlerFunc type is an adapter to allow the use of ordinary functions as action handlers.
// If f is a function with the appropriate signature, HandlerFunc(f) is a Handler that calls f.
type HandlerFunc func(action *Action, dispatch Dispatch) error

// Do calls f(action, dispatch).
func (f HandlerFunc) Do(action *Action, dispatch Dispatch) error {
	return f(action, dispatch)
}

// Handlers represents set of handler.
type Handlers struct {
	sync.Mutex
	handlers   map[Type]Handler
	dispatcher Dispatcher
	logger     *log.Logger
}

// SetDispatcher sets dispatcher
func (h *Handlers) SetDispatcher(dispatcher Dispatcher) {
	h.dispatcher = dispatcher
}

// SetLogger sets logger
func (h *Handlers) SetLogger(logger *log.Logger) {
	h.logger = logger
}

// NewHandlers returns new Handlers
func NewHandlers() *Handlers {
	return &Handlers{
		handlers: map[Type]Handler{},
		logger:   emptyLogger,
	}
}

// Handle registers the action handler for the given type.
func (h *Handlers) Handle(t Type, handler Handler) {
	// FIXME: panic if handler already exists for type
	h.handlers[t] = handler
	h.logger.Println("debug: action handler is registered for ", t)
}

// Dispatch handle action by registered handlers
func (h *Handlers) Dispatch(action *Action) error {
	h.Lock()
	defer h.Unlock()
	if handler, ok := h.handlers[action.Type]; ok {
		h.logger.Println("debug: handler has found for action:", formatAction(action))
		err := handler.Do(action, h.dispatcher.Dispatch)
		h.logger.Println("debug: finish handler execution:", formatAction(action))
		return err
	} else {
		h.logger.Printf("debug: action handler does not found for type: %s\n", action.Type)
	}
	return nil
}

func formatAction(action *Action) string {
	msg := "Type: " + string(action.Type)
	if action.Payload != nil {
		msg += fmt.Sprintf(" Payload: %#v", action.Payload)
	}

	if action.Meta != nil {
		msg += fmt.Sprintf(" Meta: %v", action.Meta)
	}

	if action.Error {
		msg = " (Error)"
	}

	return msg
}
