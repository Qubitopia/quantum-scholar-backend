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
            <div style={styles.navBrand}>Mock Project</div>
            <div style={styles.navLinks}>
                <Link to="/" style={{ color: '#fff', textDecoration: 'none' }}>Home</Link>
                <Link to="/about" style={{ color: '#fff', textDecoration: 'none' }}>About</Link>
                <Link to="/classroom" style={{ color: '#fff', textDecoration: 'none' }}>Classroom</Link>

            </div>
        </nav>
    )
}

export default navbar;