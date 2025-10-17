const express = require('express');
const path = require('path');
const app = express();

app.use(express.static('.'));

app.get('/', (req, res) => {
  res.sendFile(path.join(__dirname, 'index.html'));
});

app.get('/api/status', (req, res) => {
  res.json({
    broadcaster: { status: 'running', viewers: 2 },
    viewers: [
      { id: 'viewer1', connected: true, port: 8081 },
      { id: 'viewer2', connected: true, port: 8082 }
    ]
  });
});

const PORT = 3000;
app.listen(PORT, '0.0.0.0', () => {
  console.log(`Demo UI running on http://localhost:${PORT}`);
});