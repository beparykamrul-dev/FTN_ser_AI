import './control-plane.css'

const root = document.querySelector('#root')
const API = '/api/v1'
const sections = [
  ['overview','Overview'],['devices','Devices'],['network','Network'],['dns','FTNDNS'],['ddns','DDNS'],
  ['geoflow','GeoFlow'],['anycast','Anycast'],['providers','Providers'],['database','Database'],
  ['monitoring','Monitoring'],['alerts','Alerts'],['audit','Audit']
]

root.innerHTML=`<div class="shell"><aside class="side"><div class="brand"><span>FTN</span><small>Control Plane</small></div><nav id="nav">${sections.map(([id,label])=>`<button data-view="${id}">${label}</button>`).join('')}</nav><div class="side-note"><b>WEB CONTROL</b><br><small>CLI / console not required</small></div></aside><main class="main"><header class="top"><div><h1 id="title">Overview</h1><p>Family Time Network · unified web operations</p></div><div class="top-actions"><button id="refresh" class="tool">Refresh</button><div class="live"><i></i><span id="liveText">Live Web</span></div></div></header><section id="content"><div class="loading">Loading FTN Control Plane…</div></section></main></div>`

const content=document.querySelector('#content'),title=document.querySelector('#title'),nav=document.querySelector('#nav'),refresh=document.querySelector('#refresh')
const esc=v=>String(v??'').replace(/[&<>"']/g,c=>({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c]))
const card=(l,v,m='')=>`<article class="card"><div class="label">${esc(l)}</div><strong>${esc(v)}</strong><small>${esc(m)}</small></article>`
const table=(h,rows)=>`<div class="table-wrap"><table><thead><tr>${h.map(x=>`<th>${esc(x)}</th>`).join('')}</tr></thead><tbody>${rows.length?rows.map(r=>`<tr>${r.map(x=>`<td>${esc(x)}</td>`).join('')}</tr>`).join(''):`<tr><td colspan="${h.length}" class="empty">No data</td></tr>`}</tbody></table></div>`
const jsonBlock=v=>`<details class="raw"><summary>Inspect API response</summary><pre>${esc(JSON.stringify(v,null,2))}</pre></details>`
async function get(path){const r=await fetch(`${API}${path}`,{headers:{Accept:'application/json'},cache:'no-store'});if(!r.ok)throw new Error(`${r.status} ${r.statusText}`);return r.json()}
function renderError(e){document.querySelector('#liveText').textContent='API error';content.innerHTML=`<div class="error"><b>Web API unavailable</b><p>${esc(e.message)}</p><button class="tool" onclick="location.reload()">Retry</button></div>`}

async function overview(){
  const [s,h,c,a]=await Promise.all([get('/control/snapshot'),get('/control/health'),get('/control/capabilities'),get('/control/adapters')])
  document.querySelector('#liveText').textContent='Live Web'
  content.innerHTML=`<div class="grid">
    ${card('Control Plane',h.status||'ok','API health')}${card('FTNDNS','Authoritative','ftndns.com')}${card('DDNS','Ready','ftnddns.net')}
    ${card('GeoFlow',s.geoflow?.mode||'geo-latency-health','route selection')}${card('Data Plane',s.data_plane?.ram_first?'RAM-first':'Database','DB-free read path: '+(s.data_plane?.read_path_db_free?'Yes':'No'))}
    ${card('Agents',s.agents_count??'—','Web-connected devices')}${card('Providers',(a.adapters||[]).length,'adapter registry')}${card('Capabilities',Object.values(c).filter(Boolean).length,'enabled controls')}
  </div>
  <div class="actionbar"><button class="action" data-action="snapshot">Refresh snapshot</button><button class="action" data-action="health">Check health</button><button class="action" data-action="adapters">Check providers</button></div>
  <div class="panel"><h2>Unified Web Operations</h2><p>FTN exposes operational capabilities through this Web Control Plane. Device connectivity is outbound WSS; provider credentials remain server-side. Sensitive actions must be implemented as explicit capability-based operations with RBAC and audit.</p></div>
  ${jsonBlock({snapshot:s,health:h,capabilities:c,adapters:a})}`
  bindActions()
}

async function view(id){
  if(id==='overview')return overview()
  const map={devices:'/agents',network:'/control/network',dns:'/control/dns',ddns:'/control/ddns',geoflow:'/control/geoflow',anycast:'/control/anycast',providers:'/control/providers',database:'/control/db',monitoring:'/control/nodes',alerts:'/control/alerts',audit:'/control/audit'}
  if(id==='network'){
    try{const data=await get(map[id]);content.innerHTML=`<div class="grid">${card('Network',data.status||'Ready','web operations')}${card('Nodes',(data.nodes||[]).length,'managed nodes')}${card('Routes',(data.routes||[]).length,'route registry')}${card('Interfaces',(data.interfaces||[]).length,'interface inventory')}</div>${table(['Node','Address','Health','Latency'],(data.nodes||[]).map(x=>[x.node_id||x.id,x.address||x.endpoint,x.healthy?'Healthy':'Down',(x.latency_ms??'—')+' ms']))}${jsonBlock(data)}`;return}catch(e){return renderError(e)}
  }
  const data=await get(map[id])
  if(id==='devices'){
    const a=data.agents||[]
    content.innerHTML=`<div class="grid">${card('Devices',a.length,'registered sessions')}${card('Online',a.filter(x=>x.connected).length,'live WSS')}${card('Offline',a.filter(x=>!x.connected).length,'last known')}${card('Transport','WSS','outbound device connection')}</div>${table(['Device','Status','Last seen','Remote','Capabilities'],a.map(x=>[x.device_id,x.connected?'Online':'Offline',x.last_seen||'—',x.remote_addr||'—',(x.capabilities||[]).join(', ')||'—']))}${jsonBlock(data)}`
  }else if(id==='dns')content.innerHTML=`<div class="grid">${card('Zones',(data.zones||[]).length,'FTN authoritative mesh')}${card('DNSSEC',data.dnssec?'Enabled':'Disabled','authority policy')}${card('Nodes',(data.nodes||[]).length,'DNS nodes')}${card('Routes',(data.routes||[]).length,'Anycast routes')}</div>${table(['Node','Endpoint','Health','Latency'],(data.nodes||[]).map(x=>[x.id,x.endpoint,x.healthy?'Healthy':'Down',(x.latency_ms??'—')+' ms']))}${jsonBlock(data)}`
  else if(id==='ddns')content.innerHTML=`<div class="grid">${card('Zone',data.zone||'ftnddns.net','authoritative DDNS')}${card('Status',data.status||'Ready','engine')}${card('Records',data.records_count??'—','managed records')}${card('Access','Web API','no console')}</div>${jsonBlock(data)}`
  else if(id==='geoflow')content.innerHTML=`<div class="grid">${card('Mode',data.mode||'geo-latency-health','route engine')}${card('Health',data.health_required?'Required':'Optional','target selection')}${card('Latency',data.latency_required?'Required':'Optional','probe based')}${card('BDIX',data.bdix_dependency?'Used':'Independent','global routing')}</div><div class="panel"><h2>Route decision</h2><p>Nearest healthy target is selected from measured health, latency and load. Global 0 ms is not assumed; FTN records real measurements.</p></div>${jsonBlock(data)}`
  else if(id==='anycast')content.innerHTML=`<div class="grid">${card('Status',data.status||'Ready','route registry')}${card('Routing',data.routing||'GeoFlow','target selection')}${card('Prefixes',(data.prefixes||[]).length,'configured')}${card('BGP',(data.prefixes||[]).filter(x=>x.bgp_advertised).length,'advertised')}</div>${table(['Prefix','ASN','Enabled','BGP','Health'],(data.prefixes||[]).map(x=>[x.prefix||'Not assigned',x.asn||'—',x.enabled?'Yes':'No',x.bgp_advertised?'Yes':'No',x.healthy?'Healthy':'Not active']))}${jsonBlock(data)}`
  else if(id==='providers'){const p=data.providers||[];content.innerHTML=`<div class="grid">${p.map(x=>card(x.name,x.enabled?'Enabled':'Disabled',x.kind+(x.read_only?' · read-only':''))).join('')}</div><div class="panel"><h2>Provider strategy</h2><p>Google, Cloudflare and Render are first-class adapters. FTN remains the control authority and keeps provider credentials server-side.</p></div>${jsonBlock(data)}`}
  else if(id==='database')content.innerHTML=`<div class="grid">${card('Database',data.status||'Registry ready','PostgreSQL source of truth')}${card('Instances',(data.instances||[]).length,'web-managed')}${card('Runtime','Redis / RAM','hot state')}${card('Writes','Async persist','durable path')}</div>${jsonBlock(data)}`
  else if(id==='monitoring')content.innerHTML=`<div class="grid">${card('Nodes',(data.nodes||[]).length,'live registry')}${card('Healthy',(data.nodes||[]).filter(x=>x.healthy).length,'current health')}${card('Telemetry','Internal','no external telemetry export')}</div>${table(['Node','CPU','Memory','Latency','Status'],(data.nodes||[]).map(x=>[x.node_id,x.cpu_percent+'%',x.memory_percent+'%',x.latency_ms+' ms',x.healthy?'Healthy':'Down']))}${jsonBlock(data)}`
  else if(id==='alerts')content.innerHTML=`<div class="grid">${card('Alerts',(data.alerts||[]).length,'web notifications')}${card('Channel','FTN Web','operator surface')}</div>${table(['ID','Severity','Message'],(data.alerts||[]).map(x=>[x.id,x.severity,x.message]))}${jsonBlock(data)}`
  else content.innerHTML=`<div class="panel"><h2>Audit</h2><p>${esc(data.status||'Audit registry ready')}</p><p>Every privileged operation should be attributable to a user, role, device and provider adapter.</p></div>${jsonBlock(data)}`
}

function bindActions(){document.querySelectorAll('[data-action]').forEach(b=>b.onclick=async()=>{b.disabled=true;try{if(b.dataset.action==='snapshot')await get('/control/snapshot');if(b.dataset.action==='health')await get('/control/health');if(b.dataset.action==='adapters')await get('/control/adapters');await overview()}catch(e){renderError(e)}finally{b.disabled=false}})}
nav.addEventListener('click',async e=>{const b=e.target.closest('button');if(!b)return;document.querySelectorAll('#nav button').forEach(x=>x.classList.toggle('active',x===b));title.textContent=b.textContent;content.innerHTML='<div class="loading">Loading…</div>';try{await view(b.dataset.view)}catch(e){renderError(e)}})
refresh.onclick=async()=>{refresh.disabled=true;try{const b=document.querySelector('#nav button.active');await view(b?.dataset.view||'overview')}catch(e){renderError(e)}finally{refresh.disabled=false}}
document.querySelector('#nav button').classList.add('active')
overview().catch(renderError)
setInterval(()=>{const a=document.querySelector('#nav button.active');if(a)view(a.dataset.view).catch(renderError)},15000)
