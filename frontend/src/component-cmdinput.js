import {useEffect, useRef} from "react";
import axios from "axios";

function CmdInput(props) {
	const inputRef = useRef(null);
	
	const submit = (event) => {
		if (event.key === "Enter") {
			//console.log("send command: "+event.target.value);
			axios.post("http://localhost:8000/message", event.target.value);	
			event.target.value = ""
		}
	}
	
	useEffect(() => {
		inputRef.current.focus()
	}, []);

	return (
		<div className="CmdInput">
			<input ref={inputRef} className="CmdTextInput" type="text" onKeyPress={submit}/>
		</div>
	);
}

export default CmdInput;
