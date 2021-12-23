import {useState, useEffect} from "react";

function ReplyWindow() {
	const ws = new WebSocket("ws://localhost:8000/admin");

	return (
    <div className="ReplyWindow">
		<Replies websocket={ws}/>	
    </div>
  	);
}


function Replies(props) {
	const [messages, setMessages] = useState([]);

	useEffect(() => {
		props.websocket.onmessage = msg => {
			var message = JSON.parse(msg.data);
			console.log(message)
			setMessages((existingMessages) => [...existingMessages, message]);
		}
	});
	
	return (
    <div className="Replies">
		{messages.map((m) => <Reply key={m.HostInfo.address}  props={m}/> )}
    </div>
  	);
}

function Reply({props}) {
	return (
	<div className="Reply">
		<HostInfo props={props.HostInfo}/>
	    {props.Command}
		<div className="Output">
			{props.Result}{props.Err}
		</div>
	</div>
	);
}

function HostInfo({props}) {
	return (
	<div className="Host">
		[{props.username}@{props.hostname} on {props.address}] $
	</div>
	);
}

export default ReplyWindow;
