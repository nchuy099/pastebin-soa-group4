const db = require('../config/database');

const removeExpiredPastes = async () => {
    try {
        const [result] = await db.query(
            'DELETE FROM paste WHERE expires_at IS NOT NULL AND expires_at <= NOW()'
        );
        console.log(`Cleanup: Removed ${result.affectedRows} expired pastes`);
        return result.affectedRows;
    } catch (error) {
        console.error('Cleanup error:', error);
    }
};

module.exports = { removeExpiredPastes };
