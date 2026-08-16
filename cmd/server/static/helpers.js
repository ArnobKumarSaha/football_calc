'use strict';

function escapeHTML(s) {
  return String(s)
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;');
}

function renderTable(headers, rows) {
  if (!rows.length) return '<p class="empty">No data.</p>';
  return `
    <table>
      <thead><tr>${headers.map(h => `<th>${h}</th>`).join('')}</tr></thead>
      <tbody>${rows.join('')}</tbody>
    </table>
  `;
}

function balanceBadge(v, negThreshold = 0) {
  const cls = v < negThreshold ? 'neg' : 'pos';
  return `<span class="balance ${cls}">${v >= 0 ? '+' : ''}${v.toFixed(2)}</span>`;
}

async function ensurePlayers() {
  if (state.players.length) return;
  const res = await api('GET', '/players');
  if (!res) return;
  state.players = await res.json();
}

function playerOptions(includeAll = false) {
  const all = includeAll ? '<option value="">All players</option>' : '<option value="">— select —</option>';
  return all + state.players.map(p => `<option value="${p.id}">${p.name}</option>`).join('');
}

function playerOptionsSelected(selectedId) {
  return '<option value="">— select —</option>' +
    state.players.map(p => `<option value="${p.id}" ${p.id === selectedId ? 'selected' : ''}>${p.name}</option>`).join('');
}
