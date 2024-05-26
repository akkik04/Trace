// index.js
const express = require('express');
const cors = require('cors');
const app = express();
const bodyParser = require('body-parser');

app.use(cors()); // Enable CORS for all routes
app.use(bodyParser.json());

app.post('/logs', (req, res) => {
    const logMessage = req.body.message;
    console.log('Received log:', logMessage);
    res.sendStatus(200);
});

const PORT = process.env.PORT || 3000;
app.listen(PORT, () => {
    console.log(`Server running on port ${PORT}`);
});
