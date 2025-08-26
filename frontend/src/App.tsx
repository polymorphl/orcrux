import { useEffect, useState } from 'react';
import './App.css';
import logo from './assets/images/logo-text.png';
import Balatro from './components/Balatro';
import Wizard from './components/Wizard';
import { splitIdleColors } from './lib/colors';

function App() {
  const [color1, setColor1] = useState(splitIdleColors[0])
  const [color2, setColor2] = useState(splitIdleColors[1])

  const handleColorChange = (color1: string, color2: string) => {
    setColor1(color1)
    setColor2(color2)
  }

  // use local message to send data to the parent window
  useEffect(() => {
    const handleMessage = (event: MessageEvent) => {
      if (event.data.type === 'color-change') {
        handleColorChange(event.data.color1, event.data.color2)
      }
    }
    window.addEventListener('message', handleMessage)
    return () => window.removeEventListener('message', handleMessage)
  }, [])

  return (
    <div className="relative">
      <div id="app" className="text-white min-h-screen overflow-hidden">
        <Balatro
          isRotate={false}
          mouseInteraction={false}
          pixelFilter={600}
          color1={color1}
          color2={color2}
        />
        <div className="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 w-[620px] h-[620px] flex flex-col items-center justify-center">
          <div className='bg-gradient-to-br from-crystal-700/40 to-crystal-600/30 rounded-sm border border-crystal-500/20 shadow-lg backdrop-blur-sm w-full h-full'>
            <img src={logo} id="logo" alt="logo" className='w-[200px] mb-4 mx-auto mt-4' />
            <p className="text-crystal-100 text-sm leading-relaxed text-center mb-2 px-4">
              Secure secret sharing using Shamir's Secret Sharing algorithm.
              <br />
              Split your secrets into multiple shards or recompose them back together.
            </p>
            <hr className="mt-4 border-crystal-500/20" />
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
