const mysql = require('mysql2/promise');

const pool = mysql.createPool({
    host: process.env.DB_HOST || 'localhost',
    user: process.env.DB_USER || 'root',
    password: process.env.DB_PASSWORD || '',
    database: process.env.DB_NAME || 'pastebin',
    waitForConnections: true,
    connectionLimit: 10,  // Reduce to prevent CPU overhead
    queueLimit: 30        // Control how many requests are queued
});

module.exports = pool;