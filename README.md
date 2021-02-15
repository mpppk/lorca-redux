# lorca-fsa

lorca-fsa is a minimal library developed for [lorca](https://github.com/zserge/lorca) for message passing between client-side and server-side using _Flux Standard Action (FSA)_.
It can be used in all applications that use FSA, regardless of which View library you use, such as React, Angular, Vue, etc.

![](images/architecture.png)

## Motivation

lorca provides a way to call JavaScript from Go and vice versa via Chrome DevTools protocol as `ui.Eval` and `ui.Bind`.
Thanks to this, an application using lorca does not need to implement API like REST or GraphQL for communication between the client-side and the server-side.

However, using `ui.Eval` and` ui.Bind` in large-scale applications presents some challenges.

1. ui.Bind binds Go methods to JavaScript global variables. So when dozens or more Go methods are bound, global pollution becomes a problem.

2. Affinity with the state management system at the client. Bound methods are often responsible for processing with side effects such as file system access and HTTP requests. If the client already has a mechanism for state management, the bound methods should also be managed by that mechanism.

To solve these problems, lorca-fsa provides an FSA-based message passing mechanism through a single global JavaScript method.
The action dispatched on the frontend will be sent to the server-side, and Go application can handle it.
The server-side processing result is also notified to the frontend by dispatching an action.

## Current Status: Alpha

lorca-fsa is working fine on [my project](https://github.com/mpppk/imagine), but it needs more feedbacks. The APIs are not stable.

## Installation

```shell
$ go get github.com/mpppk/lorca-fsa
```

## Usage

```go
package main

import (
	"fmt"

	fsa "github.com/mpppk/lorca-fsa
)

func panicIfErrExist(err error) {
	if err != nil {
		panic(err)
	}
}

func newServerResponseAction() *fsa.Action {
	return &fsa.Action{ Type: "SERVER/PONG" }
}

// someActionHandler is action handler which will be registered by handlers.Handle.
// fsa.Action represents FSA, so the struct have Type, Payload, Error, and Meta.
// fsa.Dispatch is dispatcher for dispatching action to frontend.
func someActionHandler(action *fsa.Action, dispatch fsa.Dispatch) error {
	// payload is interface{}, so you need to cast to correct type.
	payload := action.Payload.(string)
	fmt.Println(payload)

	// dispatch method dispatch given action to client.
	return dispatch(newServerResponseAction())
}

func main() {
	handlers := fsa.NewHandlers()

	// you can register action handler like http.Handle
	handlers.Handle("APP/CLICK_SOME_BUTTON", fsa.HandlerFunc(someActionHandler))

	config := &fsa.LorcaConfig{
		Width:    720,
		Height:   480,
		Handlers: handlers,
	}

	// fsa.Start start lorca application and establish connection to exchange fsa
	ui, err := fsa.Start(config)
	panicIfErrExist(err)
	defer func() {
		panicIfErrExist(ui.Close())
	}()

	// Wait waits until the interrupt signal arrives or browser window is closed
	fsa.Wait(ui)
}
```

See examples for more information.

- [create-react-app + redux-toolkit sample](https://github.com/mpppk/lorca-fsa/cra)

### Note

`fsa.Start` and `fsa.Wait` are just simple wrapper for original lorca API.
If you want to handle plain lorca.UI instance, see [its implementation](https://github.com/mpppk/lorca-fsa/blob/master/lorca-fsa/util.go).

## Usage with Redux

You can integrate lorca-fsa and redux by write very small redux middleware.
By default, `window.dispatchToServer` is automatically injected, to dispatch fsa from the browser to serverside.
Similarly, `window.handleServerAction` is called if action is dispatched on serverside.

_Note: You can change these method names by `fsa.LorcaConfig`_

This is usage with redux-toolkit.

```js
// examples/cra/front/src/redux-lorca.js

export const makeLorcaMiddleware = () => (_store) => (next) => (action) => {
  dispatchToServer(action);
  next(action);
};

const makeServerActionHandler = (store) => (action) => store.dispatch(action);

export const setupServerActionHandler = (store) => {
  window.handleServerAction = makeServerActionHandler(store);
};
```

```js
// examples/cra/front/src/redux.js

import { configureStore } from "@reduxjs/toolkit";
import { makeLorcaMiddleware, setupServerActionHandler } from "./redux-lorca";

const makeStore = ({ reducer }) => {
  const store = configureStore({
    reducer,
    middleware: [makeLorcaMiddleware()],
  });
  setupServerActionHandler(store);
  return store;
};
```

See [examples/cra](https://github.com/mpppk/lorca-fsa/tree/master/examples/cra) for more details.
