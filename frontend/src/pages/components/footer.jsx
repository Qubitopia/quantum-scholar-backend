import React from 'react';
import { Container, Row, Col } from 'react-bootstrap';
const styles = {
    footer: {
        backgroundColor: '#2d3748',
        color: '#f8f9fa',
        padding: '1.5rem 0',
        marginTop: 'auto',
    },
    link: {
        color: '#f8f9fa',
        marginRight: '1rem',
        textDecoration: 'none',
    },
    linkLast: {
        color: '#f8f9fa',
        textDecoration: 'none',
    }
};
const Footer = () => (
    <footer style={styles.footer}>
        <Container>
            <Row>
                <Col md={6}>
                    <h5>Mock Project</h5>
                    <p>&copy; {new Date().getFullYear()} All rights reserved.</p>
                </Col>
                <Col md={6} className="text-md-end">
                    <a href="/privacy" className="text-light me-3">Privacy Policy</a>
                    <a href="/terms" className="text-light">Terms of Service</a>
                </Col>
            </Row>
        </Container>
    </footer>
);

export default Footer;
