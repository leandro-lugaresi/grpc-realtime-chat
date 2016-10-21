import React, { Component } from 'react'
import { Link } from 'react-router'
import { connect } from 'react-redux'
import { signup } from '../redux/actions/actions.js'


class Signup extends Component {

  handleClick() {
    const { dispatch } = this.props
    dispatch(signup({ username: this.refs.username.value, password: this.refs.password.value}))
  }

  renderErrorMsg() {
    const { auth } = this.props
    if (!auth.signUpError) return null
    return (
      <p className="danger"><strong>Error:</strong> { auth.signUpErrorMsg }</p>
    )
  }

  render() {
    const { auth } = this.props
    return (
      <div>
        <div className="form">
          <img src='public/images/chatron-logo.png' className="chatron-logo" width='250px' />
          { !auth.signUpSuccess ? this.renderErrorMsg() : null }
          <form>
            <div className="form-group">
              <input type="email" className="form-control" placeholder="Email" ref="username" />
            </div>
            <div className="form-group">
              <input type="password" className="form-control" placeholder="Password" ref="password" />
            </div>
          </form>
          <a className="form__button" onClick={ this.handleClick.bind(this) }>Signup</a>
          <Link to='/' className="form__link">Back to Login</Link>
        </div>
      </div>
    )
  }
}

function mapStateToProps(state, props) {
  return {
    auth: state.auth
  }
}

export default connect(mapStateToProps)(Signup)
