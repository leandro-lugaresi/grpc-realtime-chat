import React, { Component, PropTypes } from 'react'
import { connect } from 'react-redux'
import socket from '../lib/socket'
import { adduser } from '../redux/actions/actions'
import { ADD_USER, REMOVE_USER } from '../constants/actionTypes'

class App extends Component {

  componentDidMount() {
    const { dispatch } = this.props
    socket.on('addUser', (user) => {
      dispatch(adduser({ type: ADD_USER, user }))
    })
    socket.on('removeUser', (user) => {
      dispatch({ type: 'REMOVE_USER', user })
    })
  }

  render() {
    return (
      <div className="appWrapper">
        { this.props.children }
      </div>
    )
  }
}

App.propTypes = {
  children: PropTypes.element.isRequired
}

function mapStateToProps(state, props) {
  return {
    users: state.users,
    messages: state.messages,
  }
}

export default connect(mapStateToProps)(App)
