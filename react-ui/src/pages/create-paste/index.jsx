import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Form, Button, Alert, Row, Col } from 'react-bootstrap';
import TwoColumnLayout from '../../components/TwoColumnLayout';

export default function CreatePaste({ setHeaderTitle }) {
    const navigate = useNavigate();

    const [content, setContent] = useState('');
    const [title, setTitle] = useState('');
    const [language, setLanguage] = useState('text');
    const [expiresIn, setExpiresIn] = useState('60');
    const [visibility, setVisibility] = useState('PUBLIC');
    const [error, setError] = useState('');

    const languageOptions = [
        { label: 'Text', value: 'text' },
        { label: 'C', value: 'c' },
        { label: 'C++', value: 'cpp' },
        { label: 'C#', value: 'csharp' },
        { label: 'CSS', value: 'css' },
        { label: 'HTML', value: 'html' },
        { label: 'Java', value: 'java' },
        { label: 'JavaScript', value: 'javascript' },
        { label: 'PHP', value: 'php' },
        { label: 'Python', value: 'python' },
        { label: 'Ruby', value: 'ruby' },
        { label: 'SQL', value: 'sql' },
        { label: 'TypeScript', value: 'typescript' },
        { label: 'XML', value: 'xml' },
    ];

    const expiryOptions = [
        { label: 'No Expiry', value: '' },
        { label: '5 Minutes', value: '300' },
        { label: '10 Minutes', value: '600' },
        { label: '30 Minutes', value: '1800' },
        { label: '1 Hour', value: '3600' },
        { label: '2 Hours', value: '7200' },
        { label: '6 Hours', value: '21600' },
        { label: '12 Hours', value: '43200' },
        { label: '1 Day', value: '86400' },
        { label: '2 Days', value: '172800' },
        { label: '1 Week', value: '604800' },
    ];

    useEffect(() => {
        setHeaderTitle('Create Paste');
    }, [setHeaderTitle]);

    const handleSubmit = async () => {
        setError('');

        const payload = {
            content,
            title,
            language,
            expiresIn: expiresIn === '' ? 0 : parseInt(expiresIn, 10),
            visibility,
        };

        try {
            const res = await fetch('/create-paste/api/paste', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload),
            });
            const result = await res.json();

            if (!res.ok) {
                throw new Error(result.message || `Server error: ${res.status}`);
            }
            if (!result.data?.id) {
                throw new Error('Missing paste ID in response');
            }

            navigate(`/get-paste/${result.data.id}`);
        } catch (e) {
            setError(e.message);
        }
    };

    return (
        <TwoColumnLayout>
            <Form>
                <Row>
                    <Form.Group as={Col} xs={12} className="mb-3">
                        <Form.Label>Title</Form.Label>
                        <Form.Control
                            type="text"
                            placeholder="Untitled"
                            value={title}
                            onChange={e => setTitle(e.target.value)}
                        />
                    </Form.Group>
                </Row>
                <Row className="mb-3">
                    <Form.Group as={Col} md={4}>
                        <Form.Label>Language</Form.Label>
                        <Form.Select
                            value={language}
                            onChange={e => setLanguage(e.target.value)}
                        >
                            {languageOptions.map(opt => (
                                <option key={opt.value} value={opt.value}>
                                    {opt.label}
                                </option>
                            ))}
                        </Form.Select>
                    </Form.Group>
                    <Form.Group as={Col} md={4}>
                        <Form.Label>Visibility</Form.Label>
                        <Form.Select
                            value={visibility}
                            onChange={e => setVisibility(e.target.value)}
                        >
                            <option value="PUBLIC">Public</option>
                            <option value="UNLISTED">Unlisted</option>
                        </Form.Select>
                    </Form.Group>
                    <Form.Group as={Col} md={4}>
                        <Form.Label>Expires In</Form.Label>
                        <Form.Select
                            value={expiresIn}
                            onChange={e => setExpiresIn(e.target.value)}
                        >
                            {expiryOptions.map(opt => (
                                <option key={opt.value} value={opt.value}>
                                    {opt.label}
                                </option>
                            ))}
                        </Form.Select>
                    </Form.Group>
                </Row>
                <Row>
                    <Form.Group as={Col} xs={12} className="mb-3">
                        <Form.Label>Content</Form.Label>
                        <Form.Control
                            as="textarea"
                            rows={10}
                            className="font-monospace"
                            placeholder="Type or paste your content here..."
                            value={content}
                            onChange={e => setContent(e.target.value)}
                        />
                    </Form.Group>
                </Row>
                {error && (
                    <Row>
                        <Col xs={12} className="mb-3">
                            <Alert variant="danger">{error}</Alert>
                        </Col>
                    </Row>
                )}
                <Row>
                    <Col xs={12}>
                        <Button
                            variant="primary"
                            onClick={handleSubmit}
                            className="w-100 py-2"
                        >
                            Create Paste
                        </Button>
                    </Col>
                </Row>
            </Form>
        </TwoColumnLayout>
    );
}