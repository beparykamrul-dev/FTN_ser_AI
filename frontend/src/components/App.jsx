import React, { useState } from 'react';
import Clock from './Clock';
import Dashboard from './Dashboard';
import MaintenancePanel from './MaintenancePanel';
import WebAppBuilder from './WebAppBuilder';
import styles from './App.module.css';

function App() {
  const [activeView, setActiveView] = useState('dashboard');

  return (
    <div className={styles.app}>
      <nav className={styles.navbar}>
        <div className={styles.logo}>🌐 FTN Server AI</div>
        <div className={styles.navItems}>
          <button
            className={`${styles.navBtn} ${activeView === 'dashboard' ? styles.active : ''}`}
            onClick={() => setActiveView('dashboard')}
          >
            📊 Dashboard
          </button>
          <button
            className={`${styles.navBtn} ${activeView === 'clocks' ? styles.active : ''}`}
            onClick={() => setActiveView('clocks')}
          >
            🕐 Clocks
          </button>
          <button
            className={`${styles.navBtn} ${activeView === 'maintenance' ? styles.active : ''}`}
            onClick={() => setActiveView('maintenance')}
          >
            🔧 Maintenance
          </button>
          <button
            className={`${styles.navBtn} ${activeView === 'builder' ? styles.active : ''}`}
            onClick={() => setActiveView('builder')}
          >
            🎨 Web Builder
          </button>
        </div>
      </nav>

      <main className={styles.main}>
        {activeView === 'dashboard' && <Dashboard />}
        {activeView === 'clocks' && <Clock />}
        {activeView === 'maintenance' && <MaintenancePanel />}
        {activeView === 'builder' && <WebAppBuilder />}
      </main>
    </div>
  );
}

export default App;
