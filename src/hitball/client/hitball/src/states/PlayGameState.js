"use strict";
/**
 * Created by liangdas on 2016/12/6 0006.
 * Email :1587790525@qq.com
 * 主界面
 */
var extend = require('../utils/inherits.js');
var Hero = require('../objects/Player.js');
var Enemy = require('../objects/Enemy.js');
var Coin = require('../objects/Coin.js');
var GameMap = require('../maps/DesertMap.js');
var mqant=window.mqant
var hudText; // text to display game info
var charging=false; // are we charging the power?
var score = 0; // the score
var gameOver = false; // flag to know if the game is over
module.exports = extend(function () {}, {
    ctor: function (Rid,player,coins) {
        this.heroRid=Rid;
        this.initplayer=player;
        this.initcoins=coins;
    },
    init: function init() {
        //this.game.renderer.renderSession.roundPixels = true;
        var self=this;
        this.physics.startSystem(Phaser.Physics.ARCADE);
        this.game.stage.backgroundColor = '#454645';
        this.game.time.desiredFps = 40;
    },
    preload: function preload() {
        this.game.load.image("ball", "assets/ball.png");
        this.game.load.image("deadly", "assets/deadly.png");
        this.game.load.image("coin", "assets/coin.png");
        this.game.load.image("arrow", "assets/arrow.png");
        this.game.load.atlas('generic', 'assets/virtualjoystick/skins/generic-joystick.png', 'assets/virtualjoystick/skins/generic-joystick.json');
    },
    joinHero: function(role) {
        var hero = new Hero(this.game, role.body.x, role.body.y, this.heros, {
            arrowsGroup: this.arrows,
            rid: role.rid,
            wid: role.wid,
        });
        return  hero
    },
    joinDeadly: function(role) {
        var x = role.body.x;
        var y = role.body.y;
        var enemy = new Enemy(this.game, x, y, this.deadlys, {
            arrowsGroup: this.arrows,
            rid: role.rid,
            wid: role.wid
        });
        return  enemy
    },
    joinCoin: function(role) {
        var coin = new Coin(this.game, role.body.x, role.body.y, this.coins, {
            rid: role.rid,
            wid: role.wid,
            fireRate: role.fireRate
        });
        return  coin
    },
    create: function create() {
        // center and scale the stage
        var self=this;
        //this.game.scale.pageAlignHorizontally = true;
        //this.game.scale.pageAlignVertically = true;
        this.game.world.resize(1280,1280);
        this.game.world.setBounds(0, 0, 1280,1280);
        //this.game.stage.width = this.game.width;
        //this.game.stage.height = this.game.height;
        //this.map = new GameMap(this.game);
        //this.layer = this.map.createLayer('Ground');
        //this.layer.resizeWorld(); //调整世界的范围跟地图范围相同

        console.log(this.game.world.width + " game " + this.game.world.height);
        console.log(this.game.world.width + " game " + this.game.world.height);

        this.arrows = this.game.add.group(); //视觉观察组
        this.heros = this.game.add.group();
        this.deadlys = this.game.add.group();
        this.coins = this.game.add.group();

        //初始化所有角色
        for(var k in this.initplayer){
            var player=this.initplayer[k];
            var roleBall={
                body:{
                    x:player.X,
                    y:player.Y
                },
                rid:player.Rid,
                wid:player.Wid,
                maxPower:300
            }
            var p = this.joinHero(roleBall);
            if(this.heroRid===player.Rid){
                this.hero=p;
            }
        }
        //初始化所有金币
        for(var ck in this.initcoins){
            var coin=this.initcoins[ck];
            var roleEnemy={
                body:{
                    x:coin.X,
                    y:coin.Y
                },
                rid:coin.Id,
                wid:coin.Wid,
                maxPower:300
            }
            self.joinCoin(roleEnemy)
        }


        this.game.camera.deadzone = new Phaser.Rectangle(200, 200, this.game.stage.width - 200, this.game.stage.height - 200); //镜头跟随触发区域dead zone
        this.game.camera.follow(this.hero); //摄像机跟随人物


        // create and place the text showing speed and score
        hudText = this.game.add.text(5,0,"",{
            font: "11px Arial",
            fill: "#ffffff",
            align: "left"
        });

        // update text content
        this.updateHud();

        // listener for input down
        this.game.input.onDown.add(this.charge, this);

        //this.game.input.keyboard.addKeyCapture([Phaser.Keyboard.SPACEBAR]);
        //this.pad = this.game.plugins.add(Phaser.Plugin.VirtualJoystick);
        //this.pad.start();
        //this.buttonA = this.pad.addButtonByKey("buttonA", this.charge, this);
        //
        //this.buttonB = this.pad.addButtonByKey("buttonB", this.charge, this);
        //this.buttonC = this.pad.addButtonByKey("buttonC", this.charge, this);
        mqant.on('Hitball/OnMove', function(data) {
            var message=JSON.parse(data.payloadString);
            self.heros.forEachAlive(function(player){
                var role=message[player.rid];
                if(role!=null){
                    //有数据
                    player.OnMove(role);
                }else{
                    player.dead();
                }
            },self);
            //添加本地没有的玩家
            for(var Rid in message){
                var role=message[Rid];
                var player=self.heros.iterate("rid",Rid,Phaser.Group.RETURN_CHILD);
                if(player===null){
                    var roleEnemy={
                        body:{
                            x:role.X,
                            y:role.Y
                        },
                        rid:role.Rid,
                        wid:role.Wid,
                        maxPower:300
                    }
                    self.joinHero(roleEnemy)
                }else{
                    //存在就不处理了
                }
            }
        });

        mqant.on('Hitball/OnJoin', function(data) {
            var role=JSON.parse(data.payloadString);
            var player=self.heros.iterate("rid",role.Rid,Phaser.Group.RETURN_CHILD);
            if(player===null){
                var roleEnemy={
                    body:{
                        x:role.X,
                        y:role.Y
                    },
                    rid:role.Rid,
                    wid:role.Wid,
                    maxPower:300
                }
                self.joinHero(roleEnemy)
            }else{
                console.log("该角色已存在!")
            }
        });

        mqant.on('Hitball/OnAddCoins', function(data) {
            console.log("添加金币");
            var role=JSON.parse(data.payloadString);
            var coin=self.coins.iterate("rid",role.Id,Phaser.Group.RETURN_CHILD);
            if(coin===null){
                var roleEnemy={
                    body:{
                        x:role.X,
                        y:role.Y
                    },
                    rid:role.Id,
                    wid:role.Wid,
                    maxPower:300
                }
                self.joinCoin(roleEnemy)
            }else{
                console.log("该金币已存在!")
            }
        });
        mqant.on('Hitball/OnEatCoins', function(data) {
            console.log("吃掉金币金币");
            var role=JSON.parse(data.payloadString);
            var coin=self.coins.iterate("rid",role.Id,Phaser.Group.RETURN_CHILD);
            if(coin===null){

            }else{
                coin.dead();
            }
        });
    },

    //// the function to place a coin is similar to the one which places the enemy, but this time we don't need
    //// to place it in an array because there's only one coin on the stage
    //placeCoin:function (){
    //    var randomX=Math.random()*(this.game.width-2*this.hero.ballRadius)+this.hero.ballRadius;
    //    var randomY=Math.random()*(this.game.height-2*this.hero.ballRadius)+this.hero.ballRadius;
    //    var roleCoin={
    //        body:{
    //            x:randomX,
    //            y:randomY
    //        },
    //        rid:"003",
    //        wid:"003",
    //        maxPower:300
    //    }
    //    var coin = this.joinCoin(roleCoin)
    //},
    update: function update() {
        // the game is update only if it's not game over
        var self=this;
        if(!gameOver){
            this.game.physics.arcade.collide(this.heros, this.deadlys, this.collisionHitDeadly, null, this); //子弹与敌人
            this.game.physics.arcade.collide(this.heros, this.coins, this.collisionHitCoin, null, this); //子弹与敌人
            // when the player is charging the power, this is increased until it reaches the maximum allowed
            if(charging){
                this.hero.Power(); //积蓄力量
                // then game text is updated
                this.updateHud();
            }

            // if the player is not charging, keep rotating the arrow
            else{
                self.heros.forEachAlive(function(player){
                    player.Rotate() //旋转
                },self);
            }
            self.heros.forEachAlive(function(player){
                player.Move() //移动
            },self);

            // handle wall bounce
            //this.wallBounce();


            //虚拟键盘
            //if (this.pad.isDragging) {
            //    this.hero.move(this.pad.angle, this.pad.force);
            //} else {
            //    this.hero.move();
            //}
        }
    },
    collisionHitDeadly: function collisionHitHandler(bullet, player) {
        //游戏角色被子弹击中
        //gameOver = true;
        //window.alert("Game Over！分数：" + score + "分，去发ajax吧！");
    },
    collisionHitCoin: function collisionHitHandler(player,coin) {
        //捡到金币了
        mqant.requestNR("Hitball/HD_EatCoin",{
            "Id": coin.rid
        });
        score += 1;
        coin.dead();
        //this.placeDeadly();
        //this.placeCoin();
        this.updateHud();
    },


    // function to handle bounces. Just check for game boundary collision
    //边缘碰撞检测
    wallBounce:function (){
        if(this.hero.x<this.hero.ballRadius){
            this.hero.x=this.hero.ballRadius;
            this.hero.xSpeed*=-1
        }
        if(this.hero.y<this.hero.ballRadius){
            this.hero.y=this.hero.ballRadius;
            this.hero.ySpeed*=-1
        }
        if(this.hero.x>this.game.world.width-this.hero.ballRadius){
            this.hero.x=this.game.world.width-this.hero.ballRadius;
            this.hero.xSpeed*=-1
        }
        if(this.hero.y>this.game.world.height-this.hero.ballRadius){
            this.hero.y=this.game.world.height-this.hero.ballRadius;
            this.hero.ySpeed*=-1
        }
    },

    // simple function to get the distance between two sprites
    // does not use sqrt to save CPU
    //碰撞检测
    getDistance:function (from,to){
        var xDist = from.x-to.x
        var yDist = from.y-to.y;
        return xDist*xDist+yDist*yDist;
    },

    // when the player is charging, set the power to min power allowed
    // and wait the player to release the input to fire the ball
    charge:function (){
        this.game.input.onDown.remove(this.charge, this);
        this.game.input.onUp.add(this.fire, this);
        charging=true;
    },

    // FIRE!!
    // update ball speed according to arrow direction
    // invert arrow rotation
    // reset power and update game text
    // wait for the player to fire again
    fire:function (){
        this.game.input.onUp.remove(this.fire, this);
        this.game.input.onDown.add(this.charge, this);
        this.hero.Fire()
        this.updateHud();
        charging=false;

    },

    // function to update game text
    updateHud:function (){
        hudText.text = "Power: "+this.hero.power+" * Score: "+score
    },
    paused:function(){
        console.log(" paused ");
    },
    resumed:function(){
        console.log(" resumed ");
        //this.game.world.resize(1800,1600);
        ////this.game.stage.width = this.game.world.width;
        ////this.game.stage.height = this.game.world.height;
        //console.log(this.game.width + " game " + this.game.height);
        //console.log(this.game.world.width + " world " + this.game.world.height);
        //console.log(this.game.stage.width + " stage " + this.game.stage.height);
        ////this.game.stage.updateTransform();
        //this.game.camera.reset();
        //this.game.camera.setBoundsToWorld();
        //this.game.camera.setSize(100,100);
        //this.game.camera.deadzone = new Phaser.Rectangle(100, 100, this.game.stage.width - 200, this.game.stage.height - 200); //镜头跟随触发区域dead zone
        //this.game.camera.follow(this.hero); //摄像机跟随人物
        //this.game.camera.setPosition(this.hero.x,this.hero.y);
        //this.game.camera.update();
    },
    resize: function (width, height) {
        this.game.world.resize(1280,1280);
        this.game.world.setBounds(0, 0, 1280,1280);
        //this.game.stage.width = this.game.world.width;
        //this.game.stage.height = this.game.world.height;
        console.log(width + " resize " + height);
        console.log(this.game.width + " game " + this.game.height);
        console.log(this.game.world.width + " world " + this.game.world.height);
        console.log(this.game.stage.width + " stage " + this.game.stage.height);
        //this.layer.resize(width, height);
        //this.layer._bounds.width=width;
        //this.layer._bounds.height=height;
        //this.layer.resizeWorld(); //调整世界的范围跟地图范围相同
        //console.log(game.world.width+" world "+game.world.height);
        //game.world.resize(game.world.width, game.world.height);
        //game.world.setBounds(0,0,game.world.width, game.world.height);
        //console.log(game.world);
        //console.log(" world getBounds"+game.world.getBounds());
        //console.log(" world getLocalBounds"+game.world.getLocalBounds());
        //this.game.camera.checkBounds();
        //this.game.camera.focusOn(this.hero);
        //this.game.camera.deadzone = new Phaser.Rectangle(100, 100, this.game.stage.width - 200, this.game.stage.height - 200); //镜头跟随触发区域dead zone
        //game.camera.update();
        this.game.camera.follow(this.hero); //摄像机跟随人物
        //game.stage.width=width;
        //game.stage.height=height;
        //console.log(game.stage.width+" : "+game.stage.height);
        //game.world.resize(height, width);
        //this.game.scale.refresh();
    }
});