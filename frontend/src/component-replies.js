import {useState, useEffect, useRef} from "react";

function ReplyWindow() {
	const replyWindowRef = useRef(null)
	//const ws = new WebSocket("ws://192.168.1.178:8000/admin");
	const ws = new WebSocket("ws://localhost:8000/admin");

	const scrollToBottom = () => {
		replyWindowRef.current.scroll(0,replyWindowRef.current.scrollHeight);
	}

	return (
    <div className="ReplyWindow" ref={replyWindowRef}>
		<Replies websocket={ws} scrollToBottom={scrollToBottom} />	
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
		props.scrollToBottom();
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
		<div className="ReplyHeader"> 
			<HostInfo props={props.HostInfo}/>
			<div className="Command">{props.Command}</div>
		</div>
		<div className="Output">
			{props.Result}{props.Err}
		</div>
	</div>
	);
}

function HostInfo({props}) {
	return (
	<div className="HostInfo">
		[{props.username}@{props.hostname} on {props.address}]
	</div>
	);
}

export default ReplyWindow;
