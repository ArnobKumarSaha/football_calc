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
  `;
}

async function loadPlayerDetail(card) {
  const id = card.dataset.id;
  const name = card.querySelector('.balance-name').textContent;
  const el = document.getElementById('playerDetail');
  el.innerHTML = '<p class="loading">Loading…</p>';
  const res = await api('GET', '/players/' + id + '/balance');
  if (!res) return;
  const h = await res.json();
  const rows = (h.matches || []).map(m =>
    `<tr><td>${m.match_date}</td><td>৳${m.due != null ? m.due.toFixed(2) : '—'}</td><td>৳${m.paid.toFixed(2)}</td></tr>`
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
      ${renderTable(['Date', 'Due', 'Paid'], rows)}
    </div>
  `;
}
