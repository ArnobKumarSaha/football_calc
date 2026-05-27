'use strict';

async function payments(el) {
  if (!state.token) {
    el.innerHTML = `
      <div class="page-header"><h2>Record Payment</h2></div>
      <div class="card form-section">
        <p>Admin login required to record payments.</p>
        <button class="btn btn-primary" onclick="openLoginModal()">Login</button>
      </div>
    `;
    return;
  }

  await ensurePlayers();
  const mRes = await api('GET', '/matches?limit=200');
  const matchList = mRes ? await mRes.json() : [];
  const matchOptions = '<option value="">— no match —</option>' +
    matchList.map(m => `<option value="${m.id}">${m.date}${m.notes ? ' — ' + m.notes : ''}</option>`).join('');

  el.innerHTML = `
    <div class="page-header"><h2>Record Payment</h2></div>
    <div class="card form-section">
      <div class="form-row">
        <label>Player <select id="payPlayer">${playerOptions()}</select></label>
        <label>Amount <input type="number" id="payAmount" step="0.01" min="0" placeholder="0.00"></label>
        <label>Date <input type="date" id="payDate" value="${new Date().toISOString().slice(0,10)}"></label>
        <label>Match (optional) <select id="payMatch">${matchOptions}</select></label>
        <label>Notes <input type="text" id="payNotes" placeholder="optional"></label>
      </div>
      <button class="btn btn-primary" style="margin-top:8px" onclick="submitPayment()">Record</button>
    </div>
  `;
}

async function submitPayment() {
  const playerId = parseInt(document.getElementById('payPlayer').value);
  const amount = parseFloat(document.getElementById('payAmount').value);
  const paidAt = document.getElementById('payDate').value;
  const matchRaw = document.getElementById('payMatch').value;
  const notes = document.getElementById('payNotes').value || null;
  if (!playerId || !amount || !paidAt) { toast('Player, amount and date required', false); return; }
  const body = { player_id: playerId, amount, paid_at: paidAt, notes };
  if (matchRaw) body.match_id = parseInt(matchRaw);
  const res = await api('POST', '/payments', body);
  if (!res) return;
  const data = await res.json();
  if (!res.ok) { toast(data.error || 'Error', false); return; }
  toast('Payment #' + data.id + ' recorded');
  document.getElementById('payAmount').value = '';
  document.getElementById('payNotes').value = '';
  document.getElementById('payMatch').value = '';
}
