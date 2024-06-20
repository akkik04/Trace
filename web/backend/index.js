// index.js
const express = require('express');
const cors = require('cors');
const axios = require('axios');
const app = express();
const bodyParser = require('body-parser');

app.use(cors());
app.use(bodyParser.json());

app.post('/logs', (req, res) => {
    const logMessage = req.body.message;
    
    try{
        axios.post('http://log-collector-microservice:8080/log_collector', {
            message: logMessage
        });
    }catch (error) {
        console.error('Error processing log:');
        res.status(500).send('Error processing log');
    }
    console.log('Received log:', logMessage);
    res.sendStatus(200);
});

const PORT = process.env.PORT || 8000;
app.listen(PORT, () => {
    console.log(`Server running on port ${PORT}`);
});