import React, { useState, useEffect } from 'react';
import './App.css';

import { TopBar } from './components/TopBar';
import { SideBar } from './components/SideBar';
import { Dashboard } from './components/Dashboard';
import { Whitelist } from './components/Whitelist';

// Wails Go Functions
import { LoadSettings, SaveSettings, StartEnforcement, StopEnforcement } from '../wailsjs/go/main/App';

function App() {
  const [activeTab, setActiveTab] = useState('dashboard');
  const [config, setConfig] = useState(null);
  const [isEnforced, setIsEnforced] = useState(false);

  // get settings from TOML on startup
  useEffect(() => {
    LoadSettings().then((loadedConfig) => {
      setConfig(loadedConfig);
      setIsEnforced(loadedConfig.enforcement.is_enforced);
    }).catch(err => console.error("Failed to load settings:", err));
  }, []);

  // persistent toggle logic
  const handleToggleEnforce = (newState) => {
    setIsEnforced(newState);
    
    if (newState) {
      StartEnforcement();
    } else {
      StopEnforcement();
    }

    // update and save to toml
    if (config) {
      const updatedConfig = { ...config };
      updatedConfig.enforcement.is_enforced = newState;
      setConfig(updatedConfig);
      SaveSettings(updatedConfig);
    }
  };

  const handleUpdateConfig = (key, value) => {
    if (config) {
      const updatedConfig = { ...config };
      updatedConfig.enforcement[key] = value;
      setConfig(updatedConfig);
      SaveSettings(updatedConfig);
    }
  };

  // function to render the correct view based on Sidebar selection
  const renderContent = () => {
    switch(activeTab) {
      case 'dashboard':
        return (
          <>
            <Dashboard config={config} onUpdateConfig={handleUpdateConfig}/>
            <Whitelist />
          </>
        );
      case 'whitelist':
        return <Whitelist />;
      case 'terminal':
        return <div className="placeholder-view">System Logs coming soon...</div>;
      case 'advanced':
        return <div className="placeholder-view">Modes & Routines coming soon...</div>;
      default:
        return <Dashboard />;
    }
  };

  return (
    <div style={{ display: 'flex', height: '100vh', width: '100vw', background: '#0f172a', color: 'white' }}>
      
      {/* SIDEBAR */}
      <SideBar activeTab={activeTab} setActiveTab={setActiveTab} />

      {/* MAIN AREA */}
      <main style={{ flex: 1, display: 'flex', flexDirection: 'column', minWidth: 0 }}>
        
        {/* TOPBAR (Header) */}
        <TopBar isEnforced={isEnforced} setIsEnforced={handleToggleEnforce} />

        {/* CONTENT AREA */}
        <div style={{ padding: '32px', overflowY: 'auto', flex: 1 }}>
          {renderContent()}
        </div>
      </main>
    </div>
  );
}

export default App;