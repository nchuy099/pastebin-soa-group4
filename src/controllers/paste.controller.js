const pasteService = require('../services/paste.service');

class PasteController {
    async createPaste(req, res) {
        try {
            const { content, title, language, expires_in: expiresIn, visibility } = req.body;

            if (!content) {
                return res.render('index', { error: 'Content is required' });
            }

            const result = await pasteService.createPaste({
                content,
                title,
                language,
                expiresIn,
                visibility
            });

            res.redirect(`/paste/${result.id}`);
        } catch (error) {
            console.error('Error creating paste:', error);
            res.render('index', { error: 'Failed to create paste' });
        }
    }

    async getPaste(req, res) {
        try {
            const paste = await pasteService.getPasteById(req.params.id);
            res.render('paste', { paste });
        } catch (error) {
            console.error('Error retrieving paste:', error);
            const status = error.message.includes('expired') ? 403 :
                error.message.includes('not found') ? 404 : 500;
            const template = status === 403 ? 'expired' : 'index';
            res.status(status).render(template, {
                error: error.message,
                pasteId: req.params.id
            });
        }
    }

    async getPublicPastes(req, res) {
        try {
            const pastes = await pasteService.getPublicPastes();
            res.render('public', { pastes });
        } catch (error) {
            console.error('Error fetching public pastes:', error);
            res.render('public', { pastes: [], error: 'Failed to fetch public pastes' });
        }
    }

    showCreateForm(req, res) {
        res.render('index', { error: null });
    }

    async getMonthlyStats(req, res) {
        try {
            let month = req.params.month;

            // If no month provided, use current month
            if (!month) {
                const now = new Date();
                const year = now.getFullYear();
                const currentMonth = String(now.getMonth() + 1).padStart(2, '0');
                month = `${year}-${currentMonth}`;
            }

            // Validate month format (YYYY-MM)
            if (!/^\d{4}-\d{2}$/.test(month)) {
                throw new Error('Invalid month format. Use YYYY-MM');
            }

            const stats = await pasteService.getMonthlyStats(month);
            res.render('stats', { stats, error: null });
        } catch (error) {
            console.error('Error fetching monthly stats:', error);
            res.render('stats', {
                stats: null,
                error: error.message || 'Failed to fetch statistics'
            });
        }
    }
}

module.exports = new PasteController(); 