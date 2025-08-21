import './App.css';
import logo from './assets/images/logo-universal.png';
import Wizard from './components/Wizard';

function App() {
  return (
    <div id="app" className="bg-crystal-700 text-white min-h-screen bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900 p-4 overflow-hidden">
      <img src={logo} id="logo" alt="logo" className='w-1/6 mx-auto mt-4' />
      <p className="text-crystal-100 text-sm leading-relaxed">
        Secure secret sharing using Shamir's Secret Sharing algorithm.
        <br />
        Split your secrets into multiple shards or recompose them back together.
      </p>
      <div className="flex flex-col gap-4 items-center justify-center p-4 overflow-hidden">
        <Wizard />
      </div>
    </div>
  )
}

export default App
