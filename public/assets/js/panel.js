const cpuText = document.getElementById("cpu-usage");
const cpuBar = document.getElementById("cpu-progress");
const ramText = document.getElementById("ram-usage");
const ramBar = document.getElementById("ram-progress");
const userText = document.getElementById("user-online");
const userBar = document.getElementById("user-progress");
const netText = document.getElementById("network-speed");
const netBar = document.getElementById("network-progress");
const timeDisplay = document.getElementById("current-time");
const logContainer = document.getElementById("log-table-body");
const alertContainer = document.getElementById("alert-container");

function showAlert(message, type = "success") {
  if (!alertContainer) return;
  alertContainer.innerHTML = "";
  alertContainer.innerHTML = `
        <div class="alert alert-${type} alert-dismissible fade show shadow-sm" role="alert">
            <div>${message}</div>
            <button type="button" class="btn-close" data-bs-dismiss="alert" onclick="clearAlertContainer()"></button>
        </div>
    `;
}

function clearAlertContainer() {
  if (alertContainer) alertContainer.innerHTML = "";
}

function toggleSidebar() {
  const sidebar = document.getElementById("sidebarMenu");
  if (sidebar) {
    sidebar.classList.toggle("show");
  }
}

function updateClock() {
  const time = new Date();
  if (timeDisplay) {
    timeDisplay.innerText = time.toLocaleTimeString("fa-IR");
  }
}
setInterval(updateClock, 1000);
updateClock();

function processMetrics() {
  const cpuVal = Math.floor(Math.random() * 50) + 20;
  const ramVal = Math.floor(Math.random() * 30) + 50;
  const userVal = Math.floor(Math.random() * 200) + 100;
  const netVal = Math.floor(Math.random() * 100) + 50;

  if (cpuText && cpuBar) {
    cpuText.innerText = `${cpuVal}%`;
    cpuBar.style.width = `${cpuVal}%`;
  }
  if (ramText && ramBar) {
    ramText.innerText = `${ramVal}%`;
    ramBar.style.width = `${ramVal}%`;
  }
  if (userText && userBar) {
    userText.innerText = userVal;
    userBar.style.width = `${(userVal / 300) * 100}%`;
  }
  if (netText && netBar) {
    netText.innerText = `${netVal} Mbps`;
    netBar.style.width = `${(netVal / 150) * 100}%`;
  }
}

if (cpuText) {
  setInterval(processMetrics, 2000);
  processMetrics();
}

const staticLogs = [
  { target: "ورود موفق مدیریت", code: "200 OK", ip_address: "192.168.1.1" },
  {
    target: "پشتیبان‌گیری دیتابیس",
    code: "201 Created",
    ip_address: "10.0.0.5",
  },
  {
    target: "تلاش برای دسترسی غیرمجاز",
    code: "403 Forbidden",
    ip_address: "185.22.10.4",
  },
];

function renderLogs(logs) {
  if (!logContainer) return;
  logContainer.innerHTML = "";

  logs.forEach((entry) => {
    let time;

    if (entry.time) {
      time = new Date(entry.time).toLocaleString("fa-IR", {
        timeZone: "Asia/Tehran",
      });
    } else {
      time = new Date().toLocaleString("fa-IR", { timeZone: "Asia/Tehran" });
    }

    const row = document.createElement("tr");
    row.innerHTML = `
            <td>${time}</td>
            <td class="fw-bold">${entry.action || entry.target || "N/A"}</td>
            <td><code>${entry.code || "N/A"}</code></td>
            <td>${entry.system_info || "N/A"}</td>
            <td><code>${entry.ip_address || "N/A"}</code></td>
        `;
    logContainer.appendChild(row);
  });
}


async function loadLogs() {
  let logs = JSON.parse(localStorage.getItem("userLogs"));
  const username = localStorage.getItem("username"); // فرض بر اینکه هنگام لاگین ست کردید

  if (!logs || logs.length === 0) {
    try {
      const response = await fetch(`/api/users/logs?username=${username}`);
      const data = await response.json();

      logs = data.logs;
      console.log(logs);
      if (logs && logs.length > 0) {
        localStorage.setItem("userLogs", JSON.stringify(logs));
      } else {
        logs = staticLogs;
      }
    } catch (error) {
      logs = staticLogs;
    }
  }
  renderLogs(logs);
}

// function renderLogs() {
//   if (!logContainer) return;
//   logContainer.innerHTML = "";
//   staticLogs.forEach((entry) => {
//     const row = document.createElement("tr");
//     const stamp = new Date().toLocaleTimeString("fa-IR");
//     row.innerHTML = `
//             <td>${stamp}</td>
//             <td class="fw-bold">${entry.target}</td>
//             <td><code>${entry.code}</code></td>
//             <td>${entry.ip}</td>
//         `;
//     logContainer.appendChild(row);
//   });
// }

if (logContainer) {
  loadLogs();
}

function restartServer(serverName) {
  showAlert(`سرور ${serverName} با موفقیت راه‌اندازی مجدد شد.`, "warning");
}

function startServer(serverName) {
  showAlert(`سرور ${serverName} با موفقیت روشن و به شبکه متصل شد.`, "success");
}

function optimizeDb(dbName) {
  showAlert(
    `عملیات ایندکس‌گذاری و بهینه‌سازی ${dbName} با موفقیت انجام شد.`,
    "success",
  );
}

function clearCache() {
  showAlert("حافظه موقت رم (Redis Cache) کاملاً پاکسازی شد.", "info");
}

function saveSettings() {
  showAlert("تنظیمات هشدارهای سیستم با موفقیت به‌روزرسانی شد.", "success");
}

let userCounter = 5;

function disconnectUser(username, rowId) {
  const row = document.getElementById(rowId);
  if (row) {
    row.remove();
    showAlert(` کاربر ${username} با موفقیت حذف شد.`, "danger");
  }
}

function createUser(event) {
  event.preventDefault();

  const usernameInput = document.getElementById("new-username");
  const passwordInput = document.getElementById("new-password");
  const tableBody = document.getElementById("user-table-body");

  if (!usernameInput || !passwordInput || !tableBody) return;

  const username = usernameInput.value.trim();
  const password = passwordInput.value;

  if (username === "" || password === "") return;

  const rowId = `user-row-${userCounter}`;
  const newRow = document.createElement("tr");
  newRow.id = rowId;
  newRow.innerHTML = `
        <td class="fw-bold">${username}</td>
        <td>${password}</td>
        <td><span class="badge bg-success">فعال</span></td>
        <td>192.168.1.${userCounter}</td>
        <td>تهران، ایران</td>
        <td>Chrome / Windows</td>
        <td>
            <button class="btn btn-danger btn-sm" onclick="disconnectUser('${username}', '${rowId}')"> حذف </button>
        </td>
    `;

  tableBody.appendChild(newRow);
  showAlert(
    `کاربر جدید "${username}" با موفقیت به دیتابیس نشست‌ها اضافه شد.`,
    "success",
  );

  usernameInput.value = "";
  passwordInput.value = "";

  userCounter++;
}

const userForm = document.getElementById("add-user-form");
if (userForm) {
  userForm.addEventListener("submit", createUser);
}

function toggleServerState(serverName) {
  const btn = document.getElementById(`btn-${serverName}`);
  const badge = document.getElementById(`status-badge-${serverName}`);
  if (!btn || !badge) return;

  const statusTextNode = badge.querySelector(".status-text");
  const pingDot = badge.querySelector(".ping-dot");
  const currentTimeString = new Date().toLocaleTimeString("fa-IR");

  if (btn.classList.contains("btn-danger")) {
    btn.classList.replace("btn-danger", "btn-success");
    btn.innerText = "روشن کردن سرور";

    badge.className =
      "badge bg-danger-subtle text-danger d-flex align-items-center gap-1 border border-danger-subtle";
    if (statusTextNode) statusTextNode.innerText = "Offline";

    if (pingDot) {
      pingDot.className = "ping-dot bg-danger";
    }

    const downtimeDisplay = document.getElementById(`downtime-${serverName}`);
    if (downtimeDisplay) downtimeDisplay.innerText = currentTimeString;

    showAlert(`سرور ${serverName} خاموش شد و از مدار خارج گردید.`, "danger");
  } else {
    btn.classList.replace("btn-success", "btn-danger");
    btn.innerText = "خاموش کردن سرور";

    badge.className =
      "badge bg-success-subtle text-success d-flex align-items-center gap-1 border border-success-subtle";
    if (statusTextNode) statusTextNode.innerText = "Online";

    if (pingDot) {
      pingDot.className = "ping-dot bg-success ping-flash";
    }

    const uptimeDisplay = document.getElementById(`uptime-${serverName}`);
    if (uptimeDisplay) uptimeDisplay.innerText = currentTimeString;

    showAlert(
      `سرور ${serverName} با موفقیت روشن شد و وضعیت آن پایدار است.`,
      "success",
    );
  }
}

function logout() {
  event.preventDefault();
  localStorage.clear();
  window.location.replace("/host/");
}

function updateFullName() {
  const fullname = localStorage.getItem("fullName");
  const username = localStorage.getItem("username");
  document.getElementById("fullname").innerText = fullname;
  if (fullname == "undefined") {
    document.getElementById("fullname").innerText = username;
  } else {
    document.getElementById("fullname").innerText = fullname;
  }
}
updateFullName();


function exit() {
  localStorage.removeItem("userLogs")
  window.location.replace("/host/");
  
}