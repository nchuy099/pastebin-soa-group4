const mysql = require('mysql2/promise');

class PublicPasteModel {
    constructor() {
        this.pool = mysql.createPool({
            host: process.env.DB_HOST,
            user: process.env.DB_USER,
            password: process.env.DB_PASSWORD,
            database: process.env.DB_NAME,
            connectionLimit: 50,
            waitForConnections: true
        });
    }

    async getPublicPastes() {
        const [rows] = await this.pool.query(
            `SELECT id, content, title, language, 
              created_at, expires_at, views
       FROM pastes 
       WHERE visibility = 'public' 
         AND (expires_at IS NULL OR expires_at > NOW())
       ORDER BY created_at DESC 
       LIMIT 10`
        );
        return rows;
    }
}

module.exports = new PublicPasteModel();