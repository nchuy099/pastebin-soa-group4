import { Link } from 'react-router-dom';

export default function Header() {
    return (
        <header style={{
            padding: '1rem',
            background: '#282c34',
            display: 'flex',
            alignItems: 'center'
        }}>
            <Link to="/create">
                <img
                    src="/Pastebin.png"   /* or your own /public/logo.png */
                    alt="Logo"
                    style={{ height: '40px' }}
                />
            </Link>
        </header>
    );
}
