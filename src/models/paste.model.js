const db = require('../config/database');
const { nanoid } = require('nanoid');

class Paste {
    static async create({ content, title, language, expiresIn, visibility }) {
        const id = nanoid(8);
        let expiresAt = null;
        if (expiresIn && !isNaN(expiresIn)) {
            expiresAt = new Date(Date.now() + expiresIn * 60 * 1000);
        }

        const [result] = await db.query(
            'INSERT INTO pastes (id, content, title, language, expires_at, visibility) VALUES (?, ?, ?, ?, ?, ?)',
            [id, content, title || 'Untitled', language || 'text', expiresAt, visibility || 'public']
        );

        return { id, ...result };
    }

    static async getById(id) {
        // Get the paste
        const [pastes] = await db.query(
            'SELECT *, (CASE WHEN expires_at IS NOT NULL AND expires_at <= NOW() THEN "expired" ELSE "active" END) as status FROM pastes WHERE id = ?',
            [id]
        );

        if (pastes.length === 0) {
            throw new Error('Paste not found');
        }

        // Check if paste is expired
        if (pastes[0].expires_at && new Date(pastes[0].expires_at) <= new Date()) {
            throw new Error('This paste has expired and is no longer accessible');
        }

        // Increment views for active paste
        await db.query('UPDATE pastes SET views = views + 1 WHERE id = ?', [id]);

        return pastes[0];
    }

    static async getPublic() {
        // Get only non-expired public pastes
        const [pastes] = await db.query(
            'SELECT *, (CASE WHEN expires_at IS NOT NULL AND expires_at <= NOW() THEN "expired" ELSE "active" END) as status FROM pastes WHERE visibility = "public" AND (expires_at IS NULL OR expires_at > NOW()) ORDER BY created_at DESC LIMIT 10'
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