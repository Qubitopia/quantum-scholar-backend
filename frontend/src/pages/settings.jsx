import React from 'react';
import Navbar from './components/navbar.jsx';
const Settings = () => {
    return (
        <div style={{
            minHeight: '100vh',
            width: '100vw',
            background: '#f5f6fa',
            display: 'flex',
            flexDirection: 'column'
        }}>
            <Navbar />
            <div style={{
                flex: 1,
                display: 'flex',
                justifyContent: 'center',
                alignItems: 'center'
            }}>
                <div style={{
                    width: '100%',
                    maxWidth: 600,
                    background: '#fff',
                    borderRadius: 8,
                    boxShadow: '0 2px 8px rgba(0,0,0,0.07)',
                    padding: 32
                }}>
                    <h1>Settings</h1>
                    <form>
                        <div style={{ marginBottom: 16 }}>
                            <label>
                                <strong>Username:</strong>
                                <input type="text" name="username" style={{ marginLeft: 10 }} />
                            </label>
                        </div>
                        <div style={{ marginBottom: 16 }}>
                            <label>
                                <strong>Email:</strong>
                                <input type="email" name="email" style={{ marginLeft: 10 }} />
                            </label>
                        </div>
                        <div style={{ marginBottom: 16 }}>
                            <label>
                                <strong>Change Password:</strong>
                                <input type="password" name="password" style={{ marginLeft: 10 }} />
                            </label>
                        </div>
                        <button type="submit">Save Changes</button>
                    </form>
                </div>
            </div>
        </div>
    );
};

export default Settings;