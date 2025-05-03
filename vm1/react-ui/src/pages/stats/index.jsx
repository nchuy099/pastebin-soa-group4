import React from 'react';
import { useParams, Link } from 'react-router-dom';

export default function Stats() {
    const { id } = useParams();

    return (
        <div>
            <h1>Stats for Paste {id}</h1>
            <p>(Coming soon!)</p>
            <Link to={`/get-paste/${id}`}>
                <button>Back to Paste</button>
            </Link>
        </div>
    );
}
