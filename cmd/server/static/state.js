'use strict';

const basePath = location.pathname.replace(/[^/]*$/, '');

const state = {
  token: localStorage.getItem('token') || '',
  players: [],
};
