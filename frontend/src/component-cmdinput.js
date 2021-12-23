import {useEffect, useRef, useState} from "react";
import axios from "axios";

function CmdInput(props) {
	const [history, setHistory] = useState([]);
	const [idx, setIdx] = useState(0);
	const inputRef = useRef(null);
	
	const submit = (event) => {
		if (event.key === "Enter") {
			console.log("send command: "+event.target.value);
			axios.post("http://localhost:8000/message", event.target.value);	
			setHistory([...history, event.target.value]);
			setIdx(history.length);
			event.target.value = "";
		}
		if (event.keyCode === 38) {
			if (idx > 0) {
				var tmp = idx - 1
				setIdx(tmp);
			}
			event.target.value = history[idx];
		}
		if (event.keyCode === 40) {
			if (idx < history.length) {
				var tmp = idx + 1
				setIdx(tmp);
				console.log(history[idx]);
				console.log(idx);
				event.target.value = history[idx];
			} else {
				event.target.value = "";
			}
		}
	}
	
	useEffect(() => {
		inputRef.current.focus()
	}, []);

	return (
		<div className="CmdInput">
			<input ref={inputRef} className="CmdTextInput" type="text" onKeyUp={submit}/>
		</div>
	);
}

export default CmdInput;
