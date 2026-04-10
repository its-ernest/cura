import './SideBar.css';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { 
  faShieldHeart, 
  faTableColumns, 
  faMicrochip, 
  faListCheck,
  faBolt,
  faCode,
  faScrewdriverWrench,
  faTerminal
} from '@fortawesome/free-solid-svg-icons';

export function SideBar({ activeTab, setActiveTab }) {
  return (
    <aside className="glass-sidebar">
      <div className="sidebar-content">
        {/* Logo Section */}
        <div className="logo-container">
          <div className="logo-icon">
            <FontAwesomeIcon icon={faShieldHeart} />
          </div>
          <div className="logo-text">
            <h1 className="brand-name">Cura</h1>
            <p className="brand-subtitle">SYSTEM UTILITY</p>
          </div>
        </div>

        {/* Navigation Section */}
        <nav className="nav-list">
          <button 
            className={`nav-item ${activeTab === 'dashboard' ? 'active' : ''}`} 
            onClick={() => setActiveTab('dashboard')}
          >
            <FontAwesomeIcon icon={faTableColumns} className="nav-icon" />
            <span>Dashboard</span>
          </button>

          <button 
            className={`nav-item ${activeTab === 'advanced' ? 'active' : ''}`} 
            onClick={() => setActiveTab('advanced')}
          >
            <FontAwesomeIcon icon={faScrewdriverWrench} className="nav-icon" />
            <span>Modes &amp; Routines</span>
          </button>

          <button 
            className={`nav-item ${activeTab === 'terminal' ? 'active' : ''}`} 
            onClick={() => setActiveTab('terminal')}
          >
            <FontAwesomeIcon icon={faTerminal} className="nav-icon" />
            <span>System Logs</span>
          </button>
        </nav>
      </div>

      {/* Footer Section */}
      <div className="sidebar-footer">
        {/* Status Card */}
        <div className="status-card">
          <div className="status-header">
            <span className="status-label">PRO STATUS</span>
            <span className="status-indicator"></span>
          </div>
          <p className="status-text">System is currently optimized and protected.</p>
        </div>

        {/* Developer Credit Info */}
        <div className="dev-credit">
            <div className="dev-label-row">
                <FontAwesomeIcon icon={faCode} className="dev-code-icon" />
                <span className="dev-title">ORIGINAL DEVELOPER</span>
            </div>
            <p className="dev-name">
                ERNEST <br />
                <span className="dev-contact">+233 24 512 1140</span>
            </p>
            <p>github.com/its-ernest/cura</p>
            <p className="dev-legal">
                Open source & free to contribute.
            </p>
        </div>
      </div>
    </aside>
  );
}