"use strict";

/**
 * Created by liangdas on 2016/12/6 0006.
 * Email :1587790525@qq.com
 */
var extend = require('../utils/inherits.js');
var GameRole = require('./Role.js');
module.exports =  extend(GameRole, {
    ctor: function ctor(game, x, y, group, properties) {
        this.roleType = "enemy";
        properties = properties || {};
        this._super(game, x, y, "coin", null, properties);
        this.game = game;
        this.speed = properties.speed || 50;

        this.nextFire = 0; //下一次发射子弹的时间
        this.fireRate = 50; //发射速率 50ms
        this.tileVolume = 3; //体积与瓦片地图判断是否可以通过
        this.game.physics.arcade.enable(this);
        this.anchor.x = 0.5;
        this.anchor.y = 0.5;
        this.anchor.setTo(0.5, 0.5);
        this.checkWorldBounds = true;
        this.outOfBoundsKill = true;
        this.body.collideWorldBounds = true; //与世界边境进行物理检测
        this.inputEnabled = true;
        this.input.useHandCursor = true; //当鼠标移动到其上面时显示小手
        //this.input.enableDrag(); //可以拖动
        if (group) {
            group.add(this);
        }
    },
    move: function move() {

    },
    rebirth: function rebirth(x, y) {
        this.reset(x, y);
    }
});