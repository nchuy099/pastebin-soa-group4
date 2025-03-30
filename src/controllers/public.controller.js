const axios = require('axios');
const model = require('../models/public.model');
const logger = require('../utils/logger');

async function handlePublicRequest(req, res) {
    try {
        // Try microservice first
        const pastes = await model.getPublicPastes();
        return res.render('public', { pastes });
    } catch (error) {
        logger.error(`Microservice failed: ${error.message}`);

        // Fallback to monolith
        try {
            const response = await axios.get(
                `${process.env.MONOLITH_URL}/public`,
                { timeout: 3000 }
            );
            return res.render('public_fallback', { pastes: response.data });
        } catch (fallbackError) {
            logger.error(`Fallback failed: ${fallbackError.message}`);
            return res.status(503).render('error');
        }
    }
}

module.exports = { handlePublicRequest };