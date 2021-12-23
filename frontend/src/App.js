import './App.css';
import CmdInput from './component-cmdinput';
import ReplyWindow from './component-replies';

function App() {
	
	return (
    	<div className="App">
		<AppHeader/>
		<ReplyWindow/>
		<CmdInput/>
    	</div>
  	);
}

function AppHeader() {
	return <div className="AppHeader">cmd+ctrl</div>
}

export default App;
