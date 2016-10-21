import React, { Component } from 'react'
import { connect } from 'react-redux'
import { Link } from 'react-router'
import { loginUser } from '../redux/actions/actions'

class Login extends Component {

  handleSubmit(e) {
    const { dispatch } = this.props
    const username = this.refs.username
    const password = this.refs.password
    const creds = { username: username.value.trim(), password: password.value.trim() }
    dispatch(loginUser(creds));
  }

  render() {
    return (
      <div>
        <div className="form">
          <img src='public/images/chatron-logo.png' className="chatron-logo" width='250px' />
          <form>
            <div className="form-group">
              <input type="email" className="form-control" placeholder="Email" ref="username" />
            </div>
            <div className="form-group">
              <input type="password" className="form-control" placeholder="Password" ref="password" />
            </div>
          </form>
          <a onClick={this.handleSubmit.bind(this)} className="form__button">Login</a>
          <Link to='/signup' className="form__link">Not a member? Sign up</Link>
        </div>
      </div>

    )
  }
}

function mapStateToProps(state, props) {
    return {
        users: state.users,
    }
}

export default connect(mapStateToProps)(Login)
