const locationsData = {
  us: {
    name: "آمریکا (شیکاگو)",
    ping: "140 ms",
    flag: "fi-us",
    serversCount: 6,
    hardware: "Dual Intel Xeon Platinum",
    uptime: "99.99%",
  },
  ca: {
    name: "کانادا (تورنتو)",
    ping: "155 ms",
    flag: "fi-ca",
    serversCount: 4,
    hardware: "AMD EPYC Milian",
    uptime: "99.90%",
  },
  ir: {
    name: "ایران (تهران)",
    ping: "4 ms",
    flag: "fi-ir",
    serversCount: 5,
    hardware: "Intel Xeon 2680 v4",
    uptime: "99.99%",
  },
  de: {
    name: "آلمان (فرانکفورت)",
    ping: "32 ms",
    flag: "fi-de",
    serversCount: 6,
    hardware: "AMD Ryzen 9 7950X",
    uptime: "99.95%",
  },
  fr: {
    name: "فرانسه (پاریس)",
    ping: "45 ms",
    flag: "fi-fr",
    serversCount: 4,
    hardware: "Intel Core i9-14900K",
    uptime: "99.90%",
  },
  nl: {
    name: "هلند (آمستردام)",
    ping: "38 ms",
    flag: "fi-nl",
    serversCount: 5,
    hardware: "AMD EPYC 7402",
    uptime: "99.95%",
  },
  tr: {
    name: "ترکیه (استانبول)",
    ping: "22 ms",
    flag: "fi-tr",
    serversCount: 3,
    hardware: "Intel Xeon Gold",
    uptime: "98.85%",
  },
  ae: {
    name: "امارات (دبی)",
    ping: "15 ms",
    flag: "fi-ae",
    serversCount: 4,
    hardware: "AMD Ryzen 7",
    uptime: "99.90%",
  },
  fi: {
    name: "فنلاند (هلسینکی)",
    ping: "55 ms",
    flag: "fi-fi",
    serversCount: 3,
    hardware: "Intel Xeon Scalable",
    uptime: "99.95%",
  },
  pl: {
    name: "لهستان (ورشو)",
    ping: "48 ms",
    flag: "fi-pl",
    serversCount: 4,
    hardware: "Intel Gold 6230",
    uptime: "99.80%",
  },
  uk: {
    name: "انگلیس (لندن)",
    ping: "40 ms",
    flag: "fi-gb",
    serversCount: 5,
    hardware: "AMD Ryzen 9 5950X",
    uptime: "99.95%",
  },
  sg: {
    name: "سنگاپور (مرکزی)",
    ping: "185 ms",
    flag: "fi-sg",
    serversCount: 3,
    hardware: "Intel Xeon Silver",
    uptime: "99.75%",
  },
};

document.addEventListener("DOMContentLoaded", function () {
  initCountriesMenu();
  renderCountryServers("ir");
  updateConfig();
});

function initCountriesMenu() {
  const listContainer = document.getElementById("countriesList");
  if (!listContainer) return;
  listContainer.innerHTML = "";

  Object.keys(locationsData).forEach((key) => {
    const loc = locationsData[key];
    const isActive = key === "ir" ? "active bg-success border-success" : "";
    const textClass = key === "ir" ? "text-white" : "text-dark";

    listContainer.innerHTML += `
            <button onclick="switchLocation('${key}', this)" class="list-group-item list-group-item-action ${isActive} d-flex align-items-center justify-content-between py-3 rounded-3 mb-1 border-0 shadow-sm">
                <span class="fw-bold small ${textClass}"><span class="fi ${loc.flag} me-2 rounded-1 shadow-sm"></span>${loc.name}</span>
                <span class="badge bg-light text-dark border rounded-pill px-2 small">${loc.ping}</span>
            </button>
        `;
  });
}

function switchLocation(countryKey, element) {
  const items = document.querySelectorAll("#countriesList button");
  items.forEach((item) => {
    item.classList.remove("active", "bg-success", "border-success");
    const textSpan = item.querySelector("span.fw-bold");
    if (textSpan) {
      textSpan.classList.remove("text-white");
      textSpan.classList.add("text-dark");
    }
  });

  element.classList.add("active", "bg-success", "border-success");
  element.querySelector("span.fw-bold").classList.remove("text-dark");
  element.querySelector("span.fw-bold").classList.add("text-white");

  renderCountryServers(countryKey);
}

function renderCountryServers(countryKey) {
  const loc = locationsData[countryKey];
  const container = document.getElementById("serversContainer");
  if (!container) return;

  document.getElementById("selectedCountryName").innerText =
    `سرورهای فعال لوکیشن ${loc.name}`;
  document.getElementById("selectedCountryPing").innerText = loc.ping;
  document.getElementById("selectedCountryFlag").className =
    `fi ${loc.flag} fi-de me-2 rounded-2`;

  container.innerHTML = "";

  for (let i = 1; i <= loc.serversCount; i++) {
    const core = i * 2;
    const ram = i * 4;
    const storage = i * 40;
    const priceValue = i * 450000 ;
    const formattedPrice = priceValue.toLocaleString("fa");
    container.innerHTML += `
            <div class="col-md-6">
                <div class="card border shadow-sm p-3 rounded-4 bg-white h-100 d-flex flex-column justify-content-between">
                    <div>
                        <div class="d-flex justify-content-between align-items-center mb-2">
                            <h6 class="fw-bold text-dark mb-0">پلان ابری Enterprise-${i}</h6>
                            <span class="badge bg-success-subtle text-success rounded-pill small">آپتایم ${loc.uptime}</span>
                        </div>
                        <p class="text-muted text-start" style="font-size: 11px;"><i class="bi bi-cpu-fill me-1 text-secondary"></i>سخت‌افزار میزبان: ${loc.hardware}</p>
                        <hr class="opacity-25 my-2">
                        <ul class="list-unstyled small text-muted d-flex flex-column gap-1 text-start">
                            <li><i class="bi bi-cpu text-success me-2"></i>پردازنده: <strong>${core} هسته اختصاصی</strong></li>
                            <li><i class="bi bi-memory text-success me-2"></i>حافظه موقت: <strong>${ram} گیگابایت RAM DDR5</strong></li>
                            <li><i class="bi bi-database text-success me-2"></i>هارد دیسک: <strong>${storage} GB NVMe نسل ۴</strong></li>
                            <li><i class="bi bi-shield-check text-success me-2"></i>امنیت: <strong>فایروال پورت سخت‌افزاری</strong></li>
                        </ul>
                    </div>
                    <div class="mt-3">
                        <div class="d-flex justify-content-between align-items-center bg-light p-2 rounded-3 mb-2">
                            <span class="small text-muted fw-bold">هزینه ماهانه:</span>
                            <span class="fw-bold text-success small">${formattedPrice} تومان</span>
                        </div>
                        <button onclick="buyServer('${loc.name} - پلان Enterprise-${i }')" class="btn btn-success btn-sm w-100 rounded-pill fw-bold">خرید ماشین ابری</button>
                    </div>
                </div>
            </div>
        `;
  }
}

function updateConfig() {
  const cpuSlider = document.getElementById("cpuSlider");
  const ramSlider = document.getElementById("ramSlider");
  const diskSlider = document.getElementById("diskSlider");
  const diskTypeSelector = document.getElementById("diskTypeSelector");

  if (!cpuSlider || !ramSlider || !diskSlider || !diskTypeSelector) return;

  const cpu = parseInt(cpuSlider.value);
  const ram = parseInt(ramSlider.value);
  const disk = parseInt(diskSlider.value);
  const diskType = diskTypeSelector.value;

  document.getElementById("cpuVal").innerText = cpu + " هسته";
  document.getElementById("ramVal").innerText = ram + " گیگابایت";
  document.getElementById("diskVal").innerText = disk + " گیگابایت";

  let basePricePerMonth = cpu * 165000 + ram * 66000 + disk * 4800;

  if (diskType === "ssd") {
    basePricePerMonth = basePricePerMonth * 0.9; 
  } else if (diskType === "hdd") {
    basePricePerMonth = basePricePerMonth * 0.5; 
  }

  const selectedCycle = document.querySelector(
    'input[name="billingCycle"]:checked',
  ).value;
  const months = parseInt(selectedCycle);

  let totalPrice = basePricePerMonth * months;

  if (months === 3) {
    totalPrice = totalPrice * 0.9;
  } else if (months === 12) {
    totalPrice = totalPrice * 0.8;
  }

  document.getElementById("totalPrice").innerText =
    Math.round(totalPrice);
}

function selectOS(osName) {
  alert(`سیستم‌عامل پیش‌فرض برای راه‌اندازی سرور به [${osName}] تغییر یافت.`);
}

function orderCustomServer() {
  alert(
    "پیکربندی هوشمند سرور سفارشی شما با موفقیت تایید شد. در حال هدایت به پیش‌فاکتور خرید...",
  );
}

function buyServer(serverDetails) {
  alert(
    `سرویس [${serverDetails}] انتخاب شد. انتقال به درگاه صادرکننده لایسنس...`,
  );
}
