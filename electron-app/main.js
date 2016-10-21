'use strict';

const electron = require('electron');
// Module to control application life.
const app = electron.app;
// Module to create native browser window.
const BrowserWindow = electron.BrowserWindow;

const dorusu = require('dorusu');
const protobuf = dorusu.pb;

// Keep a global reference of the window object, if you don't, the window will
// be closed automatically when the JavaScript object is garbage collected.
let mainWindow;

function createWindow () {
  var PROTO_PATH = __dirname + '/../proto/';
  var chatpb = protobuf.requireProto(PROTO_PATH+'chat.proto', require);
  var ChatClient = chatpb.chatpb.MessageService.Client.raw;
  var client = new ChatClient({
    host: 'localhost',
    port: 50051,
    protocol: 'http:'
  });
  console.log(client);
  client.get({limit: 50, offset: 0}, function(response) {
    response.on('data', function(pb) {
      console.log('Conversations:', pb.message);
    });
    response.on('error', function(err) {
      console.error("Error calling GetConversations: ", err);
    });
  });

  // Create the browser window.
  mainWindow = new BrowserWindow({width: 800, height: 600, minWidth: 800, titleBarStyle: 'hidden' });

  // and load the index.html of the app.
  mainWindow.loadURL('file://' + __dirname + '/index.html');

  // Emitted when the window is closed.
  mainWindow.on('closed', function() {
    // Dereference the window object, usually you would store windows
    // in an array if your app supports multi windows, this is the time
    // when you should delete the corresponding element.
    mainWindow = null;
  });
}

// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
app.on('ready', createWindow);

// Quit when all windows are closed.
app.on('window-all-closed', function () {
  // On OS X it is common for applications and their menu bar
  // to stay active until the user quits explicitly with Cmd + Q
  if (process.platform !== 'darwin') {
    app.quit();
  }
});

app.on('activate', function () {
  // On OS X it's common to re-create a window in the app when the
  // dock icon is clicked and there are no other windows open.
  if (mainWindow === null) {
    createWindow();
  }
});
