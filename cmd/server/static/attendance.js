'use strict';

async function attendance(el) {
  if (!state.token) {
    el.innerHTML = `
      <div class="page-header"><h2>Attendance</h2></div>
      <div class="card form-section">
        <p>Admin login required to manage attendance.</p>
        <button class="btn btn-primary" onclick="openLoginModal()">Login</button>
      </div>
    `;
    return;
  }
  await ensurePlayers();
  const mRes = await api('GET', '/matches?limit=200');
  const matchList = mRes ? await mRes.json() : [];
  const matchOptions = matchList.map(m =>
    `<option value="${m.id}">${m.date}${m.notes ? ' — ' + m.notes : ''}</option>`
  ).join('');

  el.innerHTML = `
    <div class="page-header"><h2>Attendance</h2></div>
    <div class="card form-section">
      <h3>Add / Update Attendee</h3>
      <div class="form-row">
        <label>Match <select id="attMatchId"><option value="">— select match —</option>${matchOptions}</select></label>
        <label>Player <select id="attPlayerId">${playerOptions()}</select></label>
        <label>Guests <input type="number" id="attGuest" value="0" min="0"></label>
      </div>
      <div style="display:flex;gap:8px">
        <button class="btn btn-primary" onclick="upsertAttendee()">Add / Update</button>
        <button class="btn btn-danger" onclick="removeAttendee()">Remove</button>
      </div>
    </div>
  `;
}

async function upsertAttendee() {
  const matchId = document.getElementById('attMatchId').value;
  const playerId = document.getElementById('attPlayerId').value;
  const guestCount = parseInt(document.getElementById('attGuest').value) || 0;
  if (!matchId || !playerId) { toast('Match and player required', false); return; }
  const res = await api('POST', '/matches/' + matchId + '/attendees', {
    player_id: parseInt(playerId), guest_count: guestCount,
  });
  if (!res) return;
  const data = await res.json();
  if (!res.ok) { toast(data.error || 'Error', false); return; }
  toast('Attendee saved');
}

async function removeAttendee() {
  const matchId = document.getElementById('attMatchId').value;
  const playerId = document.getElementById('attPlayerId').value;
  if (!matchId || !playerId) { toast('Match and player required', false); return; }
  const res = await api('DELETE', '/matches/' + matchId + '/attendees/' + playerId);
  if (!res) return;
  toast(res.ok ? 'Removed' : 'Error', res.ok);
}
