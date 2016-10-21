import React, { Component, PropTypes } from 'react'
import { Link } from 'react-router'
import _ from 'lodash'

export default class FriendsList extends Component {

  render() {
    const { users, loggedInUser, toUser } = this.props
    return (
      <ul className="list-group">
        <li className="list-group-header" key="header">
          <input className="form-control" type="text" placeholder="Search for someone" />
        </li>
        {
          _.map(users, (user, userId) => {
            if (loggedInUser.user.userId === userId || !user.currentlyOnline || !user) return null
            return (
              <Link to={ `/chat/${userId}` } key={ user.username } >
                <li className={`list-group-item ${userId === toUser.userId ? 'active' : '' } `} key={ user.username }>
                  <img className="img-circle media-object pull-left" src="public/images/heisenburg.jpeg" width="32" height="32" />
                  <div className="media-body">
                    <strong>{ user.username }</strong>
                  </div>
                </li>
              </Link>
            )
          })
        }
      </ul>
    )
  }
}
