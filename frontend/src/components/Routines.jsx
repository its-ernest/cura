import { useState, useEffect } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlus, faBolt, faFileCode, faGear, faTimes, faTrash } from '@fortawesome/free-solid-svg-icons';
import { GetRoutines, ToggleRoutine, CreateRoutine } from '../../wailsjs/go/main/App';
import './Routines.css';

export function Routines() {
    const [routines, setRoutines] = useState([]);
    const [showModal, setShowModal] = useState(false);

    const initialFormState = {
        name: '',
        enabled: true,
        trigger: { type: 'app_launch', target: '' },
        actions: [{ type: 'set_memory_cap', value: 70 }],
        stop_condition: { type: 'app_close', target: '' }
    };
    const [newRoutine, setNewRoutine] = useState(initialFormState);

    const refreshRoutines = () => GetRoutines().then(setRoutines).catch(console.error);

    useEffect(() => {
        refreshRoutines();
        const interval = setInterval(refreshRoutines, 3000);
        return () => clearInterval(interval);
    }, []);

    const onToggle = (e, routineName) => {
        e.stopPropagation();
        const newState = e.target.checked;
        setRoutines(prev => prev.map(r => r.name === routineName ? { ...r, enabled: newState } : r));
        ToggleRoutine(routineName, newState).then(() => { }).catch(err => {
            console.error(err);
            refreshRoutines();
        });
    };

    const handleCreate = (e) => {
        e.preventDefault();
        CreateRoutine(newRoutine).then(() => {
            setShowModal(false);
            refreshRoutines();
            setNewRoutine(initialFormState);
        }).catch(console.error);
    };

    const addAction = () => {
        setNewRoutine({
            ...newRoutine,
            actions: [...newRoutine.actions, { type: 'temporary_exemption', value: 0 }]
        });
    };

    const removeAction = (index) => {
        const updated = newRoutine.actions.filter((_, i) => i !== index);
        setNewRoutine({ ...newRoutine, actions: updated });
    };

    const updateAction = (index, field, val) => {
        const updated = [...newRoutine.actions];
        updated[index][field] = field === 'value' ? parseInt(val) || 0 : val;
        setNewRoutine({ ...newRoutine, actions: updated });
    };

    const handleTargetChange = (val) => {
        setNewRoutine({
            ...newRoutine,
            trigger: { ...newRoutine.trigger, target: val },
            stop_condition: { ...newRoutine.stop_condition, target: val }
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
                    <p>Automated routines triggered by system events.</p>
                </div>
                <button className="create-btn" onClick={() => setShowModal(true)}>
                    <FontAwesomeIcon icon={faPlus} /> Create Pipeline
                </button>
            </div>

            {showModal && (
                <div className="modal-overlay">
                    <div className="modal-content advanced">
                        <div className="modal-header">
                            <h4 className="header-topic">Initialize Routine</h4>
                            <div>
                                <button className="close-btn" onClick={() => setShowModal(false)}>
                                    <FontAwesomeIcon icon={faTimes} />
                                </button>
                            </div>

                        </div>
                        <form onSubmit={handleCreate} className="form-scroll-area">
                            <div className="form-group">
                                <label>Routine Name</label>
                                <input required value={newRoutine.name} onChange={e => setNewRoutine({ ...newRoutine, name: e.target.value })} placeholder="e.g. Gaming Super-Mode" />
                            </div>

                            <div className="form-group">
                                <label>Target Process (.exe)</label>
                                <input required value={newRoutine.trigger.target} onChange={e => handleTargetChange(e.target.value)} placeholder="cs2.exe" />
                            </div>

                            <div className="actions-section">
                                <label className="section-label">Automated Actions</label>
                                {newRoutine.actions.map((action, idx) => (
                                    <div key={idx} className="action-builder-row">
                                        <select value={action.type} onChange={e => updateAction(idx, 'type', e.target.value)}>
                                            <option value="set_memory_cap">Set Memory Cap</option>
                                            <option value="temporary_exemption">Temporary Exemption</option>
                                            <option value="boost_priority">Boost Priority</option>
                                        </select>
                                        {action.type === 'set_memory_cap' && (
                                            <input type="number" value={action.value} onChange={e => updateAction(idx, 'value', e.target.value)} placeholder="MB" />
                                        )}
                                        <button type="button" className="remove-action" onClick={() => removeAction(idx)}>
                                            <FontAwesomeIcon icon={faTrash} />
                                        </button>
                                    </div>
                                ))}
                                <button type="button" className="add-action-btn" onClick={addAction}>
                                    + Add Action
                                </button>
                            </div>

                            <button type="submit" className="submit-btn">Deploy Pipeline</button>
                        </form>
                    </div>
                </div>
            )}

            <div className="routines-grid">
                {routines.map((r) => (
                    <div key={r.name} className={`routine-card ${r.isActive ? 'is-running' : ''} ${!r.enabled ? 'is-disabled' : ''}`}>
                        <div className="card-top">
                            <div className="routine-icon"><FontAwesomeIcon icon={faBolt} /></div>
                            <div className="routine-meta">
                                <span className="routine-name">{r.name}</span>
                                <span className="routine-trigger">Trigger: {r.trigger.target}</span>
                            </div>
                            <label className="switch">
                                <input type="checkbox" checked={r.enabled} onChange={(e) => onToggle(e, r.name)} />
                                <span className="slider round"></span>
                            </label>
                        </div>

                        <div className="card-body">
                            <div className="action-list">
                                {r.actions.map((action, idx) => (
                                    <div key={idx} className="action-item">
                                        <FontAwesomeIcon icon={faFileCode} className="action-icon" />
                                        <span>{action.type.replace('_', ' ')}: <b>{action.value || action.target || 'Active'}</b></span>
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