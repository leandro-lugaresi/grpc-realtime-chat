import React, { Component, PropTypes } from 'react'
import { connect } from 'react-redux'
import ChatroomInput from './ChatroomInput'
import _ from 'lodash'
import Message from './Message'
import { logoutUser } from '../redux/actions/actions'

class ChatroomWindow extends Component {

  render() {
    const { messages, users, params, auth, dispatch } = this.props

    return (
      <div className="form-group chatroom-window">
        <div className="form-group chatroom-content">
          { (messages && messages.length > 0) ? messages.map((message, i) => {
            return (
              <Message loggedInUser={ auth.user.userId } fromUser={ message.from } content={ message.content } key={ i }/>
            )
          }) : null }
        </div>
        <ChatroomInput loggedInUser={ auth } toUser={ params.userId }/>
      </div>
    )
  }
}

function mapStateToProps(state, props) {
  return {
    users: state.users,
    auth: state.auth,
    messages: state.messages[props.params.userId]
  }
}

export default connect(mapStateToProps)(ChatroomWindow)
