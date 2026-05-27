'use strict';

async function admin(el) {
  el.innerHTML = `
    <div class="page-header"><h2>Admin</h2></div>
    <div class="card form-section">
      <h3>Danger Zone</h3>
      <p>This will permanently delete all matches, attendees, payments, and players.</p>
      <button class="btn btn-danger" onclick="cleanupAllData()">Clear All Data</button>
    </div>
  `;
}

async function cleanupAllData() {
  if (!confirm('This will delete ALL data (players, matches, attendees, payments). Are you sure?')) return;
  if (!confirm('Final confirmation: delete everything?')) return;
  const res = await api('POST', '/admin/cleanup');
  if (!res) return;
  const data = await res.json();
  if (!res.ok) { toast(data.error || 'Error', false); return; }
  toast('All data cleared');
}
