// react-ui/src/components/Header.jsx
import { Link } from 'react-router-dom';

export default function Header({ title }) {
    return (
        <header style={{
            padding: '1rem',
            background: '#282c34',
            display: 'flex',
            alignItems: 'center',
            color: 'white'
        }}>
            <Link to="/create">
                <img
                    src="/Pastebin.png"
                    alt="Logo"
                    style={{ height: '40px', marginRight: '1rem' }}
                />
            </Link>
            {title && <h2 style={{ margin: 0, fontSize: '1.25rem' }}>{title}</h2>}
        </header>
    );
}
