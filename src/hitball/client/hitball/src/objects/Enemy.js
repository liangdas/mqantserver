"use strict";

/**
 * Created by liangdas on 2016/12/6 0006.
 * Email :1587790525@qq.com
 */
var extend = require('../utils/inherits.js');
var Player = require('./Player.js');
module.exports =  extend(Player, {
    ctor: function ctor(game, x, y, group, properties) {
        this.roleType = "deadly";
        properties = properties || {};
        this._super(game, x, y, "deadly", null, properties);
    }
});