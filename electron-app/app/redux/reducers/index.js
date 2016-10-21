import { combineReducers } from 'redux';
import { messages } from './messages';
import { users } from './users';
import { auth } from './auth';

export default combineReducers({
  messages,
  auth,
  users
});
