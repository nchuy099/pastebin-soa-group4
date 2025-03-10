const express = require('express');
const cors = require('cors');
const morgan = require('morgan');
const mysql = require('mysql2/promise');
const { nanoid } = require('nanoid');
const path = require('path');
require('dotenv').config();

const app = express();

// Middleware
app.use(cors());
app.use(express.json());
app.use(express.urlencoded({ extended: true }));
app.use(morgan('dev'));

// View engine setup
app.set('views', path.join(__dirname, 'views'));
app.set('view engine', 'ejs');

// MySQL connection configuration
const dbConfig = {
    host: process.env.DB_HOST || 'localhost',
    user: process.env.DB_USER || 'root',
    password: process.env.DB_PASSWORD || '',
    database: process.env.DB_NAME || 'pastebin'
};

// Database connection pool
const pool = mysql.createPool(dbConfig);

// Cleanup expired pastes (runs every minute)
setInterval(async () => {
    try {
        const [result] = await pool.query('DELETE FROM pastes WHERE expires_at IS NOT NULL AND expires_at < NOW()');
        console.log(`Cleaned up ${result.affectedRows} expired pastes`);
    } catch (err) {
        console.error('Error cleaning up expired pastes:', err);
    }
}, 60 * 1000);

// Routes
app.get('/', (req, res) => {
    res.render('index', { paste: null, error: null });
});

// Create a new paste
app.post('/paste', async (req, res) => {
    const { content, title, language, expires_in, visibility } = req.body;
    if (!content) {
        return res.render('index', { paste: null, error: 'Content is required' });
    }

    const id = nanoid(8);
    let expiresAt = null;
    if (expires_in && !isNaN(expires_in)) {
        expiresAt = new Date(Date.now() + expires_in * 60 * 1000);
    }

    try {
        await pool.query(
            'INSERT INTO pastes (id, content, title, language, expires_at, visibility) VALUES (?, ?, ?, ?, ?, ?)',
            [id, content, title || 'Untitled', language || 'text', expiresAt, visibility || 'public']
        );
        res.redirect(`/paste/${id}`);
    } catch (err) {
        console.error('Error creating paste:', err);
        res.render('index', { paste: null, error: 'Failed to create paste' });
    }
});

// Get a paste by ID
app.get('/paste/:id', async (req, res) => {
    const { id } = req.params;

    try {
        const [rows] = await pool.query('SELECT * FROM pastes WHERE id = ?', [id]);
        if (rows.length === 0) {
            return res.status(404).render('index', { paste: null, error: 'Paste not found' });
        }

        const paste = rows[0];
        if (paste.expires_at && new Date(paste.expires_at) < new Date()) {
            return res.status(404).render('index', { paste: null, error: 'Paste has expired' });
        }

        // Increment view count
        await pool.query('UPDATE pastes SET views = views + 1 WHERE id = ?', [id]);
        paste.views += 1;

        res.render('paste', { paste });
    } catch (err) {
        console.error('Error retrieving paste:', err);
        res.status(500).render('index', { paste: null, error: 'Server error' });
    }
});

// Get all public pastes
app.get('/public', async (req, res) => {
    try {
        const [rows] = await pool.query(
            'SELECT * FROM pastes WHERE visibility = "public" AND (expires_at IS NULL OR expires_at > NOW()) ORDER BY created_at DESC LIMIT 20'
        );
        res.render('public', { pastes: rows });
    } catch (err) {
        console.error('Error fetching public pastes:', err);
        res.status(500).render('index', { paste: null, error: 'Server error' });
    }
});

// Analytics endpoint
app.get('/analytics/:month', async (req, res) => {
    const { month } = req.params; // Format: YYYY-MM
    try {
        const [rows] = await pool.query(
            `SELECT 
                COUNT(*) as total_pastes,
                SUM(views) as total_views,
                AVG(views) as avg_views_per_paste
             FROM pastes 
             WHERE DATE_FORMAT(created_at, '%Y-%m') = ?`,
            [month]
        );
        res.json({
            month,
            statistics: rows[0]
        });
    } catch (err) {
        console.error('Error fetching analytics:', err);
        res.status(500).json({ error: 'Server error' });
    }
});

// Error handling middleware
app.use((err, req, res, next) => {
    console.error(err.stack);
    res.status(500).json({
        success: false,
        message: 'Internal Server Error',
        error: process.env.NODE_ENV === 'development' ? err.message : undefined
    });
});

const PORT = process.env.PORT || 3000;
app.listen(PORT, () => {
    console.log(`Server is running on http://localhost:${PORT}`);
}); 