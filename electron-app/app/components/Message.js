import React, { Component } from 'react'

export default class Message extends Component {
  render () {
    const { content, fromUser, loggedInUser } = this.props
    if (fromUser === loggedInUser) return <div className="bubble you">{ content }</div>
    return <div className="bubble me">{ content }</div>
  }
}
