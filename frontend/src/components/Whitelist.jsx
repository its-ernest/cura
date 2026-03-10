import React, { useState } from 'react';
import './Whitelist.css';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { 
  faSearch, 
  faTrash, 
  faPlus, 
  faShieldHalved,
  faTerminal,
  faCode,
  faGlobe
} from '@fortawesome/free-solid-svg-icons';

export function Whitelist() {
  const [searchTerm, setSearchTerm] = useState('');
  
  // Mock data for the whitelist
  const [items, setItems] = useState([
    { id: 1, name: 'Google Chrome', icon: faGlobe, status: 'Active', color: '#fbbf24' },
    { id: 2, name: 'Visual Studio Code', icon: faCode, status: 'Excluded', color: '#3b32e2' },
    { id: 3, name: 'Docker Desktop', icon: faTerminal, status: 'Paused', color: '#94a3b8' },
  ]);

  const removeItem = (id) => {
    setItems(items.filter(item => item.id !== id));
  };

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
        <button className="add-btn">
          <FontAwesomeIcon icon={faPlus} />
          <span>Add Process</span>
        </button>
      </div>

      <div className="registry-box">
        {/* Search Bar */}
        <div className="search-wrapper">
          <FontAwesomeIcon icon={faSearch} className="search-icon" />
          <input 
            type="text" 
            placeholder="Search or add process name (e.g., chrome.exe)..." 
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
          />
        </div>

        {/* Table */}
        <div className="table-container">
          <table className="whitelist-table">
            <thead>
              <tr>
                <th>PROCESS NAME</th>
                <th>ENFORCEMENT STATUS</th>
                <th className="text-right">ACTION</th>
              </tr>
            </thead>
            <tbody>
              {items.map((item) => (
                <tr key={item.id} className="table-row">
                  <td>
                    <div className="process-info">
                      <div className="icon-box" style={{ backgroundColor: `${item.color}20`, color: item.color }}>
                        <FontAwesomeIcon icon={item.icon} />
                      </div>
                      <span className="process-name">{item.name}</span>
                    </div>
                  </td>
                  <td>
                    <span className={`status-pill ${item.status.toLowerCase()}`}>
                      <span className="dot"></span>
                      {item.status}
                    </span>
                  </td>
                  <td className="text-right">
                    <button className="delete-btn" onClick={() => removeItem(item.id)}>
                      <FontAwesomeIcon icon={faTrash} />
                    </button>
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