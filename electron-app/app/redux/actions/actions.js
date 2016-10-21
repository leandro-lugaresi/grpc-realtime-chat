import { push } from 'react-router-redux'
import {
        SIGN_UP_REQUEST,
        SIGN_UP_SUCCESS,
        SIGN_UP_FAILURE,
        LOGIN_REQUEST,
        LOGIN_SUCCESS,
        LOGIN_FAILURE,
        FETCH_AVAILABLE_USERS_REQUEST,
        FETCH_AVAILABLE_USERS_SUCCESS,
        FETCH_AVAILABLE_USERS_FAILURE,
        ADD_USER,
        RECEIVE_MESSAGE,
        LOGOUT_REQUEST,
        LOGOUT_SUCCESS
      } from '../../constants/actionTypes'

const BASE_URL = 'http://localhost:3001'

/**
* Action creator which posts the users signup creds and dispatches either
* signUpFailure or signUpSuccess
* @param { user } user's credentials incl. username && password
* @returns { dispatch }
*/

export function signup(user) {

  let config = {
    method: 'POST',
    headers: { 'Accept': 'application/json', 'Content-Type': 'application/json' },
    body: JSON.stringify(user)
  }

  return dispatch => {
    dispatch(signUpRequest())
    return fetch(`${BASE_URL}/api/signup`, config)
      .then(res => res.json())
      .then(res => {
        if (res.err) return dispatch(signUpFailure({ err: res.err }))
        dispatch(signUpSuccess())
        dispatch(push('/'))
    }).catch(err => { throw err })
  }

}

const signUpRequest = (user) => {
  return { type: SIGN_UP_REQUEST }
}

const signUpSuccess = () => {
  return { type: SIGN_UP_SUCCESS }
}

const signUpFailure = (msg) => {
  return { type: SIGN_UP_FAILURE, msg }
}

/**
* Action creator which posts the users login creds and dispatches either
* loginSuccess or loginError
* @param { creds } user's credentials incl. username && password
* @returns { dispatch }
*/

export function loginUser(creds) {

  let config = {
    method: 'POST',
    headers: { 'Content-Type':'application/x-www-form-urlencoded' },
    body: `username=${creds.username}&password=${creds.password}`
  }

  return dispatch => {
    dispatch(loginRequest(creds))

    return fetch(`${BASE_URL}/api/sessions/create`, config )
      .then(response => response.json().then(user => ({ user, response })))
      .then(({ user, response }) => {
        if (!response.ok) {
          dispatch(loginError(user.message))
          return Promise.reject(user)
        } else {
          localStorage.setItem('id_token', user.id_token)
          dispatch(loginSuccess(user))
          dispatch(fetchAvailableUsers())
          dispatch(push('/chat'))
        }
    }).catch(err => console.log("Error: ", err))
  }
}

const loginRequest = (creds) => {
  return {
    type: LOGIN_REQUEST,
    isFetching: true,
    isAuthenticated: false,
    creds
  }
}

const loginSuccess = (user) => {
  return {
    type: LOGIN_SUCCESS,
    isFetching: false,
    isAuthenticated: true,
    id_token: user.id_token,
    user: user.user
  }
}

const loginError = (message) => {
  return {
    type: LOGIN_FAILURE,
    isFetching: false,
    isAuthenticated: false,
    message
  }
}


/**
* Action creator which is responsible for logging the user out. It removes the local storage token and dispatches LOGOUT_SUCCESS action
* @returns { dispatch }
*/

export function logoutUser(creds) {

  let config = {
    method: 'POST',
    headers: { 'Accept': 'application/json', 'Content-Type': 'application/json' },
    body: JSON.stringify(creds)
  }

  return dispatch => {
    dispatch(logoutRequest())
    return fetch(`${BASE_URL}/api/sessions/close`, config)
      .then(response => response.json().then(user => ({ user, response })))
      .then(({ user, response }) => {
        console.log('About to logout...')
        dispatch(logoutSuccess())
        localStorage.removeItem('id_token')
        dispatch(push('/'))
      }).catch(err => { throw err })
  }
}

const logoutRequest = () => {
  return {
    type: LOGOUT_REQUEST,
    isFetching: true,
    isAuthenticated: true
  }
}

const logoutSuccess = () => {
  return {
    type: LOGOUT_SUCCESS,
    isFetching: false,
    isAuthenticated: false
  }
}

/**
* Action creator which fetches available users and dispatches either
* fetchAvailableUsersSuccess or fetchAvailableUsersFailure
* @returns { dispatch }
*/

export function fetchAvailableUsers() {

  let config = {
    method: 'GET',
    headers: { 'Accept': 'application/json', 'Content-Type': 'application/json' },
  }

  return dispatch => {
    dispatch(fetchAvailableUsersRequest())
    return fetch(`${BASE_URL}/api/users/get`, config)
      .then(response => response.json())
      .then(response => {
        if (response.err) return dispatch(fetchAvailableUsersFailure())
        dispatch(fetchAvailableUsersSuccess(response))
      }).catch(err => { throw err })
  }

}

const fetchAvailableUsersRequest = () => {
  return { type: FETCH_AVAILABLE_USERS_REQUEST }
}

const fetchAvailableUsersSuccess = (users) => {
  return { type: FETCH_AVAILABLE_USERS_SUCCESS, users }
}

const fetchAvailableUsersFailure = () => {
  return { type: FETCH_AVAILABLE_USERS_FAILURE }
}

export function adduser(user) {
  return { type: ADD_USER, user }
}
