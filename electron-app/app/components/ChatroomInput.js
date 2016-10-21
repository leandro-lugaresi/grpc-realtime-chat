import React, { Component } from 'react'
import { sendMessage } from '../lib/emit'

export default class ChatroomInput extends Component {

  handleInputChange(e) {
    const { toUser, loggedInUser } = this.props
    if (e.keyCode === 13) {
      sendMessage({ content: this.refs.msg.value, toUser: toUser, fromUser: loggedInUser.user.userId })
      this.refs.msg.value = ''
    }
  }

  render () {
    const { toUser, loggedInUser } = this.props
    return (
      <div className="input">
        <input type="text" autoFocus="true" placeholder="Type here to chat!" ref="msg"  className="form-control input" onKeyDown={ this.handleInputChange.bind(this) } />
        <a className="input-link send" onClick={() => sendMessage({ content: this.refs.msg.value, toUser: toUser, fromUser: loggedInUser.user.userId })}><span className="icon icon-paper-plane"></span></a>
      </div>
    )
  }
}
