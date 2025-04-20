const db = require('../config/database');
const { nanoid } = require('nanoid');

class Paste {
    static async create({ content, title, language, expiresIn, visibility }) {
        let id;
        let isUnique = false;
        let attempts = 0;
        const maxAttempts = 5; // Prevent infinite loop

        while (!isUnique && attempts < maxAttempts) {
            id = nanoid(8);
            // Check if ID already exists
            const [existing] = await db.query('SELECT id FROM paste WHERE id = ?', [id]);
            if (existing.length === 0) {
                isUnique = true;
            } else {
                attempts++;
            }
        }

        if (!isUnique) {
            throw new Error('Failed to generate unique ID after multiple attempts');
        }

        let expiresAt = null;
        if (expiresIn && !isNaN(expiresIn)) {
            expiresAt = new Date(Date.now() + expiresIn * 60 * 1000); // Set expiration time
        }

        // Insert into paste table
        await db.query(
            'INSERT INTO paste (id, content, title, language, expires_at, visibility) VALUES (?, ?, ?, ?, ?, ?)',
            [id, content, title || 'Untitled', language || 'text', expiresAt, visibility || 'PUBLIC']
        );

        return { id };
    }



    static async getById(id) {
        // Get the paste
        const [pastes] = await db.query(
            'SELECT id, content, title, language, created_at, expires_at, visibility FROM paste WHERE id = ?',
            [id]
        );

        if (pastes.length === 0 || (pastes[0].expires_at && new Date(pastes[0].expires_at) <= new Date())) {
            return { status: 'not_found' };
        }

        // Insert into paste_views table to track the view, using CURRENT_TIMESTAMP
        await db.query(
            'INSERT INTO paste_views (paste_id, viewed_at) VALUES (?, CURRENT_TIMESTAMP)',
            [id]
        );

        return { status: 'active', paste: pastes[0] };
    }


    static async getPublic(page = 1, limit = 5) {
        // Get only non-expired public pastes with pagination
        const offset = (page - 1) * limit;

        // Get pastes with pagination
        const [pastes] = await db.query(
            'SELECT id, title, language, created_at FROM paste WHERE visibility = ? AND (expires_at IS NULL OR expires_at > NOW()) ORDER BY created_at DESC LIMIT ? OFFSET ?',
            ['public', limit, offset]
        );

        // Get total count for pagination
        const [countResult] = await db.query(
            'SELECT COUNT(*) as total FROM paste WHERE visibility = ? AND (expires_at IS NULL OR expires_at > NOW())',
            ['public']
        );

        const totalPastes = countResult[0].total;
        const totalPages = Math.ceil(totalPastes / limit);

        return {
            pastes: pastes,
            pagination: {
                page,
                limit,
                totalPastes,
                totalPages
            }
        };
    }

    static async getMonthlyViewStats(month) {
        const fullDate = `${month}-01`; // Chuyển 'YYYY-MM' thành 'YYYY-MM-01'

        const [[viewStats]] = await db.query(`
            SELECT
                COALESCE(SUM(view_count), 0) AS totalViews,
                COALESCE(ROUND(AVG(view_count)), 0) AS avgViewsPerPaste,
                COALESCE(MIN(view_count), 0) AS minViews,
                COALESCE(MAX(view_count), 0) AS maxViews
            FROM (
                SELECT paste_id, COUNT(*) AS view_count
                FROM paste_views
                WHERE YEAR(viewed_at) = YEAR(?) AND MONTH(viewed_at) = MONTH(?)
                GROUP BY paste_id
            ) AS monthly_paste_views
        `, [fullDate, fullDate]);

        return { month, ...viewStats };
    }


}

module.exports = Paste;