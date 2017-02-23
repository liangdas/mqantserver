#快速使用
获取 mqantserver：

	git clone https://github.com/liangdas/mqantserver

设置 mqantserver 目录到 GOPATH 后获取相关依赖：

	go get github.com/astaxie/beego
	go get github.com/gorilla/websocket
	go get github.com/streadway/amqp
	go get github.com/liangdas/mqant

编译 mqantserver：

go install server
如果一切顺利，运行 bin/server 你可以获得以下输出：

	[release] mqant 1.0.0 starting up
	[debug  ] RPCClient create success type(Gate) id(127.0.0.1:Gate)
	[debug  ] RPCClient create success type(Login) id(127.0.0.1:Login)
	[debug  ] RPCClient create success type(Chat) id(127.0.0.1:Chat)
	[release] MySelfHost 172.16.8.4
	[release] WS Listen :%!(EXTRA string=0.0.0.0:3653)
	[release] TCP Listen :%!(EXTRA string=0.0.0.0:3563)

敲击 Ctrl + C 关闭游戏服务器，服务器正常关闭输出：

	[debug  ] RPCServer close success id(127.0.0.1:Chat)
	[debug  ] RPCServer close success id(127.0.0.1:Login)
	[debug  ] RPCServer close success id(127.0.0.1:Gate)
	[debug  ] RPCClient close success type(Gate) id(127.0.0.1:Gate)
	[debug  ] RPCClient close success type(Login) id(127.0.0.1:Login)
	[debug  ] RPCClient close success type(Chat) id(127.0.0.1:Chat)
	[release] mqant closing down (signal: interrupt)

#启动网页版本客户端
编译 mqantserver：

go install client

如果一切顺利，运行 bin/client

访问地址为：http://127.0.0.1:8080/mqant/index.html


#启动python版本客户端

执行src/client/mqtt_chat_client.py即可 需要安装paho.mqtt库,请自行百度

#Demo演示说明
	1. 启动服务器
	2. 启动网页客户端	(默认房间abc 用户名 liangdas)
	3. 启动python客户端 (默认房间abc 用户名 有用户输入)
	4. 登陆成功后就可以聊天了

>mqtt_chat_client.py启动后会提醒输入用户名，而网页版本的客户端进入以后直接使用"liangdas"作为用户名,因此如果想只使用网页版本客户端演示通信的话客户修改一下用户名,防止重名

#项目目录结构
https://github.com/liangdas/mqantserver 仓库中包含了mqant框架,所用到的第三方库,聊天Demo服务端,聊天代码客户端代码

	src
		|-client
			|-mqtt_chat_client.py 聊天客户端 Python版本
			|-webclient.go			聊天客户端网页版本
		|-github.com                需要执行 go get 命令拉取
			|-astaxie.beego框架 		webclient.go用到了
			|-gorilla.websocket		websocket框架
			|-liangdas.mqant			mqant框架代码
			|-streadway.amqp			rabbitmq通信框架
		|-server						聊天服务器Demo
			|-chat						聊天模块
			|-conf						系统配置文件
			|-login					登陆模块
			|-main.go					服务器启动入口




#原理说明
因为mqant默认网关支持mqtt协议,因此网页客户端和Python都是使用的paho.mqtt的客户端库来做的

##消息路由

mqant默认网关约定

	mqtt数据包
		--topic	作为路由规则
		--payload 作为消息体	json数据
###路由规则
topic 格式组成

	serverType/handler  serverType/handler/msgid

	serverType 服务器的模块名
	handler	 模块提供的函数名
	msgid		 [可选]用来控制这条信息是否需要回复客户端
> 为了让对客户端的函数与后台模块间通信的函数进行区别提供安全性,网关默认约定了 handler名称必须已"Handler_"开头
>
> msgid的作用是用来控制这条信息是否需要回复客户端,如果不加msgid网关只会将消息路由到指定的模块,但不管模块返回信息,也不回复客户端
> 如果msgid存在,则网关会路由到模块,并且将调研得到的模块信息返回给客户端,返回客户端的信息topic与客户端发送的topic相同,这样客户端就可以根据收到信息的topic来判断与哪一条已发信息相匹配。

路由原理

	1. 网关收到mqtt消息,解析topic,得到serverType 和 handler
	2. 根据serverType查询模块的RPCClient,如果查到则通过handler来调用远程模块的相应函数
	3. 根据msgid是否存在来判断是否回复客户端

###模块handler函数定义
	(session map[string]interface{},msg map[string]interface{})(result string,err string)

	session 网关生成的该网络连接的属性,可以通过
	session:=gate.NewSession(m.app,s)
	得到具体的session对象


	msg		客户端发送的payload,是json数据

###模块开发
	type Module interface {
		GetType()(string)	//模块类型
		OnInit(app App,settings *conf.ModuleSettings)
		OnDestroy()
		Run(closeSig chan bool)
	}

只要实现了如上的函数接口都可以被认为是一个模块,普通的模块是不支持远程调用的,如果要支持可以初始化一个

	eg.
	//创建一个远程调用的RPC
	m.GetServer().OnInit(app,settings)
	//注册远程调用的函数
	m.GetServer().RegisterGO("Handler_JoinChat",m.joinChat) //我们约定所有对客户端的请求都以Handler_开头
	m.GetServer().RegisterGO("Handler_Say",m.say) //我们约定所有对客户端的请求都以Handler_开头

###conf.ModuleSettings
模块的配置信息很关键,mqant通过一个server.conf的json文件来配置。
单独一个Module的配置信息主要包括以下几个项

	{
		"Id":"127.0.0.1:Login",
		"Host":"127.0.0.1",
		"Rabbitmq":{...},
		"Settings":{...}
	}
	Id 		全局唯一,不能重复
	Host 	所属服务器的IP地址,在以后做分布式的时候可以利用Host来区分不同服务器启动哪些模块
	Rabbitmq 模块RPC远程通信配置,单机环境部署可以不用填
	Settings	模块自定义配置 例如 Gate模块配置了 TCPAddr 等信息


###RPC使用
mqant RPC本身是一个相对独立的功能,关于RPC的配置这里不做说明,只说在模块中RPC的使用

RPC角色分服务提供者和服务调用者

服务提供者使用module.Server创建服务

	module.Server
	//初始化
	OnInit(app module.App,settings *conf.ModuleSettings)
	//注册服务函数
	RegisterGO(_func string, fn interface{})
	//注册服务函数
	Register(_func string, fn interface{})

> RegisterGO与Register的区别是前者为每一条消息创建一个单独的协程来处理,后者注册的函数共用一个协程来处理所有消息,具体使用哪一种方式可以根据实际情况来定,但Register方式的函数请一定注意不要执行耗时功能,以免引起消息阻塞

服务调用者
	GetRouteServersByType(module Type string)(*module.ServerSession,error)

	module.ServerSession
	//远程调用指定函数
	Call(_func string,params ...interface{})(result interface{},err string)
	//远程调用指定函数,但不需要远程模块返回执行结果
	CallNR(_func string,params ...interface{})(err error)

###RPC原理
	RPC底层使用的消息队列进行通信的,消息请求发送一次,消息回复发送一次
	RPC的超时机制目前采用的是限定超时时间默认是10s,也就是在超时时间内还没有返回的话,调用者就直接默认为超时了,因此服务方一定保证在10s内完成所有的操作,并给调用者回复。
	调用者也可以通过CallNR方法发送无需回复的调用,也就没有超时这一项了

###Session
mqant默认网关提供了一个简单的Session管理,每一个连接对应一个Session,并且每一次远程调用Gate都会把Session里面的参数发送到远程模块,因此模块可以通过给Session设置一些参数来达到区分用户的目的

	gate.Session
	Bind //绑定UserID
	UnBind //解绑UserID
	Set(key string, value string) //设置一个参数 Push()以后才能发送到网关
	Push()	//将模块设置的Session信息发送到网关
	...
	Send(topic  string,body []byte) //给这个链接发送一个消息
	...



