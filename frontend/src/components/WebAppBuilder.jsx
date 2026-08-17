import React, { useState, useCallback } from 'react';
import styles from './WebAppBuilder.module.css';

const WebAppBuilder = () => {
  const [currentApp, setCurrentApp] = useState(null);
  const [pages, setPages] = useState([]);
  const [currentPage, setCurrentPage] = useState(null);
  const [components, setComponents] = useState([]);
  const [componentLibrary, setComponentLibrary] = useState([]);
  const [showLibrary, setShowLibrary] = useState(true);
  const [editMode, setEditMode] = useState(false);

  // Fetch component library
  const fetchComponentLibrary = async () => {
    try {
      const response = await fetch('http://localhost:8080/api/builder/components');
      const data = await response.json();
      setComponentLibrary(Object.entries(data.components).map(([key, val]) => ({
        id: key,
        ...val,
      })));
    } catch (error) {
      console.error('Error fetching library:', error);
    }
  };

  // Create new web app
  const createNewApp = async () => {
    const name = prompt('App name:');
    if (!name) return;

    try {
      const response = await fetch('http://localhost:8080/api/builder/apps', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          name,
          description: prompt('Description:') || '',
          theme: { primary: '#667eea', secondary: '#764ba2' },
        }),
      });
      const app = await response.json();
      setCurrentApp(app);
      setPages(app.pages || []);
      fetchComponentLibrary();
    } catch (error) {
      console.error('Error creating app:', error);
    }
  };

  // Add page
  const addPage = async () => {
    if (!currentApp) return;
    const name = prompt('Page name:');
    if (!name) return;

    try {
      const response = await fetch(`http://localhost:8080/api/builder/pages?app_id=${currentApp.id}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          name,
          slug: name.toLowerCase().replace(/\s/g, '-'),
          components: [],
        }),
      });
      const page = await response.json();
      setPages([...pages, page]);
      setCurrentPage(page);
      setComponents([]);
    } catch (error) {
      console.error('Error adding page:', error);
    }
  };

  // Add component
  const addComponent = async (componentType) => {
    if (!currentApp || !currentPage) return;

    try {
      const response = await fetch(
        `http://localhost:8080/api/builder/components?app_id=${currentApp.id}&page_id=${currentPage.id}`,
        {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            type: componentType,
            name: componentLibrary.find(c => c.id === componentType)?.name || componentType,
            props: {},
          }),
        }
      );
      const component = await response.json();
      setComponents([...components, component]);
    } catch (error) {
      console.error('Error adding component:', error);
    }
  };

  // Export HTML
  const exportAsHTML = () => {
    if (!currentApp) return;
    window.open(`http://localhost:8080/api/builder/export/html?id=${currentApp.id}`);
  };

  // Export JSON
  const exportAsJSON = () => {
    if (!currentApp) return;
    window.open(`http://localhost:8080/api/builder/export/json?id=${currentApp.id}`);
  };

  return (
    <div className={styles.builder}>
      <div className={styles.header}>
        <h1>🎨 Web App Builder</h1>
        <div className={styles.actions}>
          <button onClick={createNewApp} className={styles.btnPrimary}>➕ New App</button>
          {currentApp && (
            <>
              <button onClick={addPage} className={styles.btnSecondary}>📄 Add Page</button>
              <button onClick={exportAsHTML} className={styles.btnSuccess}>📥 Export HTML</button>
              <button onClick={exportAsJSON} className={styles.btnInfo}>💾 Export JSON</button>
            </>
          )}
        </div>
      </div>

      {currentApp && (
        <div className={styles.container}>
          <div className={styles.sidebar}>
            <h3>📚 Components Library</h3>
            <div className={styles.componentsList}>
              {componentLibrary.map((comp) => (
                <div
                  key={comp.id}
                  className={styles.componentItem}
                  draggable
                  onDragStart={(e) => e.dataTransfer.setData('component', comp.id)}
                  onClick={() => addComponent(comp.id)}
                >
                  <span>{comp.icon || '📦'}</span>
                  <span>{comp.name}</span>
                </div>
              ))}
            </div>
          </div>

          <div className={styles.canvas}>
            <div className={styles.pagesTabs}>
              {pages.map((page) => (
                <button
                  key={page.id}
                  className={`${styles.pageTab} ${currentPage?.id === page.id ? styles.active : ''}`}
                  onClick={() => {
                    setCurrentPage(page);
                    setComponents(page.components || []);
                  }}
                >
                  📄 {page.name}
                </button>
              ))}
            </div>

            <div className={styles.canvasArea}>
              <h2>{currentPage?.name || 'Select a page'}</h2>
              <div className={styles.componentsGrid}>
                {components.map((comp) => (
                  <div key={comp.id} className={styles.componentCard}>
                    <div className={styles.componentHeader}>
                      <h4>{comp.name}</h4>
                      <button
                        className={styles.removeBtn}
                        onClick={() => setComponents(components.filter(c => c.id !== comp.id))}
                      >
                        ✕
                      </button>
                    </div>
                    <p>Type: {comp.type}</p>
                    <div className={styles.componentPreview}>
                      <code>{JSON.stringify(comp.props, null, 2)}</code>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>

          <div className={styles.properties}>
            <h3>⚙️ Settings</h3>
            <div className={styles.setting}>
              <label>App Name:</label>
              <input
                type="text"
                defaultValue={currentApp.name}
                placeholder="Enter app name"
              />
            </div>
            <div className={styles.setting}>
              <label>Description:</label>
              <textarea
                defaultValue={currentApp.description}
                placeholder="Enter description"
              />
            </div>
            <div className={styles.setting}>
              <label>Primary Color:</label>
              <input
                type="color"
                defaultValue={currentApp.theme?.primary || '#667eea'}
              />
            </div>
            <div className={styles.setting}>
              <label>Secondary Color:</label>
              <input
                type="color"
                defaultValue={currentApp.theme?.secondary || '#764ba2'}
              />
            </div>
          </div>
        </div>
      )}

      {!currentApp && (
        <div className={styles.welcome}>
          <div className={styles.welcomeCard}>
            <h2>🚀 Welcome to Web App Builder</h2>
            <p>Build unlimited web applications with unlimited UI components</p>
            <button onClick={createNewApp} className={styles.btnLarge}>
              ✨ Create Your First App
            </button>
            <div className={styles.features}>
              <div className={styles.feature}>
                <span>📚</span>
                <p>30+ UI Components</p>
              </div>
              <div className={styles.feature}>
                <span>🎨</span>
                <p>Unlimited Customization</p>
              </div>
              <div className={styles.feature}>
                <span>📥</span>
                <p>Export HTML/JSON</p>
              </div>
              <div className={styles.feature}>
                <span>⚡</span>
                <p>Real-time Preview</p>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default WebAppBuilder;
