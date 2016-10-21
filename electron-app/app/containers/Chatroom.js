import React, { Component, PropTypes } from 'react'
import _ from 'lodash'
import { connect } from 'react-redux'
import FriendsList from '../components/FriendsList'
import socket from '../lib/socket'
import { receiveMessage } from '../redux/actions/actions'
import { RECEIVE_MESSAGE } from '../constants/actionTypes'

class Chatroom extends Component {

  componentDidMount() {
    const { auth, dispatch } = this.props
    socket.on('directMessage', (msg) => {
      dispatch({ type: RECEIVE_MESSAGE, msg, auth })
    })
  }

  handleLogout (e) {
    const { dispatch, auth } = this.props
    dispatch(logoutUser(auth.user));
  }

  render () {
    const { users, auth, params } = this.props
    const talkingTo = _.find(users, { 'userId': params.userId })
    return (
      <div className="window">
        <div className="window-content">
          <div className="pane-group">
            <div className="pane-sm sidebar">
              <img src='public/images/chatron-logo.png' className="chatron-logo" width='100px' />
              <FriendsList loggedInUser={ auth } users={ users } toUser={ params }/>
            </div>
            <div className="pane">
              <div className="chatroom-navbar">
                { talkingTo ? <span><strong>To:</strong> { talkingTo.username } </span> : null }
                <a className="pull-right" onClick={ this.handleLogout.bind(this) }>Logout</a>
              </div>
              { this.props.children }
            </div>
          </div>
        </div>
      </div>
    )
  }
}

function mapStateToProps(state, props) {
  return {
    users: state.users,
    auth: state.auth
  }
}

export default connect(mapStateToProps)(Chatroom)
