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
      <h3>Players</h3>
      <div class="form-row" style="margin-bottom:12px">
        <label>Name <input type="text" id="newPlayerName" placeholder="Player name"
               onkeydown="if(event.key==='Enter') adminCreatePlayer()"></label>
        <button class="btn btn-primary" style="align-self:flex-end" onclick="adminCreatePlayer()">Add Player</button>
      </div>
      <div id="adminPlayerList"><p class="loading">Loading…</p></div>
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
      <h3>Danger Zone</h3>
      <p>This will permanently delete all matches, attendees, payments, and players.</p>
      <button class="btn btn-danger" onclick="cleanupAllData()">Clear All Data</button>
    </div>
  `;
  await adminRefreshPlayers();
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
            <td>${p.name}</td>
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
  await adminRefreshPlayers();
}

async function adminDeletePlayer(id) {
  if (!confirm('Delete player?')) return;
  const res = await api('DELETE', '/players/' + id);
  if (!res) return;
  if (!res.ok) { toast('Error', false); return; }
  toast('Player deleted');
  state.players = [];
  await adminRefreshPlayers();
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
  await adminRefreshPlayers();
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
  await adminRefreshPlayers();
}
