import React from 'react';
import { Container, Row, Col } from 'react-bootstrap';
import PublicPastesSidebar from './PublicPastesSidebar';
import './TwoColumnLayout.css';

export default function TwoColumnLayout({ children }) {
    return (
        <Container fluid className="py-3" style={{ maxWidth: '1600px' }}>
            <Row>
                <Col xs={12} md={2} className="mb-3 sidebar-col">
                    <h5 className="mb-2">Public Pastes</h5>
                    <PublicPastesSidebar limit={5} />
                </Col>
                <Col xs={12} md={10} className="main-col">
                    {children}
                </Col>
            </Row>
        </Container>
    );
}