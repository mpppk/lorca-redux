/*global dispatchToServer*/

export const makeLorcaMiddleware = () => (_store) => (next) => (action) => {
    dispatchToServer(action)
    next(action)
}

const makeServerActionHandler = (store) => {
    return (action) => {
        store.dispatch(action);
    }
}

export const setupServerActionHandler = (store) => {
    window.handleServerAction = makeServerActionHandler(store);
}
