'use strict';

async function login() {
  const pwd = document.getElementById('pwd').value.trim();
  if (!pwd) return;
  state.token = pwd;
  const res = await api('GET', '/auth');
  if (!res) { state.token = ''; return; }
  localStorage.setItem('token', pwd);
  closeLoginModal();
  renderApp();
}

function logout() {
  state.token = '';
  state.players = [];
  localStorage.removeItem('token');
  renderApp();
}
