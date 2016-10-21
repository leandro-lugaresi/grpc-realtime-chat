var webpack = require('webpack')
const path = require('path')

module.exports = {

  context: path.resolve(__dirname, 'app'),

  entry: './app.js',

  output: {
    filename: 'bundle.js',
    path: __dirname + '/build',
    publicPath: 'http://localhost:8080/build/'
  },

  module: {
    loaders: [
      {
        test: /\.js$/,
        exclude: /node_modules/,
        loaders: ['react-hot', 'babel?presets[]=react,presets[]=es2015'],
      },
      {
        test: /\.html$/,
        loader: 'file?name=[name].[ext]',
      },
      {
        test   : /\.woff|\.woff2|\.svg|.eot|\.ttf/,
        loader: 'file-loader'
        // loader: require.resolve("file-loader") + "?name=../[path][name].[ext]"
      },
      {
        test: /\.less$/,
        loader: 'style!css!less'
      },
      {
        test: /.*\.(gif|png|svg)$/i,
        loaders: [
            'file?hash=sha512&digest=hex&name=[hash].[ext]',
            'image-webpack?{progressive:true, optimizationLevel: 7, interlaced: false, pngquant:{quality: "65-90", speed: 4}}'
        ]
      },
      {
        test: /\.(jpg)$/,
        loader: 'url?limit=25000'
      }
    ]
  }
}
