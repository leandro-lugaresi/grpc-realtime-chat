import _ from 'lodash'
import { RECEIVE_MESSAGE } from '../../constants/actionTypes'

export function messages(state = {}, action) {
  switch (action.type) {

    case RECEIVE_MESSAGE:
      const { userId } = action.auth.user
      const { content, toUser, fromUser } = action.msg
      const toUserArray = state[toUser] ? state[toUser] : []
      const fromUserArray = state[fromUser] ? state[fromUser] : []

      if ( toUser !== userId && fromUser !== userId) return state
      return Object.assign({}, state, { [toUser]: toUserArray.concat({ from: fromUser, content }) }, { [fromUser]: fromUserArray.concat({ from: fromUser, content }) })

    default:
      return state

  }
}
