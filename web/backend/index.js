const express = require('express');
const cors = require('cors');
const axios = require('axios');
const app = express();
const bodyParser = require('body-parser');
const dotenv = require('dotenv');

dotenv.config();

app.use(cors());
app.use(bodyParser.json());

app.post('/logs', (req, res) => {
    const logMessage = req.body.message;

    try {
        axios.post(process.env.LOG_COLLECTOR_MICROSERVICE_URI, {
            message: logMessage
        });
    } catch (error) {
        console.error('Error processing log:', error);
        return res.status(500).send('Error processing log');
    }
    console.log('Received log:', logMessage);
    res.sendStatus(200);
});

const PORT = process.env.PORT || 8000;
app.listen(PORT, () => {
    console.log(`Server running on port ${PORT}`);
});
