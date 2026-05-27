'use strict';

async function api(method, path, body) {
  const opts = { method, headers: { 'Content-Type': 'application/json' } };
  if (state.token) opts.headers['Authorization'] = 'Bearer ' + state.token;
  if (body !== undefined) opts.body = JSON.stringify(body);
  let res;
  try {
    res = await fetch('/api' + path, opts);
  } catch (e) {
    toast('Network error', false);
    return null;
  }
  if (res.status === 401) {
    toast('Wrong password or session expired', false);
    logout();
    return null;
  }
  return res;
}
