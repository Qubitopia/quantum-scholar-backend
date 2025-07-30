import { Link } from "react-router-dom";

const styles = {
    navbar: {
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center',
        padding: '1rem 2rem',
        background: '#2d3748',
        color: '#fff',
    },
    navBrand: {
        fontWeight: 'bold',
        fontSize: '1.5rem',
    },
    navLinks: {
        display: 'flex',
        gap: '1rem',
    },
};

function navbar() {
    return (
        <nav style={styles.navbar}>
            <div style={{ ...styles.navBrand, display: 'flex', alignItems: 'center', height: '100%' }}>Mock Project</div>
            <div style={{ ...styles.navLinks, alignItems: 'center', display: 'flex', height: '100%' }}>
                <Link to="/" style={{ color: '#fff', textDecoration: 'none', display: 'flex', alignItems: 'center', height: '100%' }}>Home</Link>
                <Link to="/about" style={{ color: '#fff', textDecoration: 'none', display: 'flex', alignItems: 'center', height: '100%' }}>About</Link>
                <Link to="/classroom" style={{ color: '#fff', textDecoration: 'none', display: 'flex', alignItems: 'center', height: '100%' }}>Classroom</Link>
                <Link to="/settings" style={{ color: '#fff', textDecoration: 'none', fontSize: '1.5rem', display: 'flex', alignItems: 'center', height: '100%' }}>⚙️</Link>
            </div>
        </nav>
    )
}

export default navbar;