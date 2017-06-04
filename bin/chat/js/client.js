var mqant=window.mqant
var username;
var users;
var roomName;
var base = 1000;
var increase = 25;
var reg = /^[a-zA-Z0-9_\u4e00-\u9fa5]+$/;
var LOGIN_ERROR = "There is no server to log in, please wait.";
var LENGTH_ERROR = "Name/Channel is too long or too short. 20 character max.";
var NAME_ERROR = "Bad character in Name/Channel. Can only have letters, numbers, Chinese characters, and '_'";
var DUPLICATE_ERROR = "Please change your name to login.";

util = {
	urlRE: /https?:\/\/([-\w\.]+)+(:\d+)?(\/([^\s]*(\?\S+)?)?)?/g,
	//  html sanitizer
	toStaticHTML: function(inputHtml) {
		inputHtml = inputHtml.toString();
		return inputHtml.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
	},
	//pads n with zeros on the left,
	//digits is minimum length of output
	//zeroPad(3, 5); returns "005"
	//zeroPad(2, 500); returns "500"
	zeroPad: function(digits, n) {
		n = n.toString();
		while(n.length < digits)
		n = '0' + n;
		return n;
	},
	//it is almost 8 o'clock PM here
	//timeString(new Date); returns "19:49"
	timeString: function(date) {
		var minutes = date.getMinutes().toString();
		var hours = date.getHours().toString();
		return this.zeroPad(2, hours) + ":" + this.zeroPad(2, minutes);
	},

	//does the argument only contain whitespace?
	isBlank: function(text) {
		var blank = /^\s*$/;
		return(text.match(blank) !== null);
	}
};

//always view the most recent message when it is added
function scrollDown(base) {
	window.scrollTo(0, base);
	$("#entry").focus();
};

// add message on board
function addMessage(from, target, text, time) {
	var name = (target == '*' ? 'all' : target);
	if(text === null) return;
	if(time == null) {
		// if the time is null or undefined, use the current time.
		time = new Date();
	} else if((time instanceof Date) === false) {
		// if it's a timestamp, interpret it
		time = new Date(time);
	}
	//every message you see is actually a table with 3 cols:
	//  the time,
	//  the person who caused the event,
	//  and the content
	var messageElement = $(document.createElement("table"));
	messageElement.addClass("message");
	// sanitize
	text = util.toStaticHTML(text);
	var content = '<tr>' + '  <td class="date">' + util.timeString(time) + '</td>' + '  <td class="nick">【' + util.toStaticHTML(from) + '】 --》 【' + name + '】: ' + '</td>' + '  <td class="msg-text">' + text + '</td>' + '</tr>';
	messageElement.html(content);
	//the log is the stream that we view
	$("#chatHistory").append(messageElement);
	base += increase;
	scrollDown(base);
};

// show tip
function tip(type, name) {
	var tip,title;
	switch(type){
		case 'online':
			tip = name + ' is online now.';
			title = 'Online Notify';
			break;
		case 'offline':
			tip = name + ' is offline now.';
			title = 'Offline Notify';
			break;
		case 'message':
			tip = name + ' is saying now.'
			title = 'Message Notify';
			break;
	}
	var pop=new Pop(title, tip);
};

// init user list
function initUserList(data) {
	users = data.users;
	for(var i = 0; i < users.length; i++) {
		var slElement = $(document.createElement("option"));
		slElement.attr("value", users[i]);
		slElement.text(users[i]);
		$("#usersList").append(slElement);
	}
};

// add user in user list
function addUser(user) {
	var slElement = $(document.createElement("option"));
	slElement.attr("value", user);
	slElement.text(user);
	$("#usersList").append(slElement);
};

// remove user from user list
function removeUser(user) {
	$("#usersList option").each(
		function() {
			if($(this).val() === user) $(this).remove();
	});
};

// set your name
function setName() {
	$("#name").text(username);
};

// set your room
function setRoom() {
	$("#room").text(roomName);
};

// show error
function showError(content) {
	$("#loginError").text(content);
	$("#loginError").show();
};

// show login panel
function showLogin() {
	$("#loginView").show();
	$("#chatHistory").hide();
	$("#toolbar").hide();
	$("#loginError").hide();
	$("#loginUser").focus();
};

// show chat panel
function showChat() {
	$("#loginView").hide();
	$("#loginError").hide();
	$("#toolbar").show();
	$("entry").focus();
	scrollDown(base);
};


$(document).ready(function() {
	//when first time into chat room.
	showLogin();

	//wait message from the server.
	mqant.on('Chat/OnChat', function(data) {
		var message=JSON.parse(data.payloadString)
		addMessage(message.from, message.target, message.msg);
		$("#chatHistory").show();
		if(message.from !== username)
			tip('message', message.from);
	});

	//update user list
	mqant.on('Chat/OnJoin', function(data) {
		var message=JSON.parse(data.payloadString)
		var user = message.user;
		tip('online', user);
		addUser(user);
	});

	//update user list
	mqant.on('Chat/OnLeave', function(data) {
		var message=JSON.parse(data.payloadString)
		var user = message.user;
		tip('offline', user);
		removeUser(user);
	});


	//handle disconect message, occours when the client is disconnect with servers


	//deal with login button click.
	$("#login").click(function() {
		username = $("#loginUser").attr("value");
		roomName = $('#channelList').val();

		if(username.length > 20 || username.length == 0 || roomName.length > 20 || roomName.length == 0) {
			showError(LENGTH_ERROR);
			return false;
		}

		if(!reg.test(username) || !reg.test(roomName)) {
			showError(NAME_ERROR);
			return false;
		}
		var useSSL = 'https:' == document.location.protocol ? false : true;
		mqant.init({
			host: window.location.hostname,
			port: 3653,
			client_id: "111",
			useSSL:useSSL,
		    onSuccess:function() {
				console.log("onConnected");
				var topic = "Login/HD_Login";
				mqant.request(topic, {
					"userName": username,
					"passWord": "Hello,anyone!"
				}, function(data) {
					var message=JSON.parse(data.payloadString)
					if(message.Error!="") {
						showError(message.Error);
						return;
					}
					mqant.request("Chat/HD_JoinChat",{
						"roomName": roomName
					}, function(data) {
						var message=JSON.parse(data.payloadString)
						if(message.Error!="") {
							showError(message.Error);
							return;
						}
						setName();
						setRoom();
						showChat();
						initUserList(message.Result);
					});


				});
			},
			onConnectionLost:function(reason) {
				showLogin();
			}
			});
	});

	//deal with chat mode.
	$("#entry").keypress(function(e) {
		var topic = "Chat/HD_Say";
		var target = $("#usersList").val();
		if(e.keyCode != 13 /* Return */ ) return;
		var msg = $("#entry").attr("value").replace("\n", "");
		if(!util.isBlank(msg)) {
			mqant.request(topic, {
				roomName: roomName,
				content: msg,
				from: username,
				target: target
			}, function(data) {

				var message=JSON.parse(data.payloadString)
				$("#entry").attr("value", ""); // clear the entry field.
				if(target != '*' && target != username) {
					if(message.Error!="") {
						addMessage(username, target,msg +" error ("+message.Error+")");
					}else{
						addMessage(username, target, msg);
					}

					$("#chatHistory").show();
				}
			});
		}
	});
});