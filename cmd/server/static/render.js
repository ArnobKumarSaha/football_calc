'use strict';

function renderApp() {
  const root = document.getElementById('app');
  if (!state.token) {
    root.innerHTML = loginHTML();
    return;
  }
  root.innerHTML = shellHTML();
  const page = currentPage();
  setActiveNav(page);
  loadPage(page);
}

function loginHTML() {
  return `
    <div class="login-wrap">
      <div class="login-card">
        <div class="login-logo">⚽ Football Calc</div>
        <p class="login-sub">Admin Login</p>
        <input type="password" id="pwd" class="login-input" placeholder="Admin password"
               onkeydown="if(event.key==='Enter') login()">
        <button class="btn btn-primary w-full" onclick="login()">Login</button>
      </div>
    </div>
  `;
}

function shellHTML() {
  const navItems = [
    { page: 'matches',    label: 'Matches' },
    { page: 'players',    label: 'Players' },
    { page: 'attendance', label: 'Attendance' },
    { page: 'payments',   label: 'Payments' },
    { page: 'balances',   label: 'Balances' },
    { page: 'reports',    label: 'Reports' },
  ];
  return `
    <div class="app-layout">
      <aside class="sidebar">
        <div class="sidebar-logo">⚽ Football Calc</div>
        <nav class="sidebar-nav">
          ${navItems.map(n =>
            `<a href="#${n.page}" class="nav-item" data-page="${n.page}">${n.label}</a>`
          ).join('')}
        </nav>
        <button class="logout-btn" onclick="logout()">Logout</button>
      </aside>
      <main class="content">
        <div id="page"></div>
      </main>
    </div>
  `;
}

renderApp();
