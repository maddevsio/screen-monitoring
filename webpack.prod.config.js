var webpack = require('webpack');
var path = require('path');
var config = require('./webpack.base.config.js');

var APP_DIR = path.resolve(__dirname, 'src/app');


config.entry = APP_DIR + '/index.jsx';
config.output.publicPath = '/public/';
config.module.loaders = config.module.loaders.concat([
    {
        test : /\.jsx?/,
        include : APP_DIR,
        exclude: '/node_modules/',
        loaders : ['babel']
    }
]);
config.plugins = config.plugins.concat([
    new webpack.DefinePlugin({
        'process.env': {
            'NODE_ENV': JSON.stringify('production')
        }
    })
]);

module.exports = config;