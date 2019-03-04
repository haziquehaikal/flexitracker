const merge = require('webpack-merge');
const common = require('./webpack.common.js');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const webpack = require('webpack')


module.exports = merge(common , {
    devtool: 'inline-source-map',
    plugins: [
        new HtmlWebpackPlugin({
            template: `${__dirname}/public/index.html`,
            title:'Production'
        }),
        new webpack.DefinePlugin({
            'process.env.NODE_ENV': JSON.stringify('development'),
        })
    ],
    devServer: {
        contentBase: './dist',
        historyApiFallback: true,
        disableHostCheck: true,
        watchOptions: {
          poll: 1500,
        },
        overlay: {
          warnings: true,
          errors: true,
        },
        open: true,
        port: 3000,
      },
});