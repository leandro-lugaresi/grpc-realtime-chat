import "./less/site.less";
import React from 'react'
import ReactDOM from 'react-dom'
import { Router, Route, IndexRoute, hashHistory } from 'react-router'
import { Provider } from 'react-redux'
import App from './containers/App'
import Login from './containers/Login'
import Signup from './containers/Signup'
import Chatroom from './containers/Chatroom'
import ChatroomWindow from './components/ChatroomWindow'
import configureStore from './redux/store'

const store = configureStore()

ReactDOM.render((
  <Provider store={ store }>
    <Router history={ hashHistory }>
      <Route path='/' component={ App } >
        <IndexRoute component={ Login } />
        <Route path='/signup' component={ Signup } />
        <Route path='/chat' component={ Chatroom } >
          <Route path='/chat/:userId' component={ ChatroomWindow } />
        </Route>
      </Route>
    </Router>
  </Provider>
), document.getElementById('app'))
