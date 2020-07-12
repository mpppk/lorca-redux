package lorcafsa

import (
	"encoding/json"
	"fmt"

	"github.com/zserge/lorca"
)

type Action struct {
	Type    Type        `json:"type"`
	Payload interface{} `json:"payload"`
	Error   bool        `json:"error"`
	Meta    interface{} `json:"meta"`
}

type Type string

type Dispatch func(action *Action) error
type Dispatcher interface {
	Dispatch(action *Action) error
}

type Handler interface {
	Do(action *Action, dispatch Dispatch) error
}

type LorcaDispatcher struct {
	ui lorca.UI
}

func NewLorcaDispatcher(ui lorca.UI) *LorcaDispatcher {
	return &LorcaDispatcher{ui: ui}
}

func (l *LorcaDispatcher) Dispatch(action *Action) error {
	if err := l.callHandleServerAction(action); err != nil {
		return fmt.Errorf("failed to dispatch action: %v: %w", action, err)
	}
	return nil
}

func (l *LorcaDispatcher) callHandleServerAction(action *Action) error {
	actionJson, err := json.Marshal(action)
	if err != nil {
		return fmt.Errorf("failed to marshal action: %v: %w", action, err)
	}
	js := fmt.Sprintf("handleServerAction(JSON.parse(`%s`))", actionJson)
	l.ui.Eval(js)
	return nil
}

type HandlerFunc func(action *Action, dispatch Dispatch) error

func (h HandlerFunc) Do(action *Action, dispatch Dispatch) error {
	return h(action, dispatch)
}

type Handlers struct {
	handlers   map[Type]Handler
	dispatcher Dispatcher
}

func (h *Handlers) SetDispatcher(dispatcher Dispatcher) {
	h.dispatcher = dispatcher
}

func NewHandlers() *Handlers {
	return &Handlers{
		handlers: map[Type]Handler{},
	}
}

func (h *Handlers) Handle(t Type, handler Handler) {
	h.handlers[t] = handler
}

func (h *Handlers) Dispatch(action *Action) error {
	if handler, ok := h.handlers[action.Type]; ok {
		return handler.Do(action, h.dispatcher.Dispatch)
	}
	return nil
}
