import React, { useState, useEffect } from 'react';
import styles from './MaintenancePanel.module.css';

function MaintenancePanel() {
  const [tasks, setTasks] = useState([]);
  const [loading, setLoading] = useState(false);
  const [selectedTask, setSelectedTask] = useState(null);
  const [backups, setBackups] = useState([]);
  const [activeTab, setActiveTab] = useState('tasks');

  useEffect(() => {
    fetchMaintenance();
    const interval = setInterval(fetchMaintenance, 5000);
    return () => clearInterval(interval);
  }, []);

  const fetchMaintenance = async () => {
    try {
      const response = await fetch('http://localhost:8080/api/maintenance/tasks');
      const data = await response.json();
      setTasks(data.tasks || []);

      if (activeTab === 'backups') {
        fetchBackups();
      }
    } catch (error) {
      console.error('Error fetching maintenance:', error);
    }
  };

  const fetchBackups = async () => {
    try {
      const response = await fetch('http://localhost:8080/api/maintenance/backups');
      const data = await response.json();
      setBackups(data.backups || []);
    } catch (error) {
      console.error('Error fetching backups:', error);
    }
  };

  const runMaintenance = async (taskId) => {
    setLoading(true);
    try {
      const endpoints = {
        1: '/api/maintenance/vacuum',
        2: '/api/maintenance/reindex',
        3: '/api/maintenance/logs',
        4: '/api/maintenance/cache',
        5: '/api/maintenance/pool',
        6: '/api/maintenance/backup',
        7: '/api/maintenance/cleanup',
        8: '/api/maintenance/memory',
        9: '/api/maintenance/security',
        10: '/api/maintenance/tuning',
      };

      const response = await fetch(`http://localhost:8080${endpoints[taskId]}`, {
        method: 'POST',
      });
      const result = await response.json();
      setSelectedTask(result);
      fetchMaintenance();
    } catch (error) {
      console.error('Error running maintenance:', error);
    } finally {
      setLoading(false);
    }
  };

  const taskIcons = {
    1: '🗂️',
    2: '🔧',
    3: '📋',
    4: '💾',
    5: '🔌',
    6: '💾',
    7: '🧹',
    8: '⚙️',
    9: '🔒',
    10: '⚡',
  };

  return (
    <div className={styles.panel}>
      <div className={styles.header}>
        <h2>🔧 Server Maintenance Control Panel</h2>
        <p>10 Advanced Maintenance Tasks</p>
      </div>

      <div className={styles.tabs}>
        <button
          className={`${styles.tab} ${activeTab === 'tasks' ? styles.active : ''}`}
          onClick={() => setActiveTab('tasks')}
        >
          📋 Tasks (10)
        </button>
        <button
          className={`${styles.tab} ${activeTab === 'backups' ? styles.active : ''}`}
          onClick={() => { setActiveTab('backups'); fetchBackups(); }}
        >
          💾 Backups
        </button>
      </div>

      {activeTab === 'tasks' && (
        <div className={styles.tasksGrid}>
          {tasks.map((task) => (
            <div key={task.id} className={`${styles.taskCard} ${styles[task.status]}`}>
              <div className={styles.taskHeader}>
                <span className={styles.icon}>{taskIcons[task.id]}</span>
                <h3>{task.name}</h3>
              </div>

              <p className={styles.description}>{task.description}</p>

              <div className={styles.status}>
                <span className={`${styles.badge} ${styles[task.success ? 'success' : 'failed']}`}>
                  {task.status === 'running' ? '⏳ Running' : task.success ? '✅ Success' : '❌ Failed'}
                </span>
                {task.duration && <span className={styles.duration}>{task.duration}</span>}
              </div>

              <div className={styles.times}>
                <small>Last: {new Date(task.last_run).toLocaleString()}</small>
                <small>Next: {new Date(task.next_run).toLocaleString()}</small>
              </div>

              <button
                className={styles.runBtn}
                onClick={() => runMaintenance(task.id)}
                disabled={loading || task.status === 'running'}
              >
                {task.status === 'running' ? '⏳ Running...' : '▶ Run Now'}
              </button>
            </div>
          ))}
        </div>
      )}

      {activeTab === 'backups' && (
        <div className={styles.backupsContainer}>
          <h3>Recent Backups</h3>
          <div className={styles.backupsList}>
            {backups.length > 0 ? (
              backups.map((backup) => (
                <div key={backup.id} className={styles.backupItem}>
                  <div className={styles.backupInfo}>
                    <h4>📦 {backup.name}</h4>
                    <p>Size: {backup.size} | Created: {new Date(backup.created_at).toLocaleString()}</p>
                  </div>
                  <div className={styles.backupStatus}>
                    <span className={styles.statusBadge}>{backup.status}</span>
                  </div>
                </div>
              ))
            ) : (
              <p>No backups found</p>
            )}
          </div>
        </div>
      )}

      {selectedTask && (
        <div className={styles.details}>
          <h3>📊 Last Execution Details</h3>
          <pre>{JSON.stringify(selectedTask, null, 2)}</pre>
        </div>
      )}
    </div>
  );
}

export default MaintenancePanel;
