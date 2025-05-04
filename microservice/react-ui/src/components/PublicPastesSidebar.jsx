import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';

export default function PublicPastesSidebar({ limit = 5 }) {
    const [pastes, setPastes]   = useState([]);
    const [error, setError]     = useState('');

    useEffect(() => {
        (async () => {
            try {
                const res  = await fetch(`/public/api/paste?limit=${limit}&page=1`);
                const json = await res.json();
                if (!res.ok) throw new Error(json.message || res.status);
                setPastes(json.data.pastes || []);
            } catch (e) {
                setError(e.message);
            }
        })();
    }, [limit]);

    if (error) {
        return <div className="alert alert-danger py-1">{error}</div>;
    }

    if (!pastes.length) {
        return <div className="text-muted">No public pastes</div>;
    }

    return (
        <ul className="list-unstyled">
            {pastes.map(p => (
                <li key={p.id} className="mb-3">
                    <Link to={`/get-paste/${p.id}`} className="text-decoration-none">
                        <strong>{p.title || p.id}</strong>
                    </Link>
                    <br/>
                    <small className="text-secondary">{p.language}</small>
                </li>
            ))}
        </ul>
    );
}
