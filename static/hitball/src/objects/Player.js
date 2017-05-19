"use strict";

/**
 * Created by liangdas on 2016/12/6 0006.
 * Email :1587790525@qq.com
 */
var extend = require('../utils/inherits.js');
var GameRole = require('./Role.js');
var mqant=window.mqant
module.exports = extend(GameRole, {
    ctor: function ctor(game, x, y, group, properties) {
        this.roleType = "hero";
        properties = properties || {};
        this.game = game;
        this.arrows = properties.arrowsGroup; //视觉观察组
        this.rotateDirection = 1; // rotate direction: 1-clockwise, 2-counterclockwise
        this.rotateSpeed = 3; // arrow rotation speed
        this.friction = 0.99; // friction affects ball speed 速度递减因子
        this.arrow; // rotating arrow
        this.minPower = 50; // minimum power applied to ball
        this.maxPower = properties.maxPower || 200; // maximum power applied to ball
        this.power=this.minPower; //力量
        this.ballRadius=10; //周长
        this.degToRad=0.0174532925; // degrees-radians conversion
        this._super(this.game, x, y, "ball", null, properties);
        this.xSpeed = 0;
        this.ySpeed = 0;
        this.game.physics.arcade.enable(this);
        this.anchor.setTo(0.5, 0.5);
        this.checkWorldBounds = true;
        this.outOfBoundsKill = true;
        this.body.collideWorldBounds = true; //与世界边境进行物理检测
        this.inputEnabled = true;
        this.input.useHandCursor = true; //当鼠标移动到其上面时显示小手
        this.powering=false;
        //this.input.enableDrag(); //可以拖动
        if (group) {
            group.add(this);
        }
    },
    dead: function () {
        if(this.alive){
            this._super();
            this.arrow.kill();
        }
    },
    getArrow: function getArrow() {
        if (this.arrow == null) {
            var arrow = this.arrows.getFirstExists(false);
            if (arrow) {
                this.arrow = arrow;
                this.arrow.reset(this.x, this.y);
            } else {
                //设置一个观察器
                this.arrow = this.game.add.sprite(this.game.world.centerX,this.game.world.centerY,"arrow");
                this.arrow.anchor.x = -1;
                this.arrow.anchor.y = 0.5;
            }
        }
        return this.arrow;
    },
    Power:function (){
        this.power++;
        this.power = Math.min(this.power,this.maxPower)
        this.powering=true;
    },
    Fire: function () {
        //发射
        //this.xSpeed += Math.cos(this.getArrow().angle*this.degToRad)*this.power/20;
        //this.ySpeed += Math.sin(this.getArrow().angle*this.degToRad)*this.power/20;
        this.power = this.minPower;
        mqant.requestNR("Hitball/HD_Fire",{
            "Rid": "001",
            "Angle": this.getArrow().angle,
            "Power": this.power,
            "X": this.x,
            "Y": this.y,
        });
        this.rotateDirection*=-1;
        this.powering=false;
    },
    Rotate: function () {
        this.getArrow().angle+=this.rotateSpeed*this.rotateDirection;
    },
    OnMove:function(player){
        this.x=player.X;
        this.y=player.Y;
        this.xSpeed=player.XSpeed;
        this.ySpeed=player.YSpeed;
        this.power=player.Power;
        this.rotateDirection=player.RotateDirection
        this.getArrow().x=this.x;
        this.getArrow().y=this.y;
        if(!this.powering){
            this.angle=player.Angle
            this.getArrow().angle=player.Angle
        }
    },
    OnRotate:function(player){
        this.angle=player.Angle
        this.getArrow().angle=player.Angle
    },
    Move: function () {
        this.x=this.x+this.xSpeed;
        this.y=this.y+this.ySpeed;
        // reduce ball speed using friction 速度递减
        this.xSpeed*=this.friction;
        this.ySpeed*=this.friction;
        // update arrow position 更新选择箭头
        this.getArrow().x=this.x;
        this.getArrow().y=this.y;
        //var self=this;
        //向服务器汇报
        //mqant.requestNR("Hitball/HD_Move",{
        //    "wid": "1",
        //    "war": "001",
        //    "x":x,
        //    "y":y
        //});
    }
});