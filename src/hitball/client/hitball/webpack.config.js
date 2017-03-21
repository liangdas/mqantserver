'use strict';
/**
 * Created by liangdas on 17/1/20.
 * Email 1587790525@qq.com
 */
var path = require('path');
var webpack = require('webpack');
var phaserModule = path.join(__dirname, '/node_modules/phaser/');
var phaser = path.join(phaserModule, 'build/custom/phaser-split.js'),
    pixi = path.join(phaserModule, 'build/custom/pixi.js'),
    p2 = path.join(phaserModule, 'build/custom/p2.js');
module.exports = {
    //devtool: 'eval-source-map',//配置生成Source Maps，选择合适的选项
    devtool:"cheap-module-source-map",
    entry:  __dirname + "/src/main.js",//已多次提及的唯一入口文件
    output: {
        path: __dirname + "/js",//打包后的文件存放的地方
        filename: "main.js"//打包后输出文件的文件名
    },

    module: {
        loaders: [
            {
                test: /\.json$/,
                loader: "json"
            },
            {
                test: /\.js$/,
                exclude: /node_modules/,
                loader: 'babel',//在webpack的module部分的loaders里进行配置即可
            },
            { test: /pixi.js/, loader: "script" }
        ]
    },

    devServer: {
        contentBase: "./js",
        colors: true,
        historyApiFallback: true,
        inline: true
    },
    resolve: {
        alias: {
            'phaser': phaser,
            'pixi.js': pixi,
            'p2': p2,
        }
    },
    plugins: [
        //这个使用uglifyJs压缩你的js代码
        new webpack.optimize.UglifyJsPlugin({
            mangle: {
                except: ['$super','_super','ctor', '$', 'exports', 'require']
                //以上变量‘$super’, ‘$’, ‘exports’ or ‘require’，不会被混淆
            },
            minimize: true,
            compress: {
                warnings: false
            }
        }),
    ]
}

