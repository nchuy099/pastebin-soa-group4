const express = require('express');
const router = express.Router();
const { body, validationResult } = require('express-validator');
const { protect } = require('../middleware/auth.middleware');
const Paste = require('../models/paste.model');

// Create a new paste
router.post('/', [
    protect,
    body('content').notEmpty().withMessage('Content is required'),
    body('title').optional().trim(),
    body('language').optional().trim(),
    body('visibility').isIn(['public', 'private', 'unlisted']).optional(),
    body('expiresAt').optional().isISO8601().toDate()
], async (req, res) => {
    try {
        const errors = validationResult(req);
        if (!errors.isEmpty()) {
            return res.status(400).json({
                success: false,
                errors: errors.array()
            });
        }

        const paste = new Paste({
            ...req.body,
            author: req.user._id
        });

        await paste.save();

        res.status(201).json({
            success: true,
            data: paste
        });
    } catch (error) {
        res.status(500).json({
            success: false,
            message: 'Server error'
        });
    }
});

// Get all public pastes
router.get('/public', async (req, res) => {
    try {
        const pastes = await Paste.find({ visibility: 'public' })
            .sort('-createdAt')
            .populate('author', 'username')
            .limit(20);

        res.json({
            success: true,
            data: pastes
        });
    } catch (error) {
        res.status(500).json({
            success: false,
            message: 'Server error'
        });
    }
});

// Get a single paste by ID
router.get('/:id', async (req, res) => {
    try {
        const paste = await Paste.findById(req.params.id)
            .populate('author', 'username');

        if (!paste) {
            return res.status(404).json({
                success: false,
                message: 'Paste not found'
            });
        }

        // Check if paste is private
        if (paste.visibility === 'private') {
            if (!req.user || paste.author._id.toString() !== req.user._id.toString()) {
                return res.status(403).json({
                    success: false,
                    message: 'Not authorized to view this paste'
                });
            }
        }

        // Increment views
        paste.views += 1;
        await paste.save();

        res.json({
            success: true,
            data: paste
        });
    } catch (error) {
        res.status(500).json({
            success: false,
            message: 'Server error'
        });
    }
});

// Update a paste
router.put('/:id', [
    protect,
    body('content').optional().notEmpty(),
    body('title').optional().trim(),
    body('visibility').optional().isIn(['public', 'private', 'unlisted'])
], async (req, res) => {
    try {
        const paste = await Paste.findById(req.params.id);

        if (!paste) {
            return res.status(404).json({
                success: false,
                message: 'Paste not found'
            });
        }

        // Make sure user owns paste
        if (paste.author.toString() !== req.user._id.toString()) {
            return res.status(403).json({
                success: false,
                message: 'Not authorized to update this paste'
            });
        }

        const updatedPaste = await Paste.findByIdAndUpdate(
            req.params.id,
            req.body,
            { new: true, runValidators: true }
        );

        res.json({
            success: true,
            data: updatedPaste
        });
    } catch (error) {
        res.status(500).json({
            success: false,
            message: 'Server error'
        });
    }
});

// Delete a paste
router.delete('/:id', protect, async (req, res) => {
    try {
        const paste = await Paste.findById(req.params.id);

        if (!paste) {
            return res.status(404).json({
                success: false,
                message: 'Paste not found'
            });
        }

        // Make sure user owns paste
        if (paste.author.toString() !== req.user._id.toString()) {
            return res.status(403).json({
                success: false,
                message: 'Not authorized to delete this paste'
            });
        }

        await paste.remove();

        res.json({
            success: true,
            message: 'Paste removed'
        });
    } catch (error) {
        res.status(500).json({
            success: false,
            message: 'Server error'
        });
    }
});

// Get user's pastes
router.get('/user/me', protect, async (req, res) => {
    try {
        const pastes = await Paste.find({ author: req.user._id })
            .sort('-createdAt');

        res.json({
            success: true,
            data: pastes
        });
    } catch (error) {
        res.status(500).json({
            success: false,
            message: 'Server error'
        });
    }
});

module.exports = router; 