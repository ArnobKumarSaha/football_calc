'use strict';

async function matches(el) {
  await ensurePlayers();
  el.innerHTML = `
    <div class="page-header">
      <h2>Matches</h2>
      ${state.token ? `<button class="btn btn-primary" onclick="toggleNewMatchForm()">+ New Match</button>` : ''}
    </div>
    <div id="newMatchForm" style="display:none" class="card form-section">
      <h3>New Match</h3>
      <div class="form-row">
        <label>Date <input type="date" id="mDate"></label>
        <label>Total Bill <input type="number" id="mBill" step="0.01" min="0" placeholder="0.00"></label>
        <label>Notes <input type="text" id="mNotes" placeholder="optional"></label>
      </div>
      <div style="margin-top:12px">
        <div style="display:flex;align-items:center;gap:8px;margin-bottom:6px">
          <strong>Attendees</strong>
          <button class="btn btn-sm" onclick="addAttendeeRow('attendeeRows')">+ Add</button>
        </div>
        <div id="attendeeRows"></div>
      </div>
      <div style="display:flex;gap:8px;margin-top:12px">
        <button class="btn btn-primary" onclick="createMatch()">Create</button>
        <button class="btn" onclick="toggleNewMatchForm()">Cancel</button>
      </div>
    </div>
    <div id="matchList"></div>
  `;
  await refreshMatchList();
}

function addAttendeeRow(containerId) {
  const container = document.getElementById(containerId);
  if (!container) return;
  const row = document.createElement('div');
  row.className = 'form-row attendee-row';
  row.style.alignItems = 'center';
  row.innerHTML = `
    <label>Player <select class="att-player">${playerOptions()}</select></label>
    <label>Guests <input type="number" class="att-guests" value="0" min="0" style="width:60px"></label>
    <button class="btn btn-sm btn-danger" onclick="this.closest('.attendee-row').remove()">✕</button>
  `;
  container.appendChild(row);
}

function toggleNewMatchForm() {
  const el = document.getElementById('newMatchForm');
  el.style.display = el.style.display === 'none' ? 'block' : 'none';
}

async function refreshMatchList() {
  const el = document.getElementById('matchList');
  if (!el) return;
  const res = await api('GET', '/matches?limit=100');
  if (!res) return;
  const list = await res.json();
  if (!list.length) { el.innerHTML = '<p class="empty">No matches yet.</p>'; return; }
  el.innerHTML = list.map(m => `
    <div class="match-card card" id="mc-${m.id}">
      <div class="match-header" onclick="toggleMatchDetail(${m.id})">
        <div>
          <span class="match-id">#${m.id}</span>
          <span class="match-date">${m.date}</span>
          ${m.notes ? `<span class="match-notes">(${escapeHTML(m.notes)})</span>` : ''}
        </div>
        <div style="display:flex;align-items:center;gap:16px">
          <span class="match-bill">৳${m.total_bill.toFixed(2)}</span>
          <div class="match-actions">
            ${state.token ? `
            <button class="btn btn-sm"
              data-match-id="${m.id}"
              data-match-date="${m.date}"
              data-match-bill="${m.total_bill}"
              data-match-notes="${escapeHTML(m.notes || '')}"
              onclick="event.stopPropagation();showEditMatchBtn(this)">Edit</button>
            <button class="btn btn-sm btn-danger" onclick="event.stopPropagation();deleteMatch(${m.id})">Delete</button>
            ` : ''}
            <span class="toggle-icon" id="ti-${m.id}">▼</span>
          </div>
        </div>
      </div>
      <div id="md-${m.id}" class="match-detail" style="display:none"></div>
    </div>
  `).join('');
}

async function toggleMatchDetail(id) {
  const detail = document.getElementById('md-' + id);
  const icon = document.getElementById('ti-' + id);
  if (detail.style.display !== 'none') {
    detail.style.display = 'none';
    icon.textContent = '▼';
    return;
  }
  detail.style.display = 'block';
  icon.textContent = '▲';
  if (detail.dataset.loaded) return;
  detail.innerHTML = '<p class="loading">Loading…</p>';
  const res = await api('GET', '/matches/' + id);
  if (!res) return;
  const md = await res.json();
  detail.dataset.loaded = '1';
  const attendeeRows = (md.attendees || []).map(a =>
    `<tr><td>${escapeHTML(a.player_name)}</td><td>${a.guest_count}</td><td>৳${a.due.toFixed(2)}</td></tr>`
  );
  const paymentRows = (md.payments || []).map(p =>
    `<tr><td>${escapeHTML(p.player_name)}</td><td>৳${p.amount.toFixed(2)}</td><td>${p.paid_at}</td><td>${escapeHTML(p.notes || '—')}</td></tr>`
  );
  detail.innerHTML = `
    <div class="detail-grid">
      <div>
        <h4>Attendees</h4>
        ${renderTable(['Player','Guests','Due'], attendeeRows)}
      </div>
      <div>
        <h4>Payments</h4>
        ${renderTable(['Player','Amount','Date','Notes'], paymentRows)}
      </div>
    </div>
  `;
}

function showEditMatchBtn(btn) {
  const id = parseInt(btn.dataset.matchId);
  const date = btn.dataset.matchDate;
  const bill = parseFloat(btn.dataset.matchBill);
  const notes = btn.dataset.matchNotes;
  showEditMatch(id, date, bill, notes);
}

async function showEditMatch(id, date, bill, notes) {
  const card = document.getElementById('mc-' + id);
  const existing = card.querySelector('.edit-form');
  if (existing) { existing.remove(); return; }

  const res = await api('GET', '/matches/' + id);
  const md = res ? await res.json() : {};
  const attendees = md.attendees || [];

  const form = document.createElement('div');
  form.className = 'edit-form card';
  form.innerHTML = `
    <h4>Edit Match #${id}</h4>
    <div class="form-row">
      <label>Date <input type="date" id="em-date-${id}" value="${date}"></label>
      <label>Bill <input type="number" id="em-bill-${id}" step="0.01" value="${bill}"></label>
      <label>Notes <input type="text" id="em-notes-${id}" value="${escapeHTML(notes)}"></label>
    </div>
    <div style="margin-top:12px">
      <div style="display:flex;align-items:center;gap:8px;margin-bottom:6px">
        <strong>Attendees</strong>
        <button class="btn btn-sm" onclick="addAttendeeRow('em-att-${id}')">+ Add</button>
      </div>
      <div id="em-att-${id}">
        ${attendees.map(a => `
          <div class="form-row attendee-row" style="align-items:center">
            <label>Player <select class="att-player">${playerOptionsSelected(a.player_id)}</select></label>
            <label>Guests <input type="number" class="att-guests" value="${a.guest_count}" min="0" style="width:60px"></label>
            <button class="btn btn-sm btn-danger" onclick="this.closest('.attendee-row').remove()">✕</button>
          </div>
        `).join('')}
      </div>
    </div>
    <div style="display:flex;gap:8px;margin-top:12px">
      <button class="btn btn-primary" onclick="updateMatch(${id})">Save</button>
      <button class="btn" onclick="this.closest('.edit-form').remove()">Cancel</button>
    </div>
  `;
  form.querySelector('#em-att-' + id).dataset.originals = JSON.stringify(attendees.map(a => a.player_id));
  card.appendChild(form);
}

async function createMatch() {
  const date = document.getElementById('mDate').value;
  const bill = parseFloat(document.getElementById('mBill').value);
  const notes = document.getElementById('mNotes').value || null;
  if (!date || !bill) { toast('Date and bill required', false); return; }
  const attendees = Array.from(document.querySelectorAll('.attendee-row')).map(row => ({
    player_id: parseInt(row.querySelector('.att-player').value),
    guest_count: parseInt(row.querySelector('.att-guests').value) || 0,
  })).filter(a => a.player_id);
  const res = await api('POST', '/matches', { date, total_bill: bill, notes, attendees });
  if (!res) return;
  const data = await res.json();
  if (!res.ok) { toast(data.error || 'Error', false); return; }
  toast('Match #' + data.id + ' created');
  toggleNewMatchForm();
  await refreshMatchList();
}

async function updateMatch(id) {
  const body = {};
  const d = document.getElementById('em-date-' + id)?.value;
  const b = document.getElementById('em-bill-' + id)?.value;
  const n = document.getElementById('em-notes-' + id)?.value;
  if (d) body.date = d;
  if (b) body.total_bill = parseFloat(b);
  if (n !== undefined) body.notes = n || null;
  const res = await api('PATCH', '/matches/' + id, body);
  if (!res) return;
  if (!res.ok) { const errData = await res.json(); toast(errData.error || 'Error', false); return; }

  const container = document.getElementById('em-att-' + id);
  if (container) {
    const originalIds = JSON.parse(container.dataset.originals || '[]');
    const rows = Array.from(container.querySelectorAll('.attendee-row'));
    const currentIds = [];
    for (const row of rows) {
      const playerId = parseInt(row.querySelector('.att-player').value);
      const guestCount = parseInt(row.querySelector('.att-guests').value) || 0;
      if (!playerId) continue;
      currentIds.push(playerId);
      await api('POST', '/matches/' + id + '/attendees', { player_id: playerId, guest_count: guestCount });
    }
    for (const pid of originalIds) {
      if (!currentIds.includes(pid)) {
        await api('DELETE', '/matches/' + id + '/attendees/' + pid);
      }
    }
  }

  toast('Match updated');
  document.querySelector('#mc-' + id + ' .edit-form')?.remove();
  const detail = document.getElementById('md-' + id);
  if (detail) { delete detail.dataset.loaded; detail.innerHTML = ''; detail.style.display = 'none'; }
  await refreshMatchList();
}

async function deleteMatch(id) {
  if (!confirm('Delete match #' + id + '?')) return;
  const res = await api('DELETE', '/matches/' + id);
  if (!res) return;
  if (!res.ok) { toast('Error deleting match', false); return; }
  toast('Match deleted');
  await refreshMatchList();
}
