var wss;
	
var myName;

var notifyNewActivityInterval;

var retryConnectionInterval;

var originalDocumentTitle = window.document.title;

$(document).ready( function() {
	
	window.onbeforeunload = confirmLeave;
	
	$(window).bind("focus", stopNewActivityNotification);
	
	$("#ChangeNameButton").click(onChangeName);
	
	$("#ChatClear").click(onChatClear);
	
	$("#ChatSend").click(function() {	
		sendChat({keyCode:13});
	});
	
	$("#ChatInput").keypress(sendChat);

	connect();
});

function connect() {

	wss = new WebSocket("wss://"+document.domain+"/chatws/chat?name=User"+Math.floor(Math.random()*100+1));
	
	wss.onopen = function(status) {

		// stop trying to connect
		clearInterval(retryConnectionInterval);
		retryConnectionInterval = null;
	}

	wss.onerror = function(status) {
		// debugger;
		console.log(status);
	}
	
	wss.onclose = function(status) {

		// debugger
		$("#ChatUsers").html("");

		var now = new Date();

		$("#ChatReceive").append("<"+ now.toLocaleTimeString() + "> You have been disconnected<br/>");

		startRetryConnection();
	}
	
	wss.onmessage = function(status) {
		
		messageDispatcher(status);		
	}
}

function startRetryConnection() {

	if(retryConnectionInterval) {
		return;
	}

	retryConnectionInterval = setInterval( function() {
		
		connect();

	}, 3000)
}

function messageDispatcher(status) {
	
	if (!status || !status.data )
	{
		return;
	}
	
	var msg;
	
	try{
	 	msg = JSON.parse(status.data);
	}
	catch(err) {
		return;
	}
	
	if (!msg) {
		return;
	}
	
	console.log(msg);
	
	if( msg.Type == "ServerChatMessage") {
		addChatMessage("<b>"+msg.Name +"</b>" + " said -> " + msg.Chat)
	}
	
	else if(msg.Type == "ServerStatusMessage") {
		addChatMessage(msg.Content);
	}
	
	else if(msg.Type == "ServerClientListMessage" ) {
		users = msg.Content.split(",");
		
		clearUsers();
		
		for( i=0; i< users.length; i++)
		{
			if( users[i].length > 0 ) { 
				addUser(users[i]);
			}
		}
	}
	
	else if(msg.Type == "ServerClientJoinMessage" ) {
		addChatMessage(msg.Name + " has connected");
		addUser(msg.Name);
	}
	
	else if(msg.Type == "ServerClientPartMessage" ) {
		addChatMessage(msg.Name + " has disconnected");
		removeUser(msg.Name);
	}
	
	else if(msg.Type == "ServerSetName") {
		setName(msg.NewName);
	}
	
	else if(msg.Type == "ServerChangeName") {
		clientChangedName(msg.OldName, msg.NewName);
	}
	
	else if(msg.Type == "ServerPingMessage") {
		sendPong(msg.Payload);
	}
}

function confirmLeave() {
	var message = "You are about to get disconnected from Chat.",
		e = e || window.event;
		// For IE and Firefox
	if (e) {
		e.returnValue = message;
		}

		// For Safari
		return message;
}

function onChatClear() {
	$("#ChatReceive").html("");
}

function stopNewActivityNotification() {
	clearInterval(notifyNewActivityInterval);
	notifyNewActivityInterval = null;
	setTimeout( function() {
		window.document.title = originalDocumentTitle;
		}, 100 );
}

function notifyNewActivity() {
	if (window.document.hasFocus() || notifyNewActivityInterval)
	{
		return;
	}
	
	notifyNewActivityInterval = setInterval( function() {
		if( window.document.title == originalDocumentTitle) {
			window.document.title = "New Activity";
		}
		else {
			window.document.title = originalDocumentTitle;
		}
	}, 1000)
}

function clientChangedName(oldName, newName) {
	$("#"+oldName).replaceWith("<div id='"+newName + "'>"+newName + "</div>");
	addChatMessage(oldName + " changed name to "+newName);
}

function setName(newName) {
	$("#"+myName).replaceWith("<div id='" + newName + "'>" + newName + "</div>");
	myName = newName;
	$("#NameInput").val(newName);
	addChatMessage("Your name is now " + newName);
}

function onChangeName() {
	var name = $("#NameInput").val();
	
	var msg = new Object();
	msg.Type = "ChangeName";
	msg.NewName = name;
	
	wss.send(JSON.stringify(msg))
}

function clearUsers() {
	$("#ChatUsers").html("");
}

function addUser(name) {
	
	var chatUsers = $("#ChatUsers");
	chatUsers.append("<div id='" + name + "'>" + name + "</div>");
}

function removeUser(name) {
	$("#"+name).remove();
}

function addChatMessage(msg) {
	var now = new Date();
	
	msg = replaceURLWithHTMLLinks(msg);

	var isAtBottom = isScrollBottom();	
					
	$("#ChatReceive").append("<"+ now.toLocaleTimeString() + "> "+msg + "<br/>");
	
	if (isAtBottom) {
		$("#ChatReceive").animate({scrollTop:$("#ChatReceive")[0].scrollHeight}, 100); 
	}
	
	notifyNewActivity();
}

function sendChat(event) {
	if( event.keyCode != 13)
	{
		return;
	}
	
	var msg = $("#ChatInput").val();
	$("#ChatInput").val("");
	
	if(!msg) {
		return;
	}
	
	var obj = new Object();
	obj.Type = "ClientChat";
	obj.Chat = msg;
			
	wss.send(JSON.stringify(obj));
}

function sendPong(payload) {
	var obj = new Object();
	obj.Type = "ClientPong";
	obj.Payload = payload;
	
	wss.send(JSON.stringify(obj));
}

function replaceURLWithHTMLLinks(text) { 
	var exp = /(\b(https?|ftp|file):\/\/[-A-Z0-9+&@#\/%?=~_|!:,.;]*[-A-Z0-9+&@#\/%=~_|])/i; 
	return text.replace(exp,"<a href='$1' target=\"_blank\">$1</a>"); 
}

function isScrollBottom() { 
		var element = document.getElementById("ChatReceive");
	if( Math.abs(element.scrollTop - element.scrollHeight + element.offsetHeight) < 20 ) {
		return true;
	}
	
	return false;
}
