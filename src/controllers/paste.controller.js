const pasteService = require('../services/paste.service');

const createPaste = async (req, res) => {
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
};

const getPaste = async (req, res) => {
    try {
        const result = await pasteService.getPasteById(req.params.id);

        if (result.status === 'not_found') {
            return res.status(404).render('index', {
                error: 'Paste not found',
                pasteId: req.params.id
            });
        }

        if (result.status === 'expired') {
            return res.status(403).render('expired', {
                error: 'This paste has expired and is no longer accessible',
                pasteId: req.params.id
            });
        }

        res.render('paste', { paste: result.paste });
    } catch (error) {
        console.error('Error retrieving paste:', error);
        res.status(500).render('index', {
            error: 'Failed to retrieve paste',
            pasteId: req.params.id
        });
    }
};

const getPublicPastes = async (req, res) => {
    try {
        const pastes = await pasteService.getPublicPastes();
        res.render('public', { pastes });
    } catch (error) {
        console.error('Error fetching public pastes:', error);
        res.render('public', { pastes: [], error: 'Failed to fetch public pastes' });
    }
};

const showCreateForm = (req, res) => {
    res.render('index', { error: null });
};

const getMonthlyStats = async (req, res) => {
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
};

module.exports = {
    createPaste,
    getPaste,
    getPublicPastes,
    showCreateForm,
    getMonthlyStats
}; 