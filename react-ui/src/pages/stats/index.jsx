import React, { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import { Line } from 'react-chartjs-2';
import { Chart as ChartJS, CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend } from 'chart.js';
import TwoColumnLayout from '../../components/TwoColumnLayout';

// Register Chart.js components
ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend);

export default function PasteStats({ setHeaderTitle }) {
    const { id } = useParams();
    const [stats, setStats] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [mode, setMode] = useState('last-7-days');

    // Fetch stats based on mode
    useEffect(() => {
        console.log(`Fetching stats for paste ${id}, mode: ${mode}`);
        setHeaderTitle?.('Paste Statistics');
        setLoading(true);
        fetch(
            `/stats/api/paste/${id}/stats?mode=${mode}`
        )
            .then(res => {
                console.log('Response status:', res.status, 'URL:', `/stats/api/paste/${id}/stats?mode=${mode}`);
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
                console.log('API response:', result);
                if (result?.data) {
                    setStats(result.data);
                    setHeaderTitle?.(`Stats for Paste ${id}`);
                } else {
                    throw new Error('Stats not found');
                }
            })
            .catch(e => {
                console.error('Fetch error:', e.message);
                setError(`Failed to load stats: ${e.message}`);
            })
            .finally(() => setLoading(false));
    }, [id, mode, setHeaderTitle]);

    // Handle mode toggle
    const handleModeChange = (newMode) => {
        setMode(newMode);
    };

    // Chart data and options
    const chartData = stats ? {
        labels: stats.timeViews.map(view => view.time),
        datasets: [
            {
                label: 'Views',
                data: stats.timeViews.map(view => view.views),
                borderColor: 'rgb(75, 192, 192)',
                backgroundColor: 'rgba(75, 192, 192, 0.2)',
                tension: 0.1,
                fill: false,
            },
        ],
    } : { labels: [], datasets: [] };

    const chartOptions = {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
            legend: {
                position: 'top',
            },
            title: {
                display: true,
                text: `View Trends (${mode.replace(/-/g, ' ')})`,
            },
        },
        scales: {
            x: {
                title: {
                    display: true,
                    text: mode.includes('days') ? 'Date (MM/DD)' : 'Time (HH:MM)',
                },
            },
            y: {
                beginAtZero: true,
                title: {
                    display: true,
                    text: 'Number of Views',
                },
            },
        },
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
                    {error}
                </div>
            </div>
        </TwoColumnLayout>
    );

    if (!stats) return (
        <TwoColumnLayout>
            <div className="my-4">
                <div className="alert alert-warning" role="alert">
                    Stats not found
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
                            Paste Statistics
                            <small className="text-muted d-block mt-1 fs-6">ID: {stats.pasteId}</small>
                        </h2>

                        {/* General Info */}
                        <div className="mb-3">
                            <h5 className="mb-1">General Information</h5>
                            <dl className="mb-0">
                                <div className="d-flex">
                                    <dt className="fw-bold me-2">Total Views (All Time):</dt>
                                    <dd>{stats.totalViewsFromCreation}</dd>
                                </div>
                                <div className="d-flex">
                                    <dt className="fw-bold me-2">Total Views ({mode.replace(/-/g, ' ')}):</dt>
                                    <dd>{stats.totalViews}</dd>
                                </div>
                                <div className="d-flex">
                                    <dt className="fw-bold me-2">Timezone:</dt>
                                    <dd>{stats.timezone}</dd>
                                </div>
                            </dl>
                        </div>

                        {/* Mode Toggles */}
                        <div className="mb-3">
                            <h5 className="mb-2">Select Time Range</h5>
                            <div className="btn-group" role="group">
                                {['last-10-minutes', 'last-24-hours', 'last-7-days', 'last-30-days'].map(m => (
                                    <button
                                        key={m}
                                        type="button"
                                        className={`btn ${mode === m ? 'btn-primary' : 'btn-outline-primary'} btn-sm`}
                                        onClick={() => handleModeChange(m)}
                                    >
                                        {m.replace(/-/g, ' ')}
                                    </button>
                                ))}
                            </div>
                        </div>

                        {/* Chart */}
                        <div className="mb-3">
                            <h5 className="mb-2">View Trends</h5>
                            <div style={{ height: '400px', position: 'relative' }}>
                                <Line data={chartData} options={chartOptions} />
                            </div>
                        </div>

                        {/* Back to Paste Link */}
                        <div className="d-flex gap-2 align-items-center">
                            <Link to={`/get-paste/${id}`} className="btn btn-secondary btn-sm">
                                Back to Paste
                            </Link>
                        </div>
                    </div>
                </div>
            </div>
        </TwoColumnLayout>
    );
}