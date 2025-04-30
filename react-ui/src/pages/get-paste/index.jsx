import React, { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';

export default function GetPaste() {
    const { id } = useParams();
    const [paste, setPaste]     = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError]     = useState('');

    useEffect(() => {
        fetch(`/get-paste/api/paste/${id}`)
            .then(res => {
                if (!res.ok) throw new Error(`Server ${res.status}`);
                return res.json();
            })
            .then(result => setPaste(result.data))
            .catch(e => setError(e.message))
            .finally(() => setLoading(false));
    }, [id]);

    const copyLink = () => {
        const url = `${window.location.origin}/get-paste/${id}`;
        navigator.clipboard.writeText(url);
        alert('Link copied to clipboard!');
    };

    if (loading) return (
        <div className="d-flex justify-content-center my-5">
            <div className="spinner-border text-primary" role="status">
                <span className="visually-hidden">Loading...</span>
            </div>
        </div>
    );

    if (error) return (
        <div className="container my-5">
            <div className="alert alert-danger" role="alert">
                Error: {error}
            </div>
        </div>
    );

    if (!paste) return (
        <div className="container my-5">
            <div className="alert alert-warning" role="alert">
                Paste not found
            </div>
        </div>
    );

    return (
        <div className="container-md my-5">
            <div className="card shadow">
                <div className="card-body p-4">
                    <h1 className="card-title mb-4">
                        {paste.title || 'Untitled Paste'}
                        <small className="text-muted d-block mt-2 fs-5">ID: {paste.id}</small>
                    </h1>

                    <div className="row g-3 mb-4">
                        <div className="col-md-6">
                            <dl className="row">
                                <dt className="col-sm-4">Language</dt>
                                <dd className="col-sm-8">{paste.language}</dd>

                                <dt className="col-sm-4">Created</dt>
                                <dd className="col-sm-8">
                                    {new Date(paste.created_at).toLocaleString()}
                                </dd>

                                {paste.expires_at && (
                                    <>
                                        <dt className="col-sm-4">Expires</dt>
                                        <dd className="col-sm-8">
                                            {new Date(paste.expires_at).toLocaleString()}
                                        </dd>
                                    </>
                                )}
                            </dl>
                        </div>
                    </div>

                    <div className="mb-4">
                        <h5 className="mb-3">Content</h5>
                        <pre className="p-3 bg-light rounded font-monospace border"
                             style={{ whiteSpace: 'pre-wrap', wordBreak: 'break-word' }}>
                            {paste.content}
                        </pre>
                    </div>

                    <div className="d-flex gap-2">
                        <button
                            onClick={copyLink}
                            className="btn btn-primary"
                        >
                            Copy Link
                        </button>

                        <Link
                            to={`/stats/${id}`}
                            className="btn btn-info text-white"
                        >
                            View Stats
                        </Link>
                    </div>
                </div>
            </div>
        </div>
    );
}
