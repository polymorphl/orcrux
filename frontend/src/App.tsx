import './App.css';
import logo from './assets/images/logo-text.png';
import Balatro from './components/Balatro';
import Wizard from './components/Wizard';

function App() {
  return (
    <div className="relative">
      <div id="app" className="text-white min-h-screen overflow-hidden">
        <Balatro
          isRotate={false}
          mouseInteraction={false}
          pixelFilter={600}
          color1="#0891B2"
          color2="#7C3AED"
        />
        <div className="absolute top-0 left-0 w-full h-full flex flex-col items-center justify-center">
          <div className='bg-gradient-to-br from-crystal-700/40 to-crystal-600/30 rounded-sm border border-crystal-500/20 shadow-lg backdrop-blur-sm'>
            <img src={logo} id="logo" alt="logo" className='w-[200px] mb-4 mx-auto mt-4' />
            <p className="text-crystal-100 text-sm leading-relaxed text-center mb-2 px-4">
              Secure secret sharing using Shamir's Secret Sharing algorithm.
              <br />
              Split your secrets into multiple shards or recompose them back together.
            </p>
            <div className="flex flex-col gap-4 items-center justify-center p-4 overflow-hidden">
              <Wizard />
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default App
