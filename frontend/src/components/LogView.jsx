import React, { useState, useEffect, useRef } from 'react';
import { GetLogs } from '../../wailsjs/go/main/App';

export const LogView = () => {
    const [logs, setLogs] = useState("");
    const logEndRef = useRef(null);

    const fetchLogs = () => {
        GetLogs(100)
            .then(log => {
                console.log("Received string length:", log.length);
                setLogs(log);
            })
            .catch(err => console.error("Wails Bridge Error:", err));
    };

    useEffect(() => {
        fetchLogs();
        const interval = setInterval(fetchLogs, 2000); // poll every 2s
        return () => clearInterval(interval);
    }, []);

    // uuto-scroll to bottom when logs update
    useEffect(() => {
        logEndRef.current?.scrollIntoView({ behavior: "smooth" });
    }, [logs]);

    return (
        <div className="terminal-container">
            <div className="terminal-header">
                <span>cura.log | System Stream</span>
            </div>
            <pre className="terminal-body">
                {logs || "Waiting for system events..."}
                <div ref={logEndRef} />
            </pre>
        </div>
    );
};