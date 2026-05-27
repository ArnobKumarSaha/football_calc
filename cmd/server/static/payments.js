'use strict';

async function payments(el) {
  await ensurePlayers();
  el.innerHTML = `
    <div class="page-header"><h2>Payments</h2></div>
    <div class="card form-section">
      <h3>Record Payment</h3>
      <div class="form-row">
        <label>Player <select id="payPlayer">${playerOptions()}</select></label>
        <label>Amount <input type="number" id="payAmount" step="0.01" min="0" placeholder="0.00"></label>
        <label>Paid At <input type="date" id="payDate"></label>
        <label>Match ID (opt) <input type="number" id="payMatchId" placeholder="leave blank"></label>
        <label>Notes <input type="text" id="payNotes" placeholder="optional"></label>
      </div>
      <button class="btn btn-primary" onclick="createPayment()">Record</button>
    </div>
    <div class="card">
      <h3>Payment History</h3>
      <div class="form-row">
        <label>Player <select id="filterPlayer">${playerOptions(true)}</select></label>
        <button class="btn" style="align-self:flex-end" onclick="listPayments()">Load</button>
      </div>
      <div id="payList"></div>
    </div>
  `;
}

async function createPayment() {
  const playerId = parseInt(document.getElementById('payPlayer').value);
  const amount = parseFloat(document.getElementById('payAmount').value);
  const paidAt = document.getElementById('payDate').value;
  const matchIdRaw = document.getElementById('payMatchId').value;
  const notes = document.getElementById('payNotes').value || null;
  if (!playerId || !amount || !paidAt) { toast('Player, amount and date required', false); return; }
  const body = { player_id: playerId, amount, paid_at: paidAt, notes };
  if (matchIdRaw) body.match_id = parseInt(matchIdRaw);
  const res = await api('POST', '/payments', body);
  if (!res) return;
  const data = await res.json();
  if (!res.ok) { toast(data.error || 'Error', false); return; }
  toast('Payment #' + data.id + ' recorded');
  document.getElementById('payAmount').value = '';
  document.getElementById('payNotes').value = '';
  document.getElementById('payMatchId').value = '';
}

async function listPayments() {
  const playerFilter = document.getElementById('filterPlayer').value;
  let path = '/payments?';
  if (playerFilter) path += 'player_id=' + playerFilter + '&';
  const res = await api('GET', path);
  if (!res) return;
  const data = await res.json();
  const el = document.getElementById('payList');
  if (!data.length) { el.innerHTML = '<p class="empty">No payments found.</p>'; return; }
  const rows = data.map(p => `
    <tr>
      <td>${p.id}</td>
      <td>${p.player_id}</td>
      <td>${p.match_id || '—'}</td>
      <td>৳${p.amount.toFixed(2)}</td>
      <td>${p.paid_at}</td>
      <td>${p.notes || '—'}</td>
    </tr>
  `);
  el.innerHTML = renderTable(['ID','Player','Match','Amount','Date','Notes'], rows);
}

async function deletePayment(id) {
  if (!confirm('Delete payment?')) return;
  const res = await api('DELETE', '/payments/' + id);
  if (!res) return;
  if (!res.ok) { toast('Error', false); return; }
  toast('Deleted');
  listPayments();
}
