import logo from './assets/images/logo-universal.png';
import './App.css';
// import { Greet } from "../wailsjs/go/main/App";
import Wizard from './components/Wizard';

function App() {
  // const updateName = (e: any) => setName(e.target.value);
  // const updateResultText = (result: string) => setResultText(result);

  // function greet() {
  //     Greet(name).then(updateResultText);
  // }

  return (
    <div id="app">
      <img src={logo} id="logo" alt="logo" />
      <Wizard />
    </div>
  )
}

export default App
