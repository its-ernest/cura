import { useState, useEffect } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlay, faPause, faBolt, faFileCode, faGear } from '@fortawesome/free-solid-svg-icons';
import { GetRoutines, ToggleRoutine } from '../../wailsjs/go/main/App';
import './Routines.css';

export function Routines() {
    const [routines, setRoutines] = useState([]);

    const refreshRoutines = () => {
        GetRoutines().then(setRoutines).catch(console.error);
    };

    useEffect(() => {
        refreshRoutines();
        // Poll every few seconds to see if a routine became "Active" in the background
        const interval = setInterval(refreshRoutines, 3000);
        return () => clearInterval(interval);
    }, []);

    const onToggle = (e, routineName) => {
        // Stop bubbling, but don't preventDefault so the checkbox actually changes
        e.stopPropagation();

        const newState = e.target.checked; // Pull the actual state of the checkbox

        // 1. Instant UI update so the slider moves
        setRoutines(prev => prev.map(r =>
            r.name === routineName ? { ...r, enabled: newState } : r
        ));

        // 2. Sync with Go
        ToggleRoutine(routineName, newState).then(() => {
            // Optional: refreshRoutines()
        }).catch(err => {
            console.error(err);
            refreshRoutines(); // Rollback on fail
        });
    };
    return (
        <div className="routines-container">
            <div className="routines-header">
                <div className="header-text">
                    <h3 className="section-title">
                        <FontAwesomeIcon icon={faGear} className="title-icon" />
                        Operational Pipelines
                    </h3>
                    <p>Automated routines triggered by system events or application launches.</p>
                </div>
            </div>

            <div className="routines-grid">
                {routines.map((r) => (
                    <div key={r.name} className={`routine-card ${r.isActive ? 'is-running' : ''} ${!r.enabled ? 'is-disabled' : ''}`}>
                        <div className="card-top">
                            <div className="routine-icon">
                                <FontAwesomeIcon icon={faBolt} />
                            </div>
                            <div className="routine-meta">
                                <span className="routine-name">{r.name}</span>
                                <span className="routine-trigger">Trigger: {r.trigger.target}</span>
                            </div>
                            <label className="switch">
                                <input
                                    type="checkbox"
                                    checked={r.enabled}
                                    onChange={(e) => onToggle(e, r.name)} // Use onChange for inputs
                                />
                                <span className="slider round"></span>
                            </label>
                        </div>

                        <div className="card-body">
                            <div className="action-list">
                                {r.actions.map((action, idx) => (
                                    <div key={idx} className="action-item">
                                        <FontAwesomeIcon icon={faFileCode} className="action-icon" />
                                        <span>{action.type.replace('_', ' ')}: <b>{action.value || action.target}</b></span>
                                    </div>
                                ))}
                            </div>
                        </div>

                        <div className="card-footer">
                            <div className={`status-indicator ${r.isActive ? 'active' : ''}`}>
                                <div className="pulse-dot"></div>
                                <span>{r.isActive ? "ACTIVE & ENFORCING" : r.enabled ? "STANDBY" : "DISABLED"}</span>
                            </div>
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
}