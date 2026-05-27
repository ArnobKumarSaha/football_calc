'use strict';

async function admin(el) {
  el.innerHTML = `
    <div class="page-header"><h2>Admin</h2></div>
    <div class="card form-section">
      <h3>Import Initial Data</h3>
      <p>Upload a text file with one entry per line: <code>&lt;name&gt; &lt;amount&gt;</code><br>
         Negative amount means the player is in debt. Only allowed when the database is empty.</p>
      <div style="display:flex;gap:10px;align-items:center;flex-wrap:wrap">
        <input type="file" id="import-file" accept=".txt,text/plain">
        <button class="btn btn-primary" onclick="importData()">Import</button>
      </div>
    </div>
    <div class="card form-section" style="margin-top:16px">
      <h3>Danger Zone</h3>
      <p>This will permanently delete all matches, attendees, payments, and players.</p>
      <button class="btn btn-danger" onclick="cleanupAllData()">Clear All Data</button>
    </div>
  `;
}

async function importData() {
  const fileInput = document.getElementById('import-file');
  if (!fileInput || !fileInput.files.length) {
    toast('Select a file first', false);
    return;
  }
  const formData = new FormData();
  formData.append('file', fileInput.files[0]);

  const res = await fetch('/api/admin/import', {
    method: 'POST',
    headers: { Authorization: 'Bearer ' + state.token },
    body: formData,
  });
  const data = await res.json();
  if (!res.ok) { toast(data.error || 'Import failed', false); return; }
  toast(`Imported ${data.imported} player(s)`);
  fileInput.value = '';
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
