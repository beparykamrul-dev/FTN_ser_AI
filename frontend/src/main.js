const root = document.querySelector('#root')

root.innerHTML = `
  <main style="font-family:system-ui;max-width:900px;margin:40px auto;padding:20px">
    <h1>FTN SER AI</h1>
    <p id="status">Checking backend…</p>
    <pre id="output"></pre>
  </main>
`

async function check() {
  const status = document.querySelector('#status')
  const output = document.querySelector('#output')
  try {
    const r = await fetch('/healthz')
    const data = await r.json()
    status.textContent = `Backend: ${data.status}`
    output.textContent = JSON.stringify(data, null, 2)
  } catch (e) {
    status.textContent = 'Backend unavailable'
    output.textContent = String(e)
  }
}
check()
