# lorca-fsa

lorca-fsa is a minimal library developed for [lorca](https://github.com/zserge/lorca) for message passing between client side and server side using Flux Standard Action (FSA).
It can be used in all applications that use FSA, regardless of which View library you use, such as React, Angular, Vue, etc.

lorca provides a way to call JavaScript from Go and vice versa via Chrome DevTools protocol.
Thanks to this, an application using lorca does not need to implement endpoints using REST API or GraphQL for communication between the front end and the server-side.
However, it is not realistic to bind many methods because the JavaScript methods called from Go must be exposed to global.
lorca-fsa provides an FSA-based message passing mechanism through a single global JavaScript method.
The action dispatched on the frontend will be sent to the server-side, and Go application can handle it.
The server-side processing result is also notified to the frontend by dispatching an action.

## Current Status: Under Development

## Installation

```shell
$ go get lorca-fsa
```

## Usage
```go
package main
import (
	"fmt"
	"io/ioutil"
	"os"

	fsa "github.com/mpppk/lorca-fsa"
)

func panicIfErrExist(err error) {
    if err != nil {
        panic(err)
    }
}

func newServerResponseAction() *fsa.Action {
	return &fsa.Action{ Type: "SERVER/PONG" }
}

// someActionHandler is action handler which will be registered by handlers.Handle
// fsa.Action represents fsa, so the struct have Type, Payload, Error, and Meta.
// fsa.Dispatch is dispatcher to dispatch action to frontend application.
func someActionHandler(action *fsa.Action, dispatch fsa.Dispatch) error {
    // payload is interface{}
    payload := action.Payload.(string)
    fmt.Println(payload)

    // action is dispatched to frontend application
	return dispatch(newServerResponseAction())
}

func main() {
	handlers := fsa.NewHandlers()

	// you can register action handler like http.Handle
	handlers.Handle("APP/CLICK_SOME_BUTTON", fsa.HandlerFunc(someActionHandler))

	config := &fsa.LorcaConfig{
		Width:            720,
		Height:           480,
		Handlers:         handlers,
	}

    // fsa.Start start lorca application and bind Go and JavaScript methods for dispatch
	ui, err := fsa.Start(config)
    panicIfErrExist(err)
	defer func() {
        panicIfErrExist(ui.Close())
	}()

    // fsa.Wait block goroutine until close UI or dispatch SIGINT signal
	fsa.Wait(ui)
}
```

See examples for more information.
* create-react-app + redux-toolkit sample
* Nextjs + redux sample
* Vue + redux sample
* Vue + Vuex sample
* Angular + ? sample

### Note
`fsa.Start` and `fsa.Wait` is just simple wrapper for lorca.New.
If you want to handle plain lorca.UI instance, see its implementation.

## redux-lorca

redux-lorca is the redux middleware to send dispatched action to server-side.

```shell script
$ yarn add redux-lorca
```

This is usage with redux-toolkit.

```js
import { configureStore } from '@reduxjs/toolkit';
import {makeLorcaMiddleware, setupServerActionHandler} from 'redux-lorca';

const makeStore = ({reducer}) => {
  const store =  configureStore({reducer, middleware: [makeLorcaMiddleware()]});
  setupServerActionHandler(store);
  return store;
}
```
