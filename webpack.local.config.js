var webpack = require('webpack');
var path = require('path');
var config = require('./webpack.base.config.js');

var APP_DIR = path.resolve(__dirname, 'src/app');

config.entry = [
    APP_DIR + '/index.jsx'
];

config.output.publicPath = '/';

config.module.loaders = config.module.loaders.concat([
    {
        test : /\.jsx?/,
        include : APP_DIR,
        exclude: '/node_modules/',
        loaders : ['babel']
    }
]);

module.exports = config;