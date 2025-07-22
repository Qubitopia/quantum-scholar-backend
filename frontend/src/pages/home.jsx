import React from 'react';
import Navbar from './components/navbar.jsx';
import Footer from './components/footer.jsx';
const styles = {
    page: {
        minHeight: '100vh',
        display: 'flex',
        flexDirection: 'column',
    },
    container: {
        flex: 1,
        minHeight: '70vh',
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        padding: '2rem',
        background: '#f7fafc',
    },
    heading: {
        fontSize: '2.5rem',
        fontWeight: 'bold',
        marginBottom: '1rem',
        color: '#2d3748',
    },
    paragraph: {
        fontSize: '1.2rem',
        color: '#4a5568',
        marginBottom: '2rem',
        maxWidth: '600px',
        textAlign: 'center',
    },
    buttonGroup: {
        display: 'flex',
        gap: '1rem',
        marginBottom: '2rem',
    },
    button: {
        padding: '0.75rem 2rem',
        fontSize: '1rem',
        borderRadius: '5px',
        border: 'none',
        cursor: 'pointer',
        fontWeight: 'bold',
        background: '#3182ce',
        color: '#fff',
        transition: 'background 0.2s',
    },
    buttonAlt: {
        background: '#38a169',
    },
    graphicPlaceholder: {
        width: '300px',
        height: '180px',
        background: '#e2e8f0',
        borderRadius: '10px',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        color: '#718096',
        fontSize: '1.1rem',
        marginBottom: '2rem',
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

const Home = () => {
    return (
        <div style={styles.page}>
            <Navbar />
            <div style={styles.container}>
                <div style={styles.graphicPlaceholder}>
                    {/* Add your graphic/image here */}
                    Graphic/Image Placeholder
                </div>
                <h1 style={styles.heading}>Welcome to Mock Project</h1>
                <p style={styles.paragraph}>
                    Mock Project is a platform designed for students and teachers to collaborate, learn, and grow together. 
                    Easily manage your courses, assignments, and communication in one place.
                </p>
                <div style={styles.buttonGroup}>
                    <a href="/login" style={{ ...styles.button }}>Login</a>
                    <a href="/signup/student" style={{ ...styles.button, ...styles.buttonAlt }}>Sign Up as Student</a>
                    <a href="/signup/teacher" style={{ ...styles.button, background: '#805ad5' }}>Sign Up as Teacher</a>
                </div>
            </div>
            <Footer />

        </div>
    );
};

export default Home;