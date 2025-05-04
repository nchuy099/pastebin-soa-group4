// react-ui/src/App.jsx
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { useState } from 'react';
import CreatePaste from './pages/create-paste';
import GetPaste from './pages/get-paste';
import Stats from './pages/stats';
import Header from './components/Header';

export default function App() {
    const [headerTitle, setHeaderTitle] = useState('Create Paste');

    return (
        <Router>
            <Header title={headerTitle} />
            <div style={{ maxWidth: 800, margin: '2rem auto' }}>
                <Routes>
                    <Route path="/create" element={<CreatePaste setHeaderTitle={setHeaderTitle} />} />
                    <Route path="/get-paste/:id" element={<GetPaste setHeaderTitle={setHeaderTitle} />} />
                    <Route path="/stats/:id" element={<Stats />} />
                    <Route path="*" element={<Navigate to="/create" replace />} />
                </Routes>
            </div>
        </Router>
    );
}
