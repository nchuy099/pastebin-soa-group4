const express = require('express');
const router = express.Router();
const pasteController = require('../controllers/paste.controller');

// Show create form
router.get('/', pasteController.showCreateForm);

// Create a new paste
router.post('/paste', pasteController.createPaste);

// Get a paste by ID
router.get('/paste/:id', pasteController.getPaste);

// Get all public pastes
router.get('/public', pasteController.getPublicPastes);

// Get monthly statistics
router.get('/stats/:month?', pasteController.getMonthlyStats);

module.exports = router; 