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

        res.status(302).redirect(`/paste/${result.id}`);
    } catch (error) {
        console.error('Error creating paste:', error.message);
        res.status(500).render('index', { error: error.message });
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
            
        res.status(200).render('paste', { paste: result.paste, error: null });
    } catch (error) {
        console.error(`Error retrieving paste ${req.params.id}:`, error.message);
        res.status(500).render('index', {
            error: error.message,
            pasteId: req.params.id
        });
    }
};

const getPublicPastes = async (req, res) => {
    try {
        // Lấy tham số page từ query string
        const page = parseInt(req.query.page) || 1;
        const result = await pasteService.getPublicPastes(page);
        res.status(200).render('public', { pastes: result.pastes, pagination: result.pagination, error: null });
    } catch (error) {
        console.error(`Error fetching public page ${req.query.page}:`, error.message);
        res.status(500).render('public', { pastes: [], error: error.message, page: req.query.page });
    }
};

const showCreateForm = (req, res) => {
    try {
        res.status(200).render('index', { error: null });
    } catch (error) {
        console.error('Error showing create form:', error.message);
        res.status(500).render('index', { error: error.message });
    }
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
        res.status(200).render('stats', { stats, error: null });
    } catch (error) {
        console.error(`Error fetching stats on ${req.params.month}:`, error.message);
        res.status(500).render('stats', {
            stats: null,
            error: error.message || 'Failed to fetch statistics',
            month: req.params.month
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