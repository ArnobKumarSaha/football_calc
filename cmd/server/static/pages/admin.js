'use strict';

async function admin(el) {
  if (!state.token) {
    el.innerHTML = `
      <div class="page-header"><h2>Admin</h2></div>
      <div class="card form-section">
        <p>Admin login required.</p>
        <button class="btn btn-primary" onclick="openLoginModal()">Login</button>
      </div>
    `;
    return;
  }
  el.innerHTML = `
    <div class="page-header"><h2>Admin</h2></div>
    <div class="card form-section">
      <h3>Add Player</h3>
      <div class="form-row">
        <label>Name <input type="text" id="newPlayerName" placeholder="Player name"
               onkeydown="if(event.key==='Enter') adminCreatePlayer()"></label>
        <button class="btn btn-primary" style="align-self:flex-end" onclick="adminCreatePlayer()">Add Player</button>
      </div>
    </div>
    <div class="card form-section" style="margin-top:16px">
      <h3 class="collapsible-header" onclick="adminTogglePlayerList(this)" style="cursor:pointer;user-select:none">
        Players <span class="collapse-arrow">▶</span>
      </h3>
      <div id="adminPlayerList" style="display:none"></div>
    </div>
    <div class="card form-section" style="margin-top:16px">
      <h3>Import Initial Data</h3>
      <p>Upload a text file with one entry per line: <code>&lt;name&gt; &lt;amount&gt;</code><br>
         Negative amount means the player is in debt. Only allowed when the database is empty.</p>
      <div style="display:flex;gap:10px;align-items:center;flex-wrap:wrap">
        <input type="file" id="import-file" accept=".txt,text/plain">
        <button class="btn btn-primary" onclick="importData()">Import</button>
      </div>
    </div>
    <div class="card form-section" style="margin-top:16px">
      <h3>Export / Backup</h3>
      <p>Download a self-contained SQL script (schema + all data) to re-initialize a fresh
         PostgreSQL with an exact clone of this database — e.g. as a KubeDB init script source.</p>
      <button class="btn btn-primary" onclick="exportData()">Download SQL Backup</button>
    </div>
    <div class="card form-section" style="margin-top:16px">
      <h3>Danger Zone</h3>
      <p>This will permanently delete all matches, attendees, payments, and players.</p>
      <button class="btn btn-danger" onclick="cleanupAllData()">Clear All Data</button>
    </div>
  `;
}

function adminTogglePlayerList(header) {
  const list = document.getElementById('adminPlayerList');
  const arrow = header.querySelector('.collapse-arrow');
  const collapsed = list.style.display === 'none';
  list.style.display = collapsed ? '' : 'none';
  arrow.textContent = collapsed ? '▼' : '▶';
  if (collapsed && !list.dataset.loaded) {
    list.dataset.loaded = '1';
    list.innerHTML = '<p class="loading">Loading…</p>';
    adminRefreshPlayers();
  }
}

function adminRefreshPlayersIfVisible() {
  const el = document.getElementById('adminPlayerList');
  if (el && el.style.display !== 'none') adminRefreshPlayers();
}

async function adminRefreshPlayers() {
  const el = document.getElementById('adminPlayerList');
  if (!el) return;
  const res = await api('GET', '/players');
  if (!res) return;
  const list = await res.json();
  state.players = list;
  if (!list.length) { el.innerHTML = '<p class="empty">No players yet.</p>'; return; }
  el.innerHTML = `
    <table>
      <thead><tr><th>ID</th><th>Name</th><th>Balance</th><th></th></tr></thead>
      <tbody>
        ${list.map(p => `
          <tr>
            <td>${p.id}</td>
            <td>${escapeHTML(p.name)}</td>
            <td>${balanceBadge(p.balance)}</td>
            <td class="row-actions">
              <button class="btn btn-sm btn-danger" onclick="adminDeletePlayer(${p.id})">Delete</button>
            </td>
          </tr>
        `).join('')}
      </tbody>
    </table>
  `;
}

async function adminCreatePlayer() {
  const name = document.getElementById('newPlayerName').value.trim();
  if (!name) { toast('Name required', false); return; }
  const res = await api('POST', '/players', { name });
  if (!res) return;
  const data = await res.json();
  if (!res.ok) { toast(data.error || 'Error', false); return; }
  toast(name + ' added');
  document.getElementById('newPlayerName').value = '';
  state.players = [];
  adminRefreshPlayersIfVisible();
}

async function adminDeletePlayer(id) {
  if (!confirm('Delete player?')) return;
  const res = await api('DELETE', '/players/' + id);
  if (!res) return;
  if (!res.ok) { toast('Error', false); return; }
  toast('Player deleted');
  state.players = [];
  adminRefreshPlayersIfVisible();
}

async function importData() {
  const fileInput = document.getElementById('import-file');
  if (!fileInput || !fileInput.files.length) {
    toast('Select a file first', false);
    return;
  }
  const formData = new FormData();
  formData.append('file', fileInput.files[0]);

  const res = await fetch('/api/admin/import', {
    method: 'POST',
    headers: { Authorization: 'Bearer ' + state.token },
    body: formData,
  });
  const data = await res.json();
  if (!res.ok) { toast(data.error || 'Import failed', false); return; }
  toast(`Imported ${data.imported} player(s)`);
  fileInput.value = '';
  state.players = [];
  adminRefreshPlayersIfVisible();
}

async function exportData() {
  let res;
  try {
    res = await fetch('/api/admin/export.sql', {
      headers: { Authorization: 'Bearer ' + state.token },
    });
  } catch (e) {
    toast('Network error', false);
    return;
  }
  if (!res.ok) { toast('Export failed', false); return; }
  const blob = await res.blob();
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = 'football-calc-backup.sql';
  a.click();
  URL.revokeObjectURL(url);
  toast('Backup downloaded');
}

async function cleanupAllData() {
  if (!confirm('This will delete ALL data (players, matches, attendees, payments). Are you sure?')) return;
  if (!confirm('Final confirmation: delete everything?')) return;
  const res = await api('POST', '/admin/cleanup');
  if (!res) return;
  const data = await res.json();
  if (!res.ok) { toast(data.error || 'Error', false); return; }
  toast('All data cleared');
  state.players = [];
  adminRefreshPlayersIfVisible();
}
