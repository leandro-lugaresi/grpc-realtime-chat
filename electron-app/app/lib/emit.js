import socket from './socket';

module.exports.sendMessage = (message) => {
  socket.emit('directMessage', message)
}
