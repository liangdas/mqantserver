'use strict';
/**
 * Created by liangdas on 2016/12/6 0006.
 * Email :1587790525@qq.com
 */
var extend = require('../utils/inherits.js');
module.exports = extend(Phaser.Sprite, {
  //ctor 可以省略  省略以后会继续执行其父构造函数 如 this._super.apply(this,arguments);
  ctor: function ctor(game, x, y, key, frame, properties) {
    properties = properties || {};
    this._super(game, x, y, key, frame);
    this.rid = properties.rid; //系统角色
    this.wid = properties.wid; //战场中的角色ID
  },
  dead: function dead() {
    if(this.alive){
      this.kill();
    }
  }
});