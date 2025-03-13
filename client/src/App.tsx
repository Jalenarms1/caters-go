import { useState, useRef } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'

function App() {
  const [count, setCount] = useState(0)
  const emailRef = useRef(null)
  const passwordRef = useRef(null)

  const submitLogin = async () => {
    const resp = await fetch("https://lalocura-go-production.up.railway.app/api/login", {
      method: "POST",
      credentials: 'include',
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({
        email: emailRef.current.value,
        password: passwordRef.current.value
      })
    })

    const data = await resp.json()

    console.log(data);
    if (data.token) {
      console.log(data.token);
      
    }
    
  }

  return (
    <>
      <div className="flex min-h-screen bg-slate-700 justify-center items-center">
        <div className="flex relative items-center flex-col bg-slate-500  z-[1] rounded-md shadow-sm shadow-zinc-700 w-[80vw] h-[75vh]">
          <div className="flex flex-col items-center w-[80%] gap-5 pt-10">
            <p className='text-zinc-100 font-semibold text-2xl mb-5'>Login</p>
            <div className="flex flex-col w-full gap-2">
              <p className='text-sm text-zinc-200'>Email</p>
              <input ref={emailRef} type="text" className='bg-white w-full rounded-sm p-1' placeholder='Email' />
            </div>
            <div className="flex flex-col w-full gap-2">
              <p className='text-sm text-zinc-200'>Password</p>
              <input ref={passwordRef} type="password" className='bg-white w-full rounded-sm p-1' placeholder='Password' />
            </div>

            <button onClick={submitLogin} className='w-full bg-slate-800 text-white p-2 rounded-full active:bg-slate-700 transition-all mt-5'>Submit</button>

          </div>

        </div>

      </div>
    </>
  )
}

export default App
