import React, { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import TwoColumnLayout from '../../components/TwoColumnLayout';

export default function GetPaste({ setHeaderTitle }) {
    const { id } = useParams();
    const [paste, setPaste] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [copySuccess, setCopySuccess] = useState(false);

    useEffect(() => {
        setHeaderTitle?.('View Paste');
        fetch(`/get-paste/api/paste/${id}`)
            .then(res => {
                if (!res.ok) {
                    throw new Error(`Server error: ${res.status}`);
                }
                const contentType = res.headers.get('content-type');
                if (!contentType?.includes('application/json')) {
                    throw new Error('Received non-JSON response from server');
                }
                return res.json();
            })
            .then(result => {
                if (result?.data) {
                    setPaste(result.data);
                    setHeaderTitle?.(result.data.title || 'View Paste');
                } else {
                    throw new Error('Paste not found');
                }
            })
            .catch(e => setError(e.message))
            .finally(() => setLoading(false));
    }, [id, setHeaderTitle]);

    const copyLink = () => {
        const url = `${window.location.origin}/get-paste/${id}`;
        navigator.clipboard.writeText(url)
            .then(() => {
                setCopySuccess(true);
                setTimeout(() => setCopySuccess(false), 3000);
            })
            .catch(err => {
                console.error('Failed to copy link:', err);
                setError('Failed to copy link');
            });
    };

    if (loading) return (
        <TwoColumnLayout>
            <div className="d-flex justify-content-center my-4">
                <div className="spinner-border text-primary" role="status">
                    <span className="visually-hidden">Loading...</span>
                </div>
            </div>
        </TwoColumnLayout>
    );

    if (error) return (
        <TwoColumnLayout>
            <div className="my-4">
                <div className="alert alert-danger" role="alert">
                    Error: {error}
                </div>
            </div>
        </TwoColumnLayout>
    );

    if (!paste) return (
        <TwoColumnLayout>
            <div className="my-4">
                <div className="alert alert-warning" role="alert">
                    Paste not found
                </div>
            </div>
        </TwoColumnLayout>
    );

    return (
        <TwoColumnLayout>
            <div className="my-3">
                <div className="card shadow h-100">
                    <div className="card-body p-3">
                        <h2 className="card-title mb-2">
                            {paste.title || 'Untitled Paste'}
                            <small className="text-muted d-block mt-1 fs-6">ID: {paste.id}</small>
                        </h2>

                        <div className="mb-2">
                            <h5 className="mb-1">Details</h5>
                            <dl className="mb-0">
                                <div className="d-flex">
                                    <dt className="fw-bold me-2">Language:</dt>
                                    <dd>{paste.language}</dd>
                                </div>
                                <div className="d-flex">
                                    <dt className="fw-bold me-2">Created:</dt>
                                    <dd>{new Date(paste.created_at).toLocaleString()}</dd>
                                </div>
                                {paste.expires_at && (
                                    <div className="d-flex">
                                        <dt className="fw-bold me-2">Expires:</dt>
                                        <dd>{new Date(paste.expires_at).toLocaleString()}</dd>
                                    </div>
                                )}
                            </dl>
                        </div>

                        <div className="mb-2 flex-grow-1">
                            <h5 className="mb-1">Content</h5>
                            <pre className="p-3 bg-light rounded font-monospace border"
                                 style={{
                                     whiteSpace: 'pre-wrap',
                                     wordBreak: 'break-word',
                                     minHeight: '300px',
                                     maxHeight: '70vh',
                                     overflowY: 'auto',
                                     width: '100%',
                                     boxSizing: 'border-box'
                                 }}>
                                {paste.content}
                            </pre>
                        </div>

                        <div className="d-flex gap-2 align-items-center flex-wrap">
                            <button
                                onClick={copyLink}
                                className="btn btn-primary btn-sm"
                            >
                                Copy Link
                            </button>
                            <Link
                                to={`/stats/${id}`}
                                className="btn btn-info text-white btn-sm"
                            >
                                View Stats
                            </Link>
                            {copySuccess && (
                                <span className="text-success ms-2">
                                    Link copied to clipboard!
                                </span>
                            )}
                        </div>
                    </div>
                </div>
            </div>
        </TwoColumnLayout>
    );
}