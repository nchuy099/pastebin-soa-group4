// react-ui/src/App.jsx
import {BrowserRouter as Router, Routes, Route, Link, Navigate} from 'react-router-dom';
import CreatePaste from './pages/create-paste';
import GetPaste    from './pages/get-paste';
import Stats       from './pages/stats';
import Header from "./components/header";

export default function App() {
    return (
        <Router>
            <Header />
            <div style={{ maxWidth: 800, margin: '2rem auto' }}>
                <Routes>
                    <Route path="/create" element={<CreatePaste />} />
                    <Route path="/get-paste/:id" element={<GetPaste />} />
                    <Route path="/stats/:id"    element={<Stats />} />
                    <Route path="*"              element={<Navigate to="/create" replace />} />
                </Routes>
            </div>
        </Router>
    );
}
