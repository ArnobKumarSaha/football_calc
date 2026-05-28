'use strict';

async function reports(el) {
  el.innerHTML = `
    <div class="page-header"><h2>Reports</h2></div>
    <div class="reports-grid">
      <div class="card">
        <h3>Monthly Summary</h3>
        <div class="form-row">
          <label>Year <input type="number" id="reportYear" value="${new Date().getFullYear()}"></label>
          <button class="btn btn-primary" style="align-self:flex-end" onclick="loadMonthly()">Load</button>
        </div>
        <div id="monthlyTable"></div>
      </div>
      <div class="card">
        <h3>CSV Exports</h3>
        <div class="csv-links">
          <a href="/api/reports/matches.csv" class="btn btn-primary" download>Matches CSV</a>
          <a href="/api/reports/payments.csv" class="btn btn-primary" download>Payments CSV</a>
        </div>
      </div>
    </div>
  `;
}

async function loadMonthly() {
  const year = document.getElementById('reportYear').value;
  const res = await api('GET', '/reports/monthly?year=' + year);
  if (!res) return;
  const data = await res.json();
  const rows = data.map(r =>
    `<tr><td>${r.month}</td><td>${r.match_count}</td><td>৳${r.total_spent.toFixed(2)}</td></tr>`
  );
  document.getElementById('monthlyTable').innerHTML = renderTable(['Month','Matches','Total Spent'], rows);
}
