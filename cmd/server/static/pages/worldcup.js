'use strict';

let _wcToken = 0;
const _wcCache = { fixtures: null, standings: null };

function worldcup(el) {
  const token = ++_wcToken;

  el.innerHTML = `
    <div class="page-header">
      <h2>World Cup 2026</h2>
      <div class="wc-search-wrap">
        <svg class="wc-search-icon" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
          <circle cx="8.5" cy="8.5" r="5.5" stroke="currentColor" stroke-width="1.7"/>
          <line x1="12.5" y1="12.5" x2="17" y2="17" stroke="currentColor" stroke-width="1.7" stroke-linecap="round"/>
        </svg>
        <input id="wc-search" type="text" class="wc-search-input" placeholder="Search by team…" oninput="wcFilterFixtures(this.value)">
      </div>
    </div>
    <div class="wc-tabs">
      <button class="wc-tab active" onclick="wcSwitchTab(this,'wc-fixtures')">Fixtures</button>
      <button class="wc-tab" onclick="wcSwitchTab(this,'wc-standings')">Standings</button>
    </div>
    <div id="wc-fixtures">
      <div id="wc-fixture-list"><p class="loading">Loading…</p></div>
    </div>
    <div id="wc-standings" style="display:none"><p class="loading">Loading…</p></div>
  `;

  wcLoadFixtures(token);
  wcLoadStandings(token);
}

function wcIsStale(token) {
  return token !== _wcToken;
}

function wcSwitchTab(btn, target) {
  document.querySelectorAll('.wc-tab').forEach(b => b.classList.remove('active'));
  btn.classList.add('active');
  ['wc-fixtures', 'wc-standings'].forEach(id => {
    const el = document.getElementById(id);
    if (el) el.style.display = id === target ? '' : 'none';
  });
}

function wcFilterFixtures(q) {
  const fixtures = _wcCache.fixtures;
  if (!fixtures) return;
  const el = document.getElementById('wc-fixture-list');
  if (!el) return;
  wcRenderFixtures(el, q.trim().toLowerCase());
}

function wcRenderFixtures(el, filter) {
  const fixtures = _wcCache.fixtures || [];
  const filtered = filter
    ? fixtures.filter(f =>
        f.home_team.name.toLowerCase().includes(filter) ||
        f.away_team.name.toLowerCase().includes(filter) ||
        f.home_team.abbr.toLowerCase().includes(filter) ||
        f.away_team.abbr.toLowerCase().includes(filter)
      )
    : fixtures;

  if (!filtered.length) {
    el.innerHTML = '<p class="empty">No matches found.</p>';
    return;
  }

  const byDate = {};
  filtered.forEach(f => {
    const d = new Date(f.date).toLocaleDateString('en-US', {
      weekday: 'short', month: 'short', day: 'numeric', timeZone: 'UTC',
    });
    (byDate[d] = byDate[d] || []).push(f);
  });

  el.innerHTML = Object.entries(byDate).map(([date, matches]) => `
    <div class="wc-date-group">
      <div class="wc-date-label">${date}</div>
      ${matches.map(wcFixtureCard).join('')}
    </div>
  `).join('');
}

async function wcLoadFixtures(token) {
  try {
    if (!_wcCache.fixtures) {
      const res = await fetch('/api/worldcup/fixtures');
      if (wcIsStale(token)) return;
      if (!res.ok) throw new Error('HTTP ' + res.status);
      const data = await res.json();
      if (wcIsStale(token)) return;
      _wcCache.fixtures = data;
    }

    const el = document.getElementById('wc-fixture-list');
    if (!el || wcIsStale(token)) return;
    wcRenderFixtures(el, '');
  } catch (e) {
    if (wcIsStale(token)) return;
    const el = document.getElementById('wc-fixture-list');
    if (el) el.innerHTML = '<p class="empty">Failed to load fixtures.</p>';
  }
}

function wcFixtureCard(f) {
  const isLive = f.status === 'in';
  const isPre  = f.status === 'pre';
  const statusLabel = isLive
    ? `<span class="wc-status live">LIVE</span>`
    : isPre
      ? `<span class="wc-status pre">${f.detail}</span>`
      : `<span class="wc-status post">FT</span>`;
  const score = isPre
    ? `<div class="wc-scores"><span class="wc-vs">vs</span></div>`
    : `<div class="wc-scores"><span>${f.home_score}</span><span class="wc-dash">–</span><span>${f.away_score}</span></div>`;

  return `
    <div class="wc-fixture card">
      <div class="wc-team home">
        <img class="wc-logo" src="${f.home_team.logo}" alt="${f.home_team.abbr}" onerror="this.style.display='none'">
        <span class="wc-team-name">${f.home_team.name}</span>
      </div>
      <div class="wc-middle">
        ${statusLabel}
        ${score}
        ${f.venue ? `<div class="wc-venue">${f.venue}</div>` : ''}
      </div>
      <div class="wc-team away">
        <img class="wc-logo" src="${f.away_team.logo}" alt="${f.away_team.abbr}" onerror="this.style.display='none'">
        <span class="wc-team-name">${f.away_team.name}</span>
      </div>
    </div>
  `;
}

async function wcLoadStandings(token) {
  try {
    if (!_wcCache.standings) {
      const res = await fetch('/api/worldcup/standings');
      if (wcIsStale(token)) return;
      if (!res.ok) throw new Error('HTTP ' + res.status);
      const data = await res.json();
      if (wcIsStale(token)) return;
      _wcCache.standings = data;
    }

    const el = document.getElementById('wc-standings');
    if (!el || wcIsStale(token)) return;

    const groups = _wcCache.standings;
    if (!groups.length) {
      el.innerHTML = '<p class="empty">No standings.</p>';
      return;
    }
    el.innerHTML = `<div class="wc-groups">${groups.map(wcGroupTable).join('')}</div>`;
  } catch (e) {
    if (wcIsStale(token)) return;
    const el = document.getElementById('wc-standings');
    if (el) el.innerHTML = '<p class="empty">Failed to load standings.</p>';
  }
}

function wcGroupTable(g) {
  const rows = (g.entries || []).map(e => `
    <tr>
      <td>
        <img class="wc-logo-sm" src="${e.team.logo}" alt="${e.team.abbr}" onerror="this.style.display='none'">
        ${e.team.name}
      </td>
      <td>${e.gp}</td><td>${e.w}</td><td>${e.d}</td><td>${e.l}</td>
      <td>${e.gf}</td><td>${e.ga}</td><td>${e.gd}</td>
      <td><strong>${e.pts}</strong></td>
    </tr>
  `).join('');

  return `
    <div class="card wc-group-card">
      <h4 class="wc-group-title">${g.name}</h4>
      <table class="wc-table">
        <thead><tr><th>Team</th><th>GP</th><th>W</th><th>D</th><th>L</th><th>GF</th><th>GA</th><th>GD</th><th>Pts</th></tr></thead>
        <tbody>${rows}</tbody>
      </table>
    </div>
  `;
}
