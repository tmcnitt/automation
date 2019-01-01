var path = require('path');
var webpack = require('webpack');

module.exports = {
  entry: './src/index.js',
  output: {
    path: __dirname,
    filename: './build/bundle.js'
  },

  devtool: 'source-map',
  module: {
    loaders: [{
      test: /.js?$/,
      loader: 'babel-loader',
      include: [
        path.resolve(__dirname, "src"),
      ],
      query: {
        presets: ['es2015', 'stage-0', 'react'],
        plugins: ['transform-class-properties', "transform-object-assign", 'babel-plugin-transform-decorators-legacy'],
      }
    }, {
      test: /\.css$/,
      loader: "style-loader!css-loader"
    }]
  },
  plugins: [
          new webpack.ProvidePlugin({
              'Promise': 'es6-promise', // Thanks Aaron (https://gist.github.com/Couto/b29676dd1ab8714a818f#gistcomment-1584602)
              'fetch': 'imports?this=>global!exports?global.fetch!whatwg-fetch'
          }),
          new webpack.ResolverPlugin(
           new webpack.ResolverPlugin.DirectoryDescriptionFilePlugin(".bower.json", ["main"])
          ),
          new webpack.IgnorePlugin(/vertx/)
   ],
  resolve: {
    root: path.resolve(__dirname, 'src'),
    modulesDirectories: ["web_modules", "node_modules", "bower_components"]
  },
};
