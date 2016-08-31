var webpack = require('webpack');
var path = require('path');
var HtmlWebpackPlugin = require('html-webpack-plugin');

var BUILD_DIR = path.resolve(__dirname, 'public');
var APP_DIR = path.resolve(__dirname, 'src/app');


var config = {
    entry: APP_DIR + '/index.jsx',
    output: {
        path: BUILD_DIR,
        filename: 'bundle-[hash].js',
        publicPath: '/'
    },
    module : {
        loaders : []
    },
    plugins: [
        new HtmlWebpackPlugin({
            title: 'Taxi operator',
            hash: true,
            inject: 'body',
            template: './template/index.html'
        })
    ]
};

module.exports = config;