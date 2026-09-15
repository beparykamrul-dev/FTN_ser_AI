const root = document.querySelector('#root')

const API = '/api/v1'
const sections = [
  ['overview', 'Overview'], ['devices', 'Devices'], ['dns', 'FTNDNS'], ['ddns', 'DDNS'],
  ['geoflow', 'GeoFlow'], ['anycast', 'Anycast'], ['providers', 'Providers'],
  ['database', 'Database'], ['monitoring', 'Monitoring'], ['alerts', 'Alerts'], ['audit', 'Audit']
]

root.innerHTML = `
  <div class="shell">
    <aside class="side">
      <div class="brand"><span>FTN</span><small>Control Plane</small></div>
      <nav id="nav">${sections.map(([id,label]) => `<button data-view="${id}">${label}</button>`).join('')}</nav>
      <div class="side-note">WEB ONLY<br><small>No console required</small></div>
    </aside>
    <main class="main">
      <header class="top"><div><h1 id="title">Overview</h1><p>Family Time Network · live control</p></div><div class="live"><i></i> Live Web</div></header>
      <section id="content"><div class="loading">Loading FTN Control Plane…</div></section>
    </main>
  </div>
`

const content = document.querySelector('#content')
const title = document.querySelector('#title')
const nav = document.querySelector('#nav')

const esc = (v) => String(v ?? '').replace(/[&<>"']/g, c => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c]))
const card = (label, value, meta='') => `<article class="card"><div class="label">${esc(label)}</div><strong>${esc(value)}</strong><small>${esc(meta)}</small></article>`
const table = (headers, rows) => `<div class="table-wrap"><table><thead><tr>${headers.map(h=>`<th>${esc(h)}</th>`).join('')}</tr></thead><tbody>${rows.map(r=>`<tr>${r.map(c=>`<td>${esc(c)}</td>`).join('')}</tr>`).join('')}</tbody></table></div>`

async function get(path) {
  const r = await fetch(`${API}${path}`, {headers:{Accept:'application/json'}})
  if (!r.ok) throw new Error(`${r.status} ${r.statusText}`)
  return r.json()
}

function renderError(e) { content.innerHTML = `<div class="error"><b>Web API unavailable</b><p>${esc(e.message)}</p><button onclick="location.reload()">Retry</button></div>` }

async function overview() {
  const [s,h] = await Promise.all([get('/control/snapshot'), get('/control/health')])
  content.innerHTML = `<div class="grid">
    ${card('Control Plane', h.status, 'API healthy')}
    ${card('FTNDNS', 'Authoritative', 'ftndns.com')}
    ${card('DDNS', 'Ready', 'ftnddns.net')}
    ${card('GeoFlow', s.geoflow?.mode || 'geo-latency-health', 'routing engine')}
    ${card('Data Plane', s.data_plane?.ram_first ? 'RAM-first' : 'Database', s.data_plane?.read_path_db_free ? 'DB-free read path' : 'DB on path')}
    ${card('Providers', (s.providers?.length || 0) + ' configured', 'Google · Cloudflare · Render · FTN')}
  </div><div class="panel"><h2>Strategic routing</h2><p>FTN remains the control plane. Google and Cloudflare are first-choice strategic helpers; Render is an application/web platform. Provider selection remains replaceable.</p></div>`
}

async function view(id) {
  if (id === 'overview') return overview()
  const map = {
    devices:'/control/agents', dns:'/control/dns', ddns:'/control/ddns', geoflow:'/control/geoflow',
    anycast:'/control/anycast', providers:'/control/providers', database:'/control/db',
    monitoring:'/control/nodes', alerts:'/control/alerts', audit:'/control/audit'
  }
  const data = await get(map[id])
  if (id === 'devices') {
    const a=data.agents||[]
    content.innerHTML=`<div class="grid">${card('Agents',a.length,'registered sessions')}${card('WebSocket',data.websocket||'/api/v1/ws/agent','outbound WSS')}${card('Access','Web','dashboard controlled')}</div>${table(['Device','Status','Last seen','Remote'],a.map(x=>[x.device_id,x.connected?'Online':'Offline',x.last_seen||'—',x.remote_addr||'—']))}`
  } else if (id === 'dns') content.innerHTML=`<div class="grid">${card('Zones',(data.zones||[]).length,'FTN authoritative mesh')}${card('DNSSEC',data.dnssec?'Enabled':'Disabled','mesh policy')}${card('Nodes',(data.nodes||[]).length,'DNS nodes')}${card('Routes',(data.routes||[]).length,'Anycast routes')}</div>${table(['Zone','Node','Endpoint','Health'],(data.nodes||[]).map(x=>['FTNDNS',x.id,x.endpoint,x.healthy?'Healthy':'Down']))}`
  else if (id === 'ddns') content.innerHTML=`<div class="grid">${card('Zone',data.zone,'authoritative DDNS')}${card('Status',data.status,'engine')}${card('Access','Web API','no console')}</div>`
  else if (id === 'geoflow') content.innerHTML=`<div class="grid">${card('Mode',data.mode||'geo-latency-health','routing')}${card('Health',data.health_required?'Required':'Optional','target selection')}${card('Latency',data.latency_required?'Required':'Optional','probe based')}${card('BDIX',data.bdix_dependency?'Used':'Independent','global routing')}</div>`
  else if (id === 'anycast') content.innerHTML=`<div class="grid">${card('Status',data.status,'route registry')}${card('Routing',data.routing,'GeoFlow')}${card('Prefixes',(data.prefixes||[]).length,'configured')}</div>${table(['Prefix','ASN','Enabled','BGP','Health'],(data.prefixes||[]).map(x=>[x.prefix||'Not assigned',x.asn||'—',x.enabled?'Yes':'No',x.bgp_advertised?'Yes':'No',x.healthy?'Healthy':'Not active']))}`
  else if (id === 'providers') { const p=data.providers||[]; content.innerHTML=`<div class="grid">${p.map(x=>card(x.name,x.enabled?'Enabled':'Disabled',x.kind)).join('')}</div>` }
  else if (id === 'database') content.innerHTML=`<div class="grid">${card('Database',data.status||'Registry ready','PostgreSQL source of truth')}${card('Instances',(data.instances||[]).length,'managed from web')}</div><pre>${esc(JSON.stringify(data,null,2))}</pre>`
  else if (id === 'monitoring') content.innerHTML=`<div class="grid">${card('Nodes',(data.nodes||[]).length,'live registry')}</div>${table(['Node','CPU','Memory','Latency','Status'],(data.nodes||[]).map(x=>[x.node_id,x.cpu_percent+'%',x.memory_percent+'%',x.latency_ms+' ms',x.healthy?'Healthy':'Down']))}`
  else if (id === 'alerts') content.innerHTML=`<div class="grid">${card('Alerts',(data.alerts||[]).length,'web notifications')}</div>${table(['ID','Severity','Message'],(data.alerts||[]).map(x=>[x.id,x.severity,x.message]))}`
  else content.innerHTML=`<div class="panel"><h2>Audit</h2><p>${esc(data.status||'Registry ready')}</p><p>All management actions are intended to be visible through the Web Control Plane.</p></div>`
}

nav.addEventListener('click', async e => {
  const b=e.target.closest('button'); if(!b) return
  document.querySelectorAll('#nav button').forEach(x=>x.classList.toggle('active',x===b))
  title.textContent=b.textContent
  content.innerHTML='<div class="loading">Loading…</div>'
  try { await view(b.dataset.view) } catch(e) { renderError(e) }
})

document.querySelector('#nav button').classList.add('active')
overview().catch(renderError)
setInterval(() => { const active=document.querySelector('#nav button.active'); if(active) view(active.dataset.view).catch(renderError) }, 15000)
