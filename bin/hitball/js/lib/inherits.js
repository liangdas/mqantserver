'use strict';
/**
 * Created by liangdas on 2016/12/6 0006.
 * Email :1587790525@qq.com
 * 这是一个 javascript 类继承函数,
 * 与以往的继承函数不同的是这个函数,
 * 可以直接从一个已存在的普通构造函数继承。
 */
var  extend= function (superCtor, prop) {
    return function () {
        var fnTest = /xyz/.test(function () {
            xyz;
        }) ? /\b_super\b/ : /.*/;

        var _super = superCtor.prototype;
        //// The base Class implementation (does nothing)
        function baseClass() {
            if (baseClass.prototype.ctor) {
                baseClass.prototype.ctor.apply(this, Array.from(arguments));
            }
        };
        // 空函数F:
        var F = function F() {};
        // 把F的原型指向Student.prototype:
        F.prototype = superCtor.prototype;
        //F.prototype = Object.create(superCtor.prototype);
        // 把PrimaryStudent的原型指向一个新的F对象，F对象的原型正好指向Student.prototype:
        baseClass.prototype = new F();
        // 把PrimaryStudent原型的构造函数修复为PrimaryStudent:
        baseClass.prototype.constructor = baseClass;

        var prototype = baseClass.prototype;
        if (!_super.ctor) {
            _super.ctor = superCtor;
        }

        // Copy the properties over onto the new prototype
        for (var name in prop) {
            // Check if we're overwriting an existing function
            prototype[name] = typeof prop[name] == "function" && typeof _super[name] == "function" && fnTest.test(prop[name]) ? function (name, fn) {
                return function () {
                    //var tmp_superclass = this.superclass;
                    var tmp = this._super;
                    // Add a new ._super() method that is the same method
                    // but on the super-class
                    this._super = _super[name];
                    //this.superclass=superCtor.bind(this);
                    // The method only need to be bound temporarily, so we
                    // remove it when we're done executing
                    var ret = fn.apply(this, Array.from(arguments));
                    this._super = tmp;
                    //this.superclass=tmp_superclass;
                    return ret;
                };
            }(name, prop[name]) : prop[name];
        }
        return baseClass;
    }();
};

window.extend=extend;