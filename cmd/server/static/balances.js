'use strict';

async function balances(el) {
  const res = await api('GET', '/players');
  if (!res) return;
  const list = await res.json();
  if (!list.length) { el.innerHTML = '<p class="empty">No players.</p>'; return; }
  el.innerHTML = `
    <div class="page-header"><h2>Balances</h2></div>
    <div class="balances-grid">
      ${list.map(p => `
        <div class="balance-card card" data-id="${p.id}" onclick="loadPlayerDetail(this)">
          <div class="balance-name">${p.name}</div>
          <div class="balance-amount">${balanceBadge(p.balance)}</div>
          <div class="balance-hint">click for history</div>
        </div>
      `).join('')}
    </div>
    <div id="playerDetail"></div>
    <div class="card bank-info">
      <h4>Bank Account Info</h4>
      <p>Account Number: 1056069500001</p>
      <p>Account Name: ARNOB KUMAR SAHA</p>
      <p>Bank Name: BRAC Bank PLC</p>
      <p>Branch Name: MADANI AVENUE BRANCH</p>
      <p>Routing Number: 060263429</p>
      <p>SWIFT Code: BRAKBDDH</p>
    </div>
  `;
}

async function loadPlayerDetail(card) {
  const id = card.dataset.id;
  const name = card.querySelector('.balance-name').textContent;
  const el = document.getElementById('playerDetail');
  el.innerHTML = '<p class="loading">Loading…</p>';
  const [balRes, payRes] = await Promise.all([
    api('GET', '/players/' + id + '/balance'),
    api('GET', '/payments?player_id=' + id),
  ]);
  if (!balRes) return;
  const h = await balRes.json();
  const allPayments = payRes ? await payRes.json() : [];

  const matchRows = (h.matches || []).map(m =>
    `<tr><td>${m.match_id}</td><td>${m.match_date}</td><td>৳${m.due != null ? m.due.toFixed(2) : '—'}</td><td>৳${m.paid.toFixed(2)}</td></tr>`
  );
  const payRows = allPayments.map(p =>
    `<tr><td>${p.paid_at}</td><td>৳${p.amount.toFixed(2)}</td><td>${p.match_id || '—'}</td></tr>`
  );
  el.innerHTML = `
    <div class="card">
      <div class="detail-header">
        <h3>${name}</h3>
        <div class="detail-stats">
          <span>Due: ৳${h.total_due.toFixed(2)}</span>
          <span>Paid: ৳${h.total_paid.toFixed(2)}</span>
          <span>Balance: ${balanceBadge(h.balance)}</span>
        </div>
      </div>
      <div class="detail-grid">
        <div>
          <h4 style="margin-bottom:8px">Match History</h4>
          ${matchRows.length ? renderTable(['Match', 'Date', 'Due', 'Paid'], matchRows) : '<p class="empty">No matches.</p>'}
        </div>
        <div>
          <h4 style="margin-bottom:8px">Payments</h4>
          ${payRows.length ? renderTable(['Date', 'Amount', 'Match'], payRows) : '<p class="empty">No payments.</p>'}
        </div>
      </div>
    </div>
  `;
}
