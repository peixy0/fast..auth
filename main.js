const app = require('express')()
const http = require('http').Server(app)
const io = require('socket.io')(http)

requests = {}

app.get('/api/fast..auth', (req, res) => {
    let requestId = req.query.r
    let token = req.header('X-Auth-Token')
    if (!requestId || !token) {
        res.status(400).end()
        return
    }
    if (!requests[requestId]) {
        res.status(404).end()
        return
    }
    requests[requestId].emit('fast..auth..token', token)
    requests[requestId].disconnect()
    delete requests[requestId]
    res.send('ok')
})

io.on('connection', (socket) => {
    let requestId = Math.floor(Math.random() * (36 ** 8)).toString(36)
    requests[requestId] = socket
    socket.emit('fast..auth..request', requestId)
    setTimeout(() => {
        delete requests[requestId]
    }, 60 * 1000)
})

http.listen(8003, () => {
    console.log('fast..auth working on 8003');
})
