'use strict';
var PlayGameState = require('./PlayGameState.js');
var guiLogin = require('../gui/login.js');
var mqant=window.mqant
module.exports = {
    preload: function preload() {
        "use strict";
        this.game.load.image('loading', 'assets/ball.png');
        this.game.load.tilemap('desert', 'assets/tilemaps/desert.json', null, Phaser.Tilemap.TILED_JSON);
        this.game.load.image('tiles', 'assets/tilemaps/tmw_desert_spacing.png');
    },
    create: function create() {
        var self=this;
        var preloadSprite = this.game.add.sprite(34, game.height / 2, 'loading');
        this.game.load.setPreloadSprite(preloadSprite);
        //this.game.scale.scaleMode = Phaser.ScaleManager.SHOW_ALL;
        this.game.scale.scaleMode = Phaser.ScaleManager.RESIZE;
        this.game.scale.setUserScale(Phaser.myScaleManager.hScale, Phaser.myScaleManager.vScale, Phaser.myScaleManager.hTrim, Phaser.myScaleManager.vTrim);
        var useSSL = 'https:' == document.location.protocol ? false : true;
        try{
            mqant.init({
                host: window.location.hostname,
                port: 3653,
                client_id: "111",
                useSSL:useSSL,
                onSuccess:function() {
                    //alert("游戏链接成功!");
                    mqant.requestNR("Hitball/HD_Join",{
                        "Rid": "001",
                    });

                    mqant.on('Hitball/OnEnter', function(data) {
                        var message=JSON.parse(data.payloadString);
                        var player=message.Player;
                        var coins=message.Coins;
                        var Rid=message.Rid;
                        try{
                            self.game.state.add('PlayGameState', new PlayGameState(Rid,player,coins), false);
                            self.game.state.start('PlayGameState');
                        }catch(e) {
                            alert(e);
                        }
                    });
                },
                onConnectionLost:function(code,reason) {
                    console.log(code)
                    alert("链接断开了:"+code);
                }
            });
        }catch (e){
            alert(e);
        }
    }
};
