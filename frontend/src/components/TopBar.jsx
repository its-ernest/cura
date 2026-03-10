import React, { useEffect } from 'react';
import { StartEnforcement, StopEnforcement } from '../../wailsjs/go/main/App';

export function TopBar({ isEnforced, setIsEnforced }) {

  useEffect(() => {
    if (isEnforced) {
      console.log("Cura Enforcer: Enabled");
      StartEnforcement();
    } else {
      console.log("Cura Enforcer: Disabled");
      StopEnforcement();
    }
  }, []);

  const handleToggle = () => {
    const nextState = !isEnforced;
    setIsEnforced(nextState);

    // call the Go backend based on the new state
    if (nextState) {
      console.log("Cura Enforcer: Enabled");
      StartEnforcement();
      isEnforced = true;
    } else {
      console.log("Cura Enforcer: Disabled");
      StopEnforcement();
      isEnforced = false;
    }
  };

  return (
    <header className="topbar">
      <div>
        <h2 style={{ fontSize: '24px', fontWeight: 'bold', margin: 0 }}>System Overview</h2>
        <p style={{ fontSize: '14px', color: 'var(--slate-400)', margin: 0 }}>
          {isEnforced ? "Enforcement Active" : "Monitoring only"}
        </p>
      </div>

      <div className="status-badge">
        <span style={{ fontSize: '14px', fontWeight: 500, color: '#cbd5e1' }}>Enforce Limits</span>
        <button
          className="toggle-btn"
          style={{
            background: isEnforced ? 'var(--emerald)' : '#475569',
            cursor: 'pointer',
            position: 'relative',
            width: '52px',
            height: '26px',
            borderRadius: '13px',
            border: 'none',
            transition: 'background 0.3s ease'
          }}
          onClick={handleToggle}
        >
          <div
            className="toggle-circle"
            style={{
              position: 'absolute',
              top: '3px',
              width: '20px',
              height: '20px',
              borderRadius: '50%',
              background: 'white',
              transition: 'all 0.3s ease',
              transform: isEnforced ? 'translateX(26px)' : 'translateX(3px)'
            }}
          />
        </button>
      </div>
    </header>
  );
}