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
            const [existing] = await db.query('SELECT id FROM pastes WHERE id = ?', [id]);
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
            expiresAt = new Date(Date.now() + expiresIn * 60 * 1000);
        }

        await db.query(
            'INSERT INTO pastes (id, content, title, language, expires_at, visibility) VALUES (?, ?, ?, ?, ?, ?)',
            [id, content, title || 'Untitled', language || 'text', expiresAt, visibility || 'public']
        );

        return { id };
    }

    static async getById(id) {
        // Get the paste
        const [pastes] = await db.query(
            'SELECT id, content, title, language, created_at, expires_at, views, visibility, status FROM pastes WHERE id = ?',
            [id]
        );

        if (pastes.length === 0) {
            return { status: 'not_found' };
        }

        // Check if paste is expired
        if (pastes[0].expires_at && new Date(pastes[0].expires_at) <= new Date()) {
            return { status: 'expired', paste: pastes[0] };
        }

        // Increment views for active paste
        await db.query('UPDATE pastes SET views = views + 1 WHERE id = ?', [id]);

        return { status: 'active', paste: pastes[0] };
    }

    static async getPublic() {
        // Get only non-expired public pastes
        const [pastes] = await db.query(
            'SELECT id, content, title, language, created_at, expires_at, views, visibility FROM pastes WHERE visibility = ? AND (expires_at IS NULL OR expires_at > NOW()) ORDER BY created_at DESC LIMIT 10',
            ['public']
        );
        return pastes;
    }

    static async getMonthlyStats(month) {
        const [stats] = await db.query(`
            SELECT 
                COUNT(*) as totalPastes,
                COALESCE(SUM(views), 0) as totalViews,
                COALESCE(ROUND(AVG(views)), 0) as avgViewsPerPaste,
                COALESCE(MIN(views), 0) as minViews,
                COALESCE(MAX(views), 0) as maxViews,
                SUM(CASE WHEN expires_at IS NULL OR expires_at > NOW() THEN 1 ELSE 0 END) as activePastes,
                SUM(CASE WHEN expires_at IS NOT NULL AND expires_at <= NOW() THEN 1 ELSE 0 END) as expiredPastes
            FROM pastes 
            WHERE DATE_FORMAT(created_at, '%Y-%m') = ?
        `, [month]);

        return { month, ...stats[0] };
    }
}

module.exports = Paste; 