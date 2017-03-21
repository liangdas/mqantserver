'use strict';
/**
 * Created by liangdas on 2016/12/6 0006.
 * Email :1587790525@qq.com
 */

var Bullet = extend(GameRole, {
    ctor: function ctor(game, x, y, key, frame, group, properties) {
        properties = properties || {};
        this.roleType = "bullet";
        this.bulletType = null;
        this._super(game, x, y, key, frame);
        game.physics.arcade.enable(this);
        this.speed = properties.speed || 200;
        this.tileVolume = 1; //体积与瓦片地图判断是否可以通过
        //this.lifespan = 200; //能发射的长度
        this.checkWorldBounds = true;
        this.outOfBoundsKill = true;

        this.exists = false;
        this.visible = false;
        this.events.onOutOfBounds.add(this.resetBullet, this);

        if (group) {
            group.add(this);
        }
    },
    //  如果子弹飞出屏幕 就调用这个回调
    resetBullet: function resetBullet(bullet) {
        bullet.kill();
    },
    //重置 位置 角度 转动(0 3.12) 速度 重力
    rebirth: function rebirth(x, y, angle, rotation, gx, gy) {
        gx = gx || 0;
        gy = gy || 0;
        this.reset(x, y);
        this.scale.set(1);
        //this.lifespan = 200; //能发射的长度
        this.rotation = rotation; //设置子弹的角度
        this.game.physics.arcade.velocityFromAngle(angle, this.speed, this.body.velocity);
        this.angle = angle;
        this.body.gravity.set(gx, gy);
    },
    hit: function hit() {
        this.explode();
        this.dead();
    },
    //爆炸效果
    explode: function explode() {
        var boom = this.game.add.sprite(this.x, this.y, 'explode', 0);
        boom.anchor.setTo(0.5, 0.5);
        boom.width = 40; //设置对象比例
        boom.height = 40;
        var anim = boom.animations.add('boom', [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16], 20);
        anim.play();
        anim.onComplete.add(function () {
            boom.destroy();
        });
    }
});

var BulletFactory = function BulletFactory() {};
BulletFactory.prototype = {
    bulletTypes: ["bullet0", "bullet2", "bullet9", "bullet10"],
    createBullet: function createBullet(bulletType, game, group, properties) {
        bulletType = bulletType || this.randomBullet();
        var bullet = null;
        if (bulletType === "bullet0") {
            var b = new Bullet0(game, group, properties);
            b.name = 'bullet';
            bullet = b;
        } else if (bulletType === "bullet2") {
            var b = new Bullet2(game, group, properties);
            b.name = 'bullet';
            bullet = b;
        } else if (bulletType === "bullet9") {
            var b = new Bullet9(game, group, properties);
            b.name = 'bullet';
            bullet = b;
        } else if (bulletType === "bullet10") {
            var b = new Bullet10(game, group, properties);
            b.name = 'bullet';
            bullet = b;
        }
        return bullet;
    },
    randomBullet: function randomBullet() {
        var bulletIndex = Math.floor(Math.random() * (this.bulletTypes.length - 1));
        var bulletType = this.bulletTypes[bulletIndex];
        return bulletType;
    }
};

var Bullet0 = extend(Bullet, {
    ctor: function ctor(game, group, properties) {
        this._super.call(this, game, 0, 0, "bullet0", null, group, properties);
        this.bulletType = "bullet0";
    }
});

var Bullet2 = extend(Bullet, {
    ctor: function ctor(game, group, properties) {
        this._super.call(this, game, 0, 0, "bullet2", null, group, properties);
        this.bulletType = "bullet2";
    }
});
var Bullet9 = extend(Bullet, {
    ctor: function ctor(game, group, properties) {
        this._super.call(this, game, 0, 0, "bullet9", null, group, properties);
        this.bulletType = "bullet9";
    }
});
var Bullet10 = extend(Bullet, {
    ctor: function ctor(game, group, properties) {
        this._super.call(this, game, 0, 0, "bullet10", null, group, properties);
        this.bulletType = "bullet10";
    }
});