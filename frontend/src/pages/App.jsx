import { useState } from 'react'
import reactLogo from '../assets/react.svg'
import viteLogo from '../../public/vite.svg'
import '../../public/App.css'

function App() {
  const [count, setCount] = useState(0)

  return (
    <>
      <div>
        <a href="https://vite.dev" target="_blank">
          <img src={viteLogo} className="logo" alt="Vite logo" />
        </a>
        <a href="https://react.dev" target="_blank">
          <img src={reactLogo} className="logo react" alt="React logo" />
        </a>
      </div>
      <h1>  
Pemanfaatan Algoritma BFS dan DFS dalam Pencarian Recipe  pada Permainan
Little Alchemy 2        
      </h1>
      <div className="card">
        <button onClick={() => setCount((count) => count + 1)}>
          count is {count}
        </button>
        <p>
          Lukas Raja Agripa | <code>[REDACTED] </code> | 13523158
        </p>
      </div>
      <p className="read-the-docs">
        BENER MASIH KUAT GA BOS ! 
      </p>
    </>
  )
}

export default App
