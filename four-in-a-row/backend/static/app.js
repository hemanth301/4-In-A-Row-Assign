const connectBtn = document.getElementById('connect');
const usernameInput = document.getElementById('username');
const boardDiv = document.getElementById('board');
let socket;
let gameState;

connectBtn.onclick = () => {
  const username = usernameInput.value || 'anon'
  socket = new WebSocket(`ws://${location.host}/ws?username=${username}`)
  socket.onmessage = (ev) => {
    const data = JSON.parse(ev.data)
    gameState = data
    render()
  }
}

function render() {
  if (!gameState) return
  boardDiv.innerHTML = ''
  const table = document.createElement('table')
  for (let r = 5; r >= 0; r--) {
    const tr = document.createElement('tr')
    for (let c = 0; c < 7; c++) {
      const td = document.createElement('td')
      td.style.width = '40px'
      td.style.height = '40px'
      td.style.border = '1px solid #000'
      td.style.textAlign = 'center'
      td.style.cursor = 'pointer'
      const v = gameState.board[r] ? gameState.board[r][c] : 0
      td.innerText = v === 0 ? '' : (v===1? 'X':'O')
      td.onclick = () => { drop(c) }
      tr.appendChild(td)
    }
    table.appendChild(tr)
  }
  boardDiv.appendChild(table)
}

function drop(col) {
  socket.send(JSON.stringify({action: 'drop', column: col}))
}
