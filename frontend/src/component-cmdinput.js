import axios from "axios";

function CmdInput(props) {
	
	const submit = (event) => {
		if (event.key === "Enter") {
			console.log("send command: "+event.target.value);
			axios.post("http://localhost:8000/message", event.target.value);	
			event.target.value = ""
		}
	}

	return (
		<div className="CmdInput">
			<input type="text" onKeyPress={submit}/>
		</div>
	);
}

export default CmdInput;
