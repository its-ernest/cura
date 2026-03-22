import React, { useState, useEffect } from 'react';
import './App.css';

import { TopBar } from './components/TopBar';
import { SideBar } from './components/SideBar';
import { Dashboard } from './components/Dashboard';
import { Whitelist } from './components/Whitelist';
import { LogView } from './components/LogView';

// Wails Go Functions
import * as runtime from '../wailsjs/runtime';
import { LoadSettings, SaveSettings, StartEnforcement, StopEnforcement, TriggerUpdate } from '../wailsjs/go/main/App';

function App() {
  const [activeTab, setActiveTab] = useState('dashboard');
  const [config, setConfig] = useState(null);
  const [isEnforced, setIsEnforced] = useState(false);

  // update state
  const [updateInfo, setUpdateInfo] = useState(null);
  const [showUpdate, setShowUpdate] = useState(false);
  const [isUpdating, setIsUpdating] = useState(false);

  // get settings from toml on startup
  useEffect(() => {
    LoadSettings().then((loadedConfig) => {
      setConfig(loadedConfig);
      setIsEnforced(loadedConfig.enforcement.is_enforced);
    }).catch(err => console.error("Failed to load settings:", err));

    runtime.EventsOn("update_available", (release) => {
      setUpdateInfo(release);
      setShowUpdate(true);
    });
    return () => runtime.EventsOff("update_available");
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

  const handleInstallUpdate = async () => {
    if (isUpdating) return;

    setIsUpdating(true);
    try {
      const result = await TriggerUpdate(updateInfo);
      if (result !== "Success") {
        alert("Update failed: " + result);
        setIsUpdating(false);
      }
      // if successful, the Go backend will call os.Exit(0), closing the app.
    } catch (err) {
      console.error("Installation error:", err);
      setIsUpdating(false);
    }
  };

  // function to render the correct view based on Sidebar selection
  const renderContent = () => {
    switch (activeTab) {
      case 'dashboard':
        return (
          <>
            <Dashboard config={config} onUpdateConfig={handleUpdateConfig} />
            <Whitelist />
          </>
        );
      case 'whitelist':
        return <Whitelist />;
      case 'terminal':
        return <LogView />;
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

      {showUpdate && updateInfo && (
        <div className="update-overlay">
          <div className={`update-modal ${isUpdating ? 'modal-busy' : ''}`}>
            <div className="update-badge">NEW EVOLUTION</div>
            <h2>Upgrade to {updateInfo.tag_name}?</h2>
            <div className="update-body">
              <p>{isUpdating ? "Downloading and extracting new binary..." : "The following enhancements are ready for deploy:"}</p>
              {!isUpdating && <pre className="release-notes">{updateInfo.body}</pre>}
              {isUpdating && <div className="loader-bar"><div className="loader-fill"></div></div>}
            </div>
            <div className="update-footer">
              <button
                className="btn-ignore"
                disabled={isUpdating}
                onClick={() => setShowUpdate(false)}
              >
                Keep Current
              </button>
              <button
                className="btn-upgrade"
                disabled={isUpdating}
                onClick={handleInstallUpdate}
              >
                {isUpdating ? "Deploying..." : "Install & Restart"}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

export default App;