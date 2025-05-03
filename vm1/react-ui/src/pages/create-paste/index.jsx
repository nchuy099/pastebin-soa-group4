// react-ui/src/pages/CreatePaste/index.jsx
import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

export default function CreatePaste() {
    const navigate = useNavigate();

    const [content, setContent]       = useState('');
    const [title, setTitle]           = useState('');
    const [language, setLanguage]     = useState('text');
    const [expiresIn, setExpiresIn]   = useState('60');
    const [visibility, setVisibility] = useState('PUBLIC');
    // const [pasteLink, setPasteLink]   = useState('');
    const [error, setError]           = useState('');

    const languageOptions = [
        { label: 'Text', value: 'text' },
        { label: 'C', value: 'c' },
        { label: 'C++', value: 'cpp' },
        { label: 'C#', value: 'csharp' },
        { label: 'CSS', value: 'css' },
        { label: 'HTML', value: 'html' },
        { label: 'Java', value: 'java' },
        { label: 'JavaScript', value: 'javascript' },
        { label: 'PHP', value: 'php' },
        { label: 'Python', value: 'python' },
        { label: 'Ruby', value: 'ruby' },
        { label: 'SQL', value: 'sql' },
        { label: 'TypeScript', value: 'typescript' },
        { label: 'XML', value: 'xml' },
    ];

    const expiryOptions = [
        { label: 'No Expiry', value: '' },
        { label: '5 Minutes', value: '300' },
        { label: '10 Minutes', value: '600' },
        { label: '30 Minutes', value: '1800' },
        { label: '1 Hour', value: '3600' },
        { label: '2 Hours', value: '7200' },
        { label: '6 Hours', value: '21600' },
        { label: '12 Hours', value: '43200' },
        { label: '1 Day', value: '86400' },
        { label: '2 Days', value: '172800' },
        { label: '1 Week', value: '604800' },
    ];

    const handleSubmit = async () => {
        setError('');
        // setPasteLink('');

        const payload = {
            content,
            title,
            language,
            expiresIn: expiresIn === '' ? 0 : parseInt(expiresIn, 10),
            visibility,
        };

        try {
            const res = await fetch('/create-paste/api/paste', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload),
            });
            const result = await res.json();

            if (!res.ok) {
                throw new Error(result.message || `Server error: ${res.status}`);
            }
            if (!result.data?.id) {
                throw new Error('Missing paste ID in response');
            }

            navigate(`/get-paste/${result.data.id}`);
        } catch (e) {
            setError(e.message);
        }
    };

    return (
        <div className="container-md my-5">
            <div className="card shadow">
                <div className="card-body p-4">
                    <h1 className="card-title text-center mb-4">Create a New Paste</h1>

                    <div className="row g-4">
                        {/* Title */}
                        <div className="col-12">
                            <label className="form-label">Title</label>
                            <input
                                type="text"
                                className="form-control"
                                placeholder="Untitled"
                                value={title}
                                onChange={e => setTitle(e.target.value)}
                            />
                        </div>

                        {/* Language, Visibility, Expiry */}
                        <div className="col-md-4">
                            <label className="form-label">Language</label>
                            <select
                                className="form-select"
                                value={language}
                                onChange={e => setLanguage(e.target.value)}
                            >
                                {languageOptions.map(opt => (
                                    <option key={opt.value} value={opt.value}>
                                        {opt.label}
                                    </option>
                                ))}
                            </select>
                        </div>

                        <div className="col-md-4">
                            <label className="form-label">Visibility</label>
                            <select
                                className="form-select"
                                value={visibility}
                                onChange={e => setVisibility(e.target.value)}
                            >
                                <option value="PUBLIC">Public</option>
                                <option value="UNLISTED">Unlisted</option>
                            </select>
                        </div>

                        <div className="col-md-4">
                            <label className="form-label">Expires In</label>
                            <select
                                className="form-select"
                                value={expiresIn}
                                onChange={e => setExpiresIn(e.target.value)}
                            >
                                {expiryOptions.map(opt => (
                                    <option key={opt.value} value={opt.value}>
                                        {opt.label}
                                    </option>
                                ))}
                            </select>
                        </div>

                        {/* Content */}
                        <div className="col-12">
                            <label className="form-label">Content</label>
                            <textarea
                                className="form-control font-monospace"
                                value={content}
                                onChange={e => setContent(e.target.value)}
                                placeholder="Type or paste your content here..."
                                rows="8"
                            />
                        </div>

                        {/* Error Message */}
                        {error && (
                            <div className="col-12">
                                <div className="alert alert-danger">{error}</div>
                            </div>
                        )}

                        {/* Submit Button */}
                        <div className="col-12">
                            <button
                                onClick={handleSubmit}
                                className="btn btn-primary w-100 py-2"
                            >
                                Create Paste
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}
