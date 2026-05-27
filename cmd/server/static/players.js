'use strict';

async function players(el) {
  el.innerHTML = `
    <div class="page-header">
      <h2>Players</h2>
      <button class="btn btn-primary" onclick="toggleNewPlayerForm()">+ Add Player</button>
    </div>
    <div id="newPlayerForm" style="display:none" class="card form-section">
      <div class="form-row">
        <label>Name <input type="text" id="newPlayerName" placeholder="Player name"
               onkeydown="if(event.key==='Enter') createPlayer()"></label>
        <button class="btn btn-primary" style="align-self:flex-end" onclick="createPlayer()">Add</button>
        <button class="btn" style="align-self:flex-end" onclick="toggleNewPlayerForm()">Cancel</button>
      </div>
    </div>
    <div id="playerList"><p class="loading">Loading…</p></div>
  `;
  await refreshPlayerList();
}

function toggleNewPlayerForm() {
  const el = document.getElementById('newPlayerForm');
  el.style.display = el.style.display === 'none' ? 'block' : 'none';
}

async function refreshPlayerList() {
  const el = document.getElementById('playerList');
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
              <button class="btn btn-sm" onclick="showEditPlayer(${p.id}, '${p.name.replace(/'/g,"\\'")}')">Edit</button>
              <button class="btn btn-sm btn-danger" onclick="deletePlayer(${p.id})">Delete</button>
            </td>
          </tr>
          <tr id="pedit-${p.id}" style="display:none">
            <td colspan="4">
              <div class="inline-edit">
                <input type="text" id="pname-${p.id}" value="${p.name}"
                       onkeydown="if(event.key==='Enter') updatePlayer(${p.id})">
                <button class="btn btn-primary btn-sm" onclick="updatePlayer(${p.id})">Save</button>
                <button class="btn btn-sm" onclick="document.getElementById('pedit-${p.id}').style.display='none'">Cancel</button>
              </div>
            </td>
          </tr>
        `).join('')}
      </tbody>
    </table>
  `;
}

function showEditPlayer(id, name) {
  const row = document.getElementById('pedit-' + id);
  row.style.display = row.style.display === 'none' ? 'table-row' : 'none';
  if (row.style.display !== 'none') document.getElementById('pname-' + id)?.focus();
}

async function createPlayer() {
  const name = document.getElementById('newPlayerName').value.trim();
  if (!name) { toast('Name required', false); return; }
  const res = await api('POST', '/players', { name });
  if (!res) return;
  const data = await res.json();
  if (!res.ok) { toast(data.error || 'Error', false); return; }
  toast(name + ' added');
  document.getElementById('newPlayerName').value = '';
  toggleNewPlayerForm();
  await refreshPlayerList();
}

async function updatePlayer(id) {
  const name = document.getElementById('pname-' + id)?.value.trim();
  if (!name) { toast('Name required', false); return; }
  const res = await api('PATCH', '/players/' + id, { name });
  if (!res) return;
  if (!res.ok) { const d = await res.json(); toast(d.error || 'Error', false); return; }
  toast('Player updated');
  await refreshPlayerList();
}

async function deletePlayer(id) {
  if (!confirm('Delete player?')) return;
  const res = await api('DELETE', '/players/' + id);
  if (!res) return;
  if (!res.ok) { toast('Error', false); return; }
  toast('Player deleted');
  state.players = [];
  await refreshPlayerList();
}
