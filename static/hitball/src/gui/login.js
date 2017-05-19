'use strict';

/**
 * Created by liangdas on 16/12/19.
 * Email 1587790525@qq.com
 */
module.exports ={
    id: 'myWindow',

    component: 'Window',

    padding: 4,

    //component position relative to parent
    position: { x: 10, y: 10 },

    width: 500,
    height: 500,

    layout: [1, 5],
    children: [null, {
        id: 'username',
        text: 'liangdas',
        component: 'Input',
        position: 'center',
        width: 300,
        height: 50
    }, {
        id: 'passwd',
        text: '123456',
        component: 'Input',
        position: 'center',
        width: 300,
        height: 50
    }, {
        id: 'warName',
        text: 'ys',
        component: 'Input',
        position: 'center',
        width: 300,
        height: 50
    }, {
        id: 'btn1',
        text: 'Get Text Value',
        component: 'Button',
        position: 'center',
        width: 200,
        height: 100
    }]
};