const Paste = require('../models/paste.model');

class PasteService {
    async createPaste(pasteData) {
        try {
            return await Paste.create(pasteData);
        } catch (error) {
            console.error('Create paste error:', error);
            throw new Error('Failed to create paste');
        }
    }

    async getPasteById(id) {
        try {
            const paste = await Paste.getById(id);
            if (!paste) {
                throw new Error('Paste not found');
            }
            return paste;
        } catch (error) {
            console.error('Get paste by id error:', error);
            throw error;
        }
    }

    async getPublicPastes() {
        try {
            return await Paste.getPublic();
        } catch (error) {
            console.error('Get public pastes error:', error);
            throw new Error('Failed to fetch public pastes');
        }
    }

    async getMonthlyStats(month) {
        try {
            return await Paste.getMonthlyStats(month);
        } catch (error) {
            console.error('Get monthly stats error:', error);
            throw new Error('Failed to fetch monthly statistics');
        }
    }

    async getLast5MonthsStats() {
        try {
            return await Paste.getLast5MonthsStats();
        } catch (error) {
            console.error('Get last 5 months stats error:', error);
            throw new Error('Failed to fetch statistics for last 5 months');
        }
    }
}

module.exports = new PasteService(); 