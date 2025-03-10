const mongoose = require('mongoose');

const pasteSchema = new mongoose.Schema({
    title: {
        type: String,
        required: false,
        trim: true,
        default: 'Untitled'
    },
    content: {
        type: String,
        required: [true, 'Content is required']
    },
    language: {
        type: String,
        default: 'text'
    },
    visibility: {
        type: String,
        enum: ['public', 'private', 'unlisted'],
        default: 'public'
    },
    expiresAt: {
        type: Date,
        default: null
    },
    author: {
        type: mongoose.Schema.Types.ObjectId,
        ref: 'User',
        required: false
    },
    views: {
        type: Number,
        default: 0
    }
}, {
    timestamps: true
});

// Index for efficient querying
pasteSchema.index({ visibility: 1, createdAt: -1 });
pasteSchema.index({ author: 1, createdAt: -1 });

module.exports = mongoose.model('Paste', pasteSchema); 