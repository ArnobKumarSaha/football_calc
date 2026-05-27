'use strict';

const routes = { matches, players, attendance, payments, balances, reports, admin };

window.addEventListener('hashchange', () => {
  const page = currentPage();
  setActiveNav(page);
  loadPage(page);
});

function currentPage() {
  return location.hash.slice(1) || 'balances';
}

function setActiveNav(page) {
  document.querySelectorAll('.nav-item').forEach(el => {
    el.classList.toggle('active', el.dataset.page === page);
  });
}

function loadPage(page) {
  const el = document.getElementById('page');
  if (!el) return;
  el.innerHTML = '<p class="loading">Loading…</p>';
  const fn = routes[page];
  if (fn) fn(el);
  else el.innerHTML = '<p>Page not found.</p>';
}
