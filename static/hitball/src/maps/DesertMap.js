'use strict';

/**
 * Created by liangdas on 2016/12/6 0006.
 * Email :1587790525@qq.com
 */
var extend = require('../utils/inherits.js');
var GameMap = require('./GameMap.js');
module.exports = extend(GameMap, {
    ctor: function ctor(game, key, tileWidth, tileHeight, width, height) {
        this._super.call(this, game, "desert");
        this.addTilesetImage('Desert', 'tiles');

        //设置瓦片地图中哪些索引可以碰撞检测
        //只要瓦片设置了mesh属性的都是需要检测碰撞的
        //碰撞检测的规则可以参考README.md
        var setCollisions = [];
        for (var key in this.tilesets[0].tileProperties) {
            var tile = this.tilesets[0].tileProperties[key];
            if ("mesh" in tile) {
                setCollisions.push(parseInt(key) + 1); //墙的坐标是从 1开始的
            }
        }
        this.setCollision(setCollisions, true); //墙
        this.setTileIndexCallback(setCollisions, this.hitTile, this, "Ground");
    },
    //有物体撞击指定的tile了
    hitTile: function hitTile(sprite, tile) {
        if ("mesh" in tile.properties) {
            if (tile.properties.mesh - sprite.tileVolume > 0) {
                //可以通过
                return false;
            }
            //无法通过 判断碰撞的精灵类型
            return true; //返回 true 精灵无法穿过 返回 false 精灵可以穿过
        } else {
            //透明的可以直接穿过
            return false;
        }
    }
});