'use strict';

function login() {
  const pwd = document.getElementById('pwd').value.trim();
  if (!pwd) return;
  state.token = pwd;
  localStorage.setItem('token', pwd);
  renderApp();
}

function logout() {
  state.token = '';
  state.players = [];
  localStorage.removeItem('token');
  renderApp();
}
