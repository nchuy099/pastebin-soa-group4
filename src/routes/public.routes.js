const express = require('express');
const router = express.Router();
const controller = require('../controllers/public.controller');

// Main endpoint
router.get('/public', controller.handlePublicRequest);

// Health check endpoint
router.get('/health', (req, res) => res.json({
    status: 'ok',
    version: process.env.APP_VERSION
}));

module.exports = router;