import React from 'react';
import Navbar from './components/navbar.jsx';
import Footer from './components/footer.jsx';

const styles = {
    page: {
        minHeight: '100vh',
        display: 'flex',
        flexDirection: 'column',
        background: 'linear-gradient(135deg, #f8fafc 0%, #e9ecef 100%)',
        fontFamily: 'Inter, Segoe UI, Arial, sans-serif',
    },
    container: {
        flex: 1,
        minHeight: '70vh',
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        padding: '3rem 1rem',
        background: 'rgba(255,255,255,0.85)',
        boxShadow: '0 8px 32px 0 rgba(31, 38, 135, 0.12)',
        borderRadius: '24px',
        margin: '2rem auto',
        maxWidth: '600px',
        backdropFilter: 'blur(2px)',
    },
    heading: {
        fontSize: '2.8rem',
        fontWeight: 800,
        marginBottom: '1.2rem',
        color: '#2d3748',
        letterSpacing: '-1px',
        textAlign: 'center',
    },
    paragraph: {
        fontSize: '1.18rem',
        color: '#4a5568',
        marginBottom: '2.2rem',
        maxWidth: '520px',
        textAlign: 'center',
        lineHeight: 1.7,
    },
    buttonGroup: {
        display: 'flex',
        gap: '1.2rem',
        marginBottom: '2.5rem',
        flexWrap: 'wrap',
        justifyContent: 'center',
    },
    button: {
        padding: '0.85rem 2.2rem',
        fontSize: '1.05rem',
        borderRadius: '8px',
        border: 'none',
        cursor: 'pointer',
        fontWeight: 600,
        background: 'linear-gradient(90deg, #3182ce 0%, #63b3ed 100%)',
        color: '#fff',
        boxShadow: '0 2px 8px rgba(49,130,206,0.08)',
        transition: 'transform 0.15s, box-shadow 0.15s, background 0.2s',
        textDecoration: 'none',
        outline: 'none',
    },
    buttonAlt: {
        background: 'linear-gradient(90deg, #38a169 0%, #68d391 100%)',
    },
    buttonTeacher: {
        background: 'linear-gradient(90deg, #805ad5 0%, #b794f4 100%)',
    },
    buttonHover: {
        transform: 'translateY(-2px) scale(1.03)',
        boxShadow: '0 4px 16px rgba(49,130,206,0.16)',
    },
    graphicPlaceholder: {
        width: '320px',
        height: '180px',
        background: 'linear-gradient(135deg, #e2e8f0 60%, #cbd5e1 100%)',
        borderRadius: '18px',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        color: '#718096',
        fontSize: '1.15rem',
        marginBottom: '2.2rem',
        boxShadow: '0 2px 12px rgba(160,174,192,0.10)',
        position: 'relative',
        overflow: 'hidden',
    },
    graphicSVG: {
        width: '100%',
        height: '100%',
        objectFit: 'cover',
        position: 'absolute',
        top: 0,
        left: 0,
        opacity: 0.7,
    },
    footer: {
        background: '#2d3748',
        color: '#fff',
        textAlign: 'center',
        padding: '1rem 0',
        fontSize: '1rem',
        marginTop: '2rem',
    },
};

const Button = ({ href, style, children }) => {
    const [hover, setHover] = React.useState(false);
    return (
        <a
            href={href}
            style={{
                ...styles.button,
                ...style,
                ...(hover ? styles.buttonHover : {}),
            }}
            onMouseEnter={() => setHover(true)}
            onMouseLeave={() => setHover(false)}
        >
            {children}
        </a>
    );
};

const Home = () => {
    return (
        <div style={styles.page}>
            <Navbar />
            <div style={styles.container}>
                <div style={styles.graphicPlaceholder}>
                    {/* Example SVG illustration */}
                    <svg style={styles.graphicSVG} viewBox="0 0 320 180" fill="none" xmlns="http://www.w3.org/2000/svg">
                        <ellipse cx="160" cy="90" rx="120" ry="70" fill="#90cdf4" opacity="0.3"/>
                        <rect x="60" y="60" width="200" height="60" rx="18" fill="#3182ce" opacity="0.15"/>
                        <circle cx="100" cy="100" r="24" fill="#38a169" opacity="0.18"/>
                        <circle cx="220" cy="80" r="18" fill="#805ad5" opacity="0.18"/>
                    </svg>
                    <span style={{position: 'relative', zIndex: 1}}>QuantumScholar</span>
                </div>
                <h1 style={styles.heading}>Welcome to QuantumScholar</h1>
                <p style={styles.paragraph}>
                    QuantumScholar is a modern platform for students and teachers to collaborate, learn, and grow together.<br />
                    Effortlessly manage your courses, assignments, and communication—all in one place.
                </p>
                <div style={styles.buttonGroup}>
                    <Button href="/login">Login</Button>
                    <Button href="/signup/student" style={styles.buttonAlt}>Sign Up as Student</Button>
                    <Button href="/signup/teacher" style={styles.buttonTeacher}>Sign Up as Teacher</Button>
                </div>
            </div>
            <Footer />
        </div>
    );
};

export default Home;