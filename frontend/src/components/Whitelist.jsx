import React, { useState, useEffect } from 'react';
import './Whitelist.css';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import {faSearch, faTrash, faPlus, faShieldHalved, faCube } from '@fortawesome/free-solid-svg-icons';

// Wails bindings
import { GetAppMap, ToggleExemption, RemoveApp, SelectProcess } from '../../wailsjs/go/main/App';

export function Whitelist() {
  const [searchTerm, setSearchTerm] = useState('');
  const [apps, setApps] = useState([]);

  // fetch data from Go backend
  const refreshApps = () => {
    GetAppMap().then((map) => {
      // convert Go Map { "Name": {directory, is_exempt} } to array for React
      const appArray = Object.keys(map).map(name => ({
        name: name,
        directory: map[name].directory,
        isExempt: map[name].is_exempt
      }));
      setApps(appArray);
    }).catch(console.error);
  };

  useEffect(() => {
    refreshApps();
  }, []);

  const handleToggle = (name) => {
    ToggleExemption(name).then(refreshApps);
  };

  const handleRemove = (name) => {
    RemoveApp(name).then(refreshApps);
  };

  const handleAddProcess = async () => {
    try {
      const newProcessName = await SelectProcess();
      if (newProcessName) {
        // refresh the list to show the new entry
        refreshApps();
        // optional: show a small toast or notification
        console.log(`Added ${newProcessName} to registry.`);
      }
    } catch (err) {
      console.error("Failed to add process:", err);
    }
  };

  // filter based on search
  const filteredApps = apps.filter(app =>
    app.name.toLowerCase().includes(searchTerm.toLowerCase())
  );

  return (
    <div className="whitelist-container">
      
      <div className="whitelist-header">
        <div className="header-text">
          <h3 className="section-title">
            <FontAwesomeIcon icon={faShieldHalved} className="title-icon" />
            Exemption Registry
          </h3>
          <p>Processes listed here bypass all performance caps and resource limits.</p>
        </div>

        <button className="add-btn" onClick={handleAddProcess}>
          <FontAwesomeIcon icon={faPlus} />
          <span>Add Process</span>
        </button>
      </div>

      <div className="registry-box">
        <div className="search-wrapper">
          <FontAwesomeIcon icon={faSearch} className="search-icon" />
          <input
            type="text"
            placeholder="Search installed apps..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
          />

        </div>

        <div className="table-container">
          <table className="whitelist-table">
            <thead>
              <tr>
                <th>PROCESS NAME</th>
                <th>STATUS</th>
                <th className="text-right">ACTIONS</th>
              </tr>
            </thead>
            <tbody>
              {filteredApps.map((app) => (
                <tr key={app.name} className="table-row">
                  <td>
                    <div className="process-info">
                      <div className="icon-box">
                        <FontAwesomeIcon icon={faCube} />
                      </div>
                      <div className="name-stack">
                        <span className="process-name">{app.name}</span>
                        <span className="process-path">{app.directory}</span>
                      </div>
                    </div>
                  </td>
                  <td>
                    {/* Status is now just a label, not a button */}
                    <span className={`status-pill ${app.isExempt ? 'excluded' : 'active'}`}>
                      <span className="dot"></span>
                      {app.isExempt ? 'Exempt' : 'Monitored'}
                    </span>
                  </td>
                  <td className="text-right action-cell">
                    {/* Contextual Toggle Button */}
                    {app.isExempt ? (
                      <button
                        className="action-btn unwhitelist"
                        onClick={() => handleToggle(app.name)}
                        title="Monitor this process"
                      >
                        Remove Protection
                      </button>
                    ) : (
                      <button
                        className="action-btn whitelist"
                        onClick={() => handleToggle(app.name)}
                        title="Add to exemptions"
                      >
                        Add Protection
                      </button>
                    )}

                    {/* not needed now
                    <button
                      className="delete-btn"
                      onClick={() => handleRemove(app.name)}
                      title="Remove from list"
                    >
                      <FontAwesomeIcon icon={faTrash} />
                    </button>
                    */}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}