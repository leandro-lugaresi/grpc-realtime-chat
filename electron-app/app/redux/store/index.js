import { createStore, applyMiddleware, compose } from 'redux'
import createLogger from 'redux-logger'
import thunk from 'redux-thunk'
import { hashHistory } from 'react-router'
import { syncHistory } from 'react-router-redux'
import rootReducer from '../reducers'

const reduxRouterMiddleware = syncHistory(hashHistory)

const logger = createLogger({
  level: 'info',
  collapsed: true
});

const middleware = applyMiddleware(thunk, reduxRouterMiddleware, logger)

export default function configureStore(initialState) {
  const store = createStore(rootReducer, initialState, middleware)

  if (module.hot) {
    module.hot.accept('../reducers', () =>
      store.replaceReducer(require('../reducers'))
    )
  }

  return store
}
