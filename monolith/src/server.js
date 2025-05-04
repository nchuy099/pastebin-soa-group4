const express = require('express');
const cors = require('cors');
const morgan = require('morgan');
const path = require('path');
const pasteService = require('./services/paste.service');
const { removeExpiredPastes } = require('./services/cleanup.service');
require('dotenv').config();

const app = express();

// Middleware
app.use(cors());
app.use(express.json());
app.use(express.urlencoded({ extended: true }));
app.use(morgan('dev'));

// Set view engine
app.set('view engine', 'ejs');
app.set('views', path.join(__dirname, 'views'));

// Routes
app.use('/', require('./routes/paste.routes'));

// Error handling middleware
app.use((err, req, res, next) => {
    console.error(err.stack);
    res.status(500).json({
        success: false,
        message: 'Internal Server Error',
        error: process.env.NODE_ENV === 'development' ? err.message : undefined
    });
});

// Start server
const PORT = process.env.PORT || 3000;
const CLEANUP_INTERVAL_MINS = process.env.CLEANUP_INTERVAL_MINS || 1;

app.listen(PORT, () => {
    console.log(`Server is running on port ${PORT}`);

    // Run cleanup service every hour
    setInterval(removeExpiredPastes, CLEANUP_INTERVAL_MINS * 60 * 1000);
});