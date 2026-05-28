'use strict';

function renderApp() {
  const root = document.getElementById('app');
  root.innerHTML = shellHTML();
  const page = currentPage();
  setActiveNav(page);
  loadPage(page);
}

function loginModalHTML() {
  return `
    <div class="modal-overlay" id="loginModal" onclick="if(event.target===this)closeLoginModal()">
      <div class="login-card">
        <div class="login-logo">⚽ Football Calc</div>
        <p class="login-sub">Admin Login</p>
        <input type="password" id="pwd" class="login-input" placeholder="Admin password"
               onkeydown="if(event.key==='Enter') login()">
        <button class="btn btn-primary w-full" onclick="login()">Login</button>
        <button class="btn w-full" style="margin-top:8px" onclick="closeLoginModal()">Cancel</button>
      </div>
    </div>
  `;
}

function openLoginModal() {
  if (document.getElementById('loginModal')) return;
  document.body.insertAdjacentHTML('beforeend', loginModalHTML());
  document.getElementById('pwd').focus();
}

function closeLoginModal() {
  document.getElementById('loginModal')?.remove();
}

function shellHTML() {
  const mainNav = [
    { page: 'matches',   label: 'Matches' },
    { page: 'payments',  label: 'Payments' },
    { page: 'balances',  label: 'Balances' },
    { page: 'worldcup',  label: 'World Cup 2026' },
  ];
  const bottomNav = [
    { page: 'reports', label: 'Reports' },
    ...(state.token ? [{ page: 'admin', label: 'Admin' }] : []),
  ];
  const authBtn = state.token
    ? `<button class="logout-btn" onclick="logout()">Logout</button>`
    : `<button class="logout-btn" onclick="openLoginModal()">Login</button>`;
  return `
    <div class="app-layout">
      <aside class="sidebar">
        <div class="sidebar-logo">⚽ Football Calc</div>
        <nav class="sidebar-nav">
          ${mainNav.map(n =>
            `<a href="#${n.page}" class="nav-item" data-page="${n.page}">${n.label}</a>`
          ).join('')}
        </nav>
        <nav class="sidebar-nav sidebar-nav-bottom">
          ${bottomNav.map(n =>
            `<a href="#${n.page}" class="nav-item" data-page="${n.page}">${n.label}</a>`
          ).join('')}
          ${authBtn}
        </nav>
      </aside>
      <main class="content">
        <div id="page"></div>
      </main>
    </div>
    <div class="sidebar-overlay" id="sidebarOverlay" onclick="closeSidebar()"></div>
    <button class="hamburger" id="hamburgerBtn" onclick="toggleSidebar()">&#9776;</button>
  `;
}

renderApp();
