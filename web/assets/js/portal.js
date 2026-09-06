/**
 * Lalatina Systems — User Self-Service Portal Controller
 * Pure Native Modular Vanilla JS (Zero heavy dependencies)
 */

(function () {
  'use strict';

  // --- State Management ---
  const state = {
    token: null,
    profile: null,
    activeTab: 'clash',
    telemetryTimer: null,
  };

  // --- DOM Elements Cache ---
  const el = {
    // Views
    loginView: document.getElementById('portal-login-view'),
    dashboardView: document.getElementById('portal-dashboard-view'),
    loadingOverlay: document.getElementById('loading-indicator'),
    toastContainer: document.getElementById('toast-container'),

    // Login Form & Actions
    loginForm: document.getElementById('login-form'),
    tokenInput: document.getElementById('token-input'),
    btnLogin: document.getElementById('btn-login'),
    btnTrial: document.getElementById('btn-create-trial'),

    // Profile & Header Elements
    headerToken: document.getElementById('display-header-token'),
    userTokenDisplay: document.getElementById('display-token'),
    userUUIDDisplay: document.getElementById('display-uuid'),
    userExpiredDisplay: document.getElementById('display-expired'),
    userStatusBadge: document.getElementById('display-status-badge'),
    userQuotaDisplay: document.getElementById('display-quota'),
    quotaMeterTrack: document.getElementById('quota-meter-track'),
    btnLogout: document.getElementById('btn-logout'),

    // Config Switcher
    selectVpn: document.getElementById('select-vpn'),
    selectServer: document.getElementById('select-server'),
    selectRelay: document.getElementById('select-relay'),
    btnSaveConfig: document.getElementById('btn-save-config'),

    // Security Controls
    btnResetUUID: document.getElementById('btn-reset-uuid'),
    btnChangeToken: document.getElementById('btn-change-token'),
    adblockToggle: document.getElementById('adblock-toggle'),

    // Invocations & Tabs
    tabBtns: document.querySelectorAll('.portal-tab'),
    invocationText: document.getElementById('invocation-text'),
    btnCopyConfig: document.getElementById('btn-copy-config'),
    btnOpenApp: document.getElementById('btn-open-app'),
    qrcodeContainer: document.getElementById('qrcode-box'),

    // Telemetry Elements
    cpuValue: document.getElementById('telemetry-cpu-val'),
    cpuMeter: document.getElementById('telemetry-cpu-meter'),
    ramValue: document.getElementById('telemetry-ram-val'),
    ramMeter: document.getElementById('telemetry-ram-meter'),
    diskValue: document.getElementById('telemetry-disk-val'),
    diskMeter: document.getElementById('telemetry-disk-meter'),
    uptimeValue: document.getElementById('telemetry-uptime-val'),

    // Wildcards Table
    wildcardsList: document.getElementById('wildcards-table-body'),
  };

  // --- UI Feedback & Notification (Toast) ---
  function showToast(message, type = 'info') {
    if (!el.toastContainer) return;
    const toast = document.createElement('div');
    toast.className = `toast ${type === 'error' ? 'toast-error' : ''}`;
    toast.setAttribute('role', 'alert');
    
    const icon = type === 'error' ? '✕' : (type === 'success' ? '✓' : 'ℹ');
    toast.innerHTML = `<span style="font-weight:700; color:${type === 'error' ? 'var(--accent-ruby)' : 'var(--accent-signal)'}">${icon}</span> <span>${escapeHtml(message)}</span>`;

    el.toastContainer.appendChild(toast);
    setTimeout(() => {
      toast.style.opacity = '0';
      toast.style.transform = 'translateY(-10px)';
      toast.style.transition = 'all 0.3s ease';
      setTimeout(() => toast.remove(), 300);
    }, 4000);
  }

  function escapeHtml(str) {
    if (!str) return '';
    return str.toString().replace(/[&<>"']/g, function (m) {
      return { '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' }[m];
    });
  }

  function setLoading(isLoading) {
    if (el.loadingOverlay) {
      el.loadingOverlay.style.display = isLoading ? 'flex' : 'none';
    }
  }

  // --- Token Extraction & Initialization ---
  function initAuth() {
    const params = new URLSearchParams(window.location.search);
    const urlToken = params.get('token');

    if (urlToken && urlToken.trim().length >= 4) {
      state.token = urlToken.trim();
      localStorage.setItem('latina_portal_token', state.token);
      // Clean query parameter from URL without page reload
      window.history.replaceState({}, document.title, window.location.pathname);
    } else {
      state.token = localStorage.getItem('latina_portal_token');
    }

    if (state.token) {
      loadProfile(state.token);
    } else {
      showLoginView();
    }
  }

  function showLoginView() {
    if (state.telemetryTimer) {
      clearInterval(state.telemetryTimer);
      state.telemetryTimer = null;
    }
    if (el.loginView) el.loginView.style.display = 'block';
    if (el.dashboardView) el.dashboardView.style.display = 'none';
    if (el.tokenInput) el.tokenInput.focus();
  }

  function showDashboardView() {
    if (el.loginView) el.loginView.style.display = 'none';
    if (el.dashboardView) el.dashboardView.style.display = 'block';
    startTelemetryPolling();
  }

  // --- API Communications ---
  async function loadProfile(token) {
    setLoading(true);
    try {
      const res = await fetch(`/api/v1/portal/profile?token=${encodeURIComponent(token)}`, {
        headers: { 'Accept': 'application/json' }
      });

      if (!res.ok) {
        if (res.status === 401 || res.status === 404) {
          throw new Error('Token tidak valid atau akun tidak ditemukan');
        }
        const errData = await res.json().catch(() => ({}));
        throw new Error(errData.error || `Gagal mengambil profil (HTTP ${res.status})`);
      }

      const profile = await res.json();
      state.profile = profile;
      state.token = profile.token;
      localStorage.setItem('latina_portal_token', profile.token);

      renderProfile(profile);
      showDashboardView();
      loadWildcards();
      showToast('Autentikasi berhasil', 'success');
    } catch (err) {
      showToast(err.message, 'error');
      localStorage.removeItem('latina_portal_token');
      state.token = null;
      showLoginView();
    } finally {
      setLoading(false);
    }
  }

  // --- Profile Rendering ---
  function renderProfile(data) {
    if (el.headerToken) el.headerToken.textContent = data.token;
    if (el.userTokenDisplay) el.userTokenDisplay.textContent = data.token;
    if (el.userUUIDDisplay) el.userUUIDDisplay.textContent = data.password || '—';
    if (el.userExpiredDisplay) {
      el.userExpiredDisplay.textContent = `${data.expired} (${data.expired_human})`;
    }

    if (el.userStatusBadge) {
      if (data.is_expired) {
        el.userStatusBadge.textContent = 'KEDALUWARSA';
        el.userStatusBadge.style.color = 'var(--accent-ruby)';
        el.userStatusBadge.style.borderColor = 'var(--accent-ruby-border)';
      } else {
        el.userStatusBadge.textContent = 'AKTIF';
        el.userStatusBadge.style.color = 'var(--status-active)';
        el.userStatusBadge.style.borderColor = 'rgba(16, 185, 129, 0.4)';
      }
    }

    if (el.userQuotaDisplay) {
      el.userQuotaDisplay.textContent = data.quotaHuman || 'Unlimited';
    }

    renderSegmentedQuota(data);
    populateConfigDropdowns(data);

    if (el.adblockToggle) {
      el.adblockToggle.checked = Boolean(data.adblock);
    }

    renderActiveInvocation();
  }

  function renderSegmentedQuota(data) {
    if (!el.quotaMeterTrack) return;
    el.quotaMeterTrack.innerHTML = '';
    const totalSegments = 20;

    // By default for unlimited accounts, display full active segments
    let filledSegments = totalSegments;
    if (data.is_expired) {
      filledSegments = 0;
    }

    for (let i = 0; i < totalSegments; i++) {
      const seg = document.createElement('div');
      seg.className = 'meter-segment';
      if (i < filledSegments) {
        seg.classList.add(data.is_expired ? 'filled-ruby' : 'filled');
      }
      el.quotaMeterTrack.appendChild(seg);
    }
  }

  function populateConfigDropdowns(data) {
    // Protocol
    if (el.selectVpn) {
      el.selectVpn.value = (data.vpn || 'vmess').toLowerCase();
    }

    // Servers
    if (el.selectServer && Array.isArray(data.servers)) {
      el.selectServer.innerHTML = '';
      data.servers.forEach(s => {
        const opt = document.createElement('option');
        opt.value = s.code;
        opt.textContent = `${s.code.toUpperCase()} — ${s.name || s.domain} (${s.domain})`;
        if (s.code === data.server_code) opt.selected = true;
        el.selectServer.appendChild(opt);
      });
    }

    // Relays
    if (el.selectRelay && Array.isArray(data.relays)) {
      el.selectRelay.innerHTML = '';
      data.relays.forEach(r => {
        const opt = document.createElement('option');
        opt.value = r;
        opt.textContent = r;
        if (r === data.relay || (!data.relay && r === 'Tanpa Relay')) opt.selected = true;
        el.selectRelay.appendChild(opt);
      });
    }
  }

  // --- Invocations & Tabs Rendering ---
  function renderActiveInvocation() {
    if (!state.profile) return;
    const p = state.profile;
    let textToDisplay = '';
    let appLaunchUrl = '';

    switch (state.activeTab) {
      case 'clash':
        textToDisplay = p.clash_url || p.sub_url;
        appLaunchUrl = p.clash_url;
        break;
      case 'singbox':
        textToDisplay = p.singbox_url || p.sub_url;
        appLaunchUrl = p.singbox_url;
        break;
      case 'shadowrocket':
        textToDisplay = p.shadowrocket_url || p.sub_url;
        appLaunchUrl = p.shadowrocket_url;
        break;
      case 'raw':
      default:
        textToDisplay = p.raw_uri || p.sub_url;
        appLaunchUrl = p.raw_uri;
        break;
    }

    if (el.invocationText) {
      el.invocationText.textContent = textToDisplay;
    }

    if (el.btnOpenApp) {
      if (appLaunchUrl && appLaunchUrl.includes('://')) {
        el.btnOpenApp.href = appLaunchUrl;
        el.btnOpenApp.style.display = 'inline-flex';
      } else {
        el.btnOpenApp.style.display = 'none';
      }
    }

    // Update QR Code
    renderQRCode(p.raw_uri || p.sub_url);
  }

  function renderQRCode(text) {
    if (!el.qrcodeContainer || !text) return;
    el.qrcodeContainer.innerHTML = '';
    if (typeof window.QRCode === 'function') {
      try {
        new window.QRCode(el.qrcodeContainer, {
          text: text,
          width: 200,
          height: 200,
          colorDark: '#0c0a08',
          colorLight: '#ffffff',
          correctLevel: window.QRCode.CorrectLevel.M,
        });
      } catch (err) {
        el.qrcodeContainer.innerHTML = '<p style="font-size:0.75rem;color:var(--text-muted);">QR Code tidak dapat digenerate</p>';
      }
    }
  }

  // --- Telemetry Polling ---
  async function fetchTelemetry() {
    try {
      const res = await fetch('/api/v1/portal/status');
      if (!res.ok) return;
      const data = await res.json();

      // CPU
      if (el.cpuValue && typeof data.cpu_usage === 'number') {
        el.cpuValue.textContent = `${data.cpu_usage.toFixed(1)}%`;
        renderMeterSegments(el.cpuMeter, data.cpu_usage);
      }

      // RAM
      if (el.ramValue && typeof data.ram_usage === 'number') {
        el.ramValue.textContent = `${data.ram_usage.toFixed(1)}%`;
        renderMeterSegments(el.ramMeter, data.ram_usage);
      }

      // Disk
      if (el.diskValue && typeof data.disk_usage === 'number') {
        el.diskValue.textContent = `${data.disk_usage.toFixed(1)}%`;
        renderMeterSegments(el.diskMeter, data.disk_usage);
      }

      // Uptime
      if (el.uptimeValue && data.uptime) {
        el.uptimeValue.textContent = data.uptime;
      }
    } catch (e) {
      // Passive polling fails gracefully
    }
  }

  function renderMeterSegments(container, percentage) {
    if (!container) return;
    container.innerHTML = '';
    const total = 10;
    const filledCount = Math.round((percentage / 100) * total);

    for (let i = 0; i < total; i++) {
      const s = document.createElement('div');
      s.className = 'meter-segment';
      if (i < filledCount) {
        s.classList.add(percentage > 85 ? 'filled-ruby' : 'filled');
      }
      container.appendChild(s);
    }
  }

  function startTelemetryPolling() {
    fetchTelemetry();
    if (state.telemetryTimer) clearInterval(state.telemetryTimer);
    state.telemetryTimer = setInterval(fetchTelemetry, 12000);
  }

  // --- Wildcards Fetching ---
  async function loadWildcards() {
    if (!el.wildcardsList) return;
    try {
      const res = await fetch('/api/v1/portal/wildcards');
      if (!res.ok) return;
      const list = await res.json();
      if (!Array.isArray(list)) return;

      el.wildcardsList.innerHTML = '';
      list.forEach(w => {
        const row = document.createElement('tr');
        row.innerHTML = `
          <td style="font-weight:600; color:var(--text-primary); font-family:var(--font-mono);">${escapeHtml(w)}</td>
          <td><span class="mono-label">SNI BUG HOST</span></td>
          <td style="text-align:right;">
            <button type="button" class="btn btn-outline btn-sm copy-wildcard-btn" data-host="${escapeHtml(w)}">Salin Host</button>
          </td>
        `;
        el.wildcardsList.appendChild(row);
      });

      // Bind copy buttons
      el.wildcardsList.querySelectorAll('.copy-wildcard-btn').forEach(btn => {
        btn.addEventListener('click', (e) => {
          const host = e.currentTarget.getAttribute('data-host');
          copyToClipboard(host, 'SNI Bug Host disalin ke clipboard');
        });
      });
    } catch (e) {
      // Graceful fallback
    }
  }

  // --- Clipboard Helper ---
  function copyToClipboard(text, successMsg = 'Tersalin ke clipboard') {
    if (!text) return;
    if (navigator.clipboard && navigator.clipboard.writeText) {
      navigator.clipboard.writeText(text).then(() => {
        showToast(successMsg, 'success');
      }).catch(() => fallbackCopy(text, successMsg));
    } else {
      fallbackCopy(text, successMsg);
    }
  }

  function fallbackCopy(text, successMsg) {
    const textarea = document.createElement('textarea');
    textarea.value = text;
    textarea.style.position = 'fixed';
    textarea.style.opacity = '0';
    document.body.appendChild(textarea);
    textarea.select();
    try {
      document.execCommand('copy');
      showToast(successMsg, 'success');
    } catch (err) {
      showToast('Gagal menyalin teks', 'error');
    }
    document.body.removeChild(textarea);
  }

  // --- Event Bindings ---
  function setupEventListeners() {
    // 1. Login Form Submission
    if (el.loginForm) {
      el.loginForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const inputVal = el.tokenInput ? el.tokenInput.value.trim() : '';
        if (inputVal.length < 4) {
          showToast('Masukkan minimal 4 karakter token', 'error');
          return;
        }
        await loadProfile(inputVal);
      });
    }

    // 2. 1-Click Trial Generation
    if (el.btnTrial) {
      el.btnTrial.addEventListener('click', async () => {
        if (!confirm('Buat akun uji coba dengan masa aktif 24 jam dan batas kecepatan 2 Mbps?')) return;
        setLoading(true);
        try {
          const res = await fetch('/api/v1/trial', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' }
          });

          if (!res.ok) {
            const errData = await res.json().catch(() => ({}));
            throw new Error(errData.error || errData.message || 'Gagal membuat akun uji coba');
          }

          const data = await res.json();
          const trialToken = data.token || (data.config && data.config.token) || data.id;

          if (trialToken) {
            showToast('Akun uji coba 24 jam berhasil dibuat', 'success');
            await loadProfile(trialToken);
          } else if (data.vmess_link || data.raw_uri) {
            // If trial only returns raw link without user entry
            showToast('Akun uji coba berhasil dibuat', 'success');
            if (el.tokenInput) el.tokenInput.value = data.id || '';
            copyToClipboard(data.vmess_link || data.raw_uri, 'Link VMess Uji Coba disalin ke clipboard');
          } else {
            showToast('Akun uji coba berhasil diaktifkan', 'success');
          }
        } catch (err) {
          showToast(err.message, 'error');
        } finally {
          setLoading(false);
        }
      });
    }

    // 3. Logout
    if (el.btnLogout) {
      el.btnLogout.addEventListener('click', () => {
        if (confirm('Keluar dari portal mandiri?')) {
          localStorage.removeItem('latina_portal_token');
          state.token = null;
          state.profile = null;
          showLoginView();
          showToast('Sesi telah diakhiri');
        }
      });
    }

    // 4. Invocations Tabs
    el.tabBtns.forEach(btn => {
      btn.addEventListener('click', () => {
        el.tabBtns.forEach(b => b.classList.remove('active'));
        btn.classList.add('active');
        state.activeTab = btn.getAttribute('data-tab');
        renderActiveInvocation();
      });
    });

    // 5. Copy Invocation / Config
    if (el.btnCopyConfig) {
      el.btnCopyConfig.addEventListener('click', () => {
        const text = el.invocationText ? el.invocationText.textContent : '';
        copyToClipboard(text, 'Konfigurasi berhasil disalin ke clipboard');
      });
    }

    // 6. Save Configuration (Protocol / Server / Relay)
    if (el.btnSaveConfig) {
      el.btnSaveConfig.addEventListener('click', async () => {
        if (!state.token) return;
        const vpn = el.selectVpn ? el.selectVpn.value : 'vmess';
        const serverCode = el.selectServer ? el.selectServer.value : '';
        const relayVal = el.selectRelay ? el.selectRelay.value : 'Tanpa Relay';

        setLoading(true);
        try {
          const res = await fetch('/api/v1/portal/update-config', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
              token: state.token,
              vpn: vpn,
              server_code: serverCode,
              relay: relayVal,
            }),
          });

          if (!res.ok) {
            const errData = await res.json().catch(() => ({}));
            throw new Error(errData.error || 'Gagal memperbarui konfigurasi');
          }

          showToast('Konfigurasi berhasil disimpan dan service di-reload', 'success');
          // Reload profile to update URLs and raw URI
          await loadProfile(state.token);
        } catch (err) {
          showToast(err.message, 'error');
        } finally {
          setLoading(false);
        }
      });
    }

    // 7. Reset UUID
    if (el.btnResetUUID) {
      el.btnResetUUID.addEventListener('click', async () => {
        if (!confirm('Peringatan: Mereset UUID akan memutuskan koneksi VPN yang sedang berjalan hingga konfigurasi baru diimpor ke aplikasi Anda. Lanjutkan?')) return;
        setLoading(true);
        try {
          const res = await fetch('/api/v1/portal/reset-uuid', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ token: state.token }),
          });

          if (!res.ok) {
            const errData = await res.json().catch(() => ({}));
            throw new Error(errData.error || 'Gagal mereset UUID');
          }

          showToast('UUID baru berhasil di-generate', 'success');
          await loadProfile(state.token);
        } catch (err) {
          showToast(err.message, 'error');
        } finally {
          setLoading(false);
        }
      });
    }

    // 8. Change Token
    if (el.btnChangeToken) {
      el.btnChangeToken.addEventListener('click', async () => {
        if (!confirm('Peringatan: Token akses login Anda akan diganti dengan 8-karakter baru. Simpan token baru setelah rotasi selesai. Lanjutkan?')) return;
        setLoading(true);
        try {
          const res = await fetch('/api/v1/portal/change-token', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ token: state.token }),
          });

          if (!res.ok) {
            const errData = await res.json().catch(() => ({}));
            throw new Error(errData.error || 'Gagal mengganti token');
          }

          const data = await res.json();
          if (data.token) {
            state.token = data.token;
            localStorage.setItem('latina_portal_token', data.token);
            showToast(`Token baru: ${data.token}. Harap catat token ini.`, 'success');
            await loadProfile(data.token);
          }
        } catch (err) {
          showToast(err.message, 'error');
        } finally {
          setLoading(false);
        }
      });
    }

    // 9. Toggle Adblock
    if (el.adblockToggle) {
      el.adblockToggle.addEventListener('change', async (e) => {
        if (!state.token) return;
        const targetChecked = e.target.checked;
        setLoading(true);
        try {
          const res = await fetch('/api/v1/portal/toggle-adblock', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
              token: state.token,
              adblock: targetChecked,
            }),
          });

          if (!res.ok) {
            e.target.checked = !targetChecked; // Revert checkbox
            const errData = await res.json().catch(() => ({}));
            throw new Error(errData.error || 'Gagal mengubah status AdBlock');
          }

          showToast(`Proteksi AdBlock ${targetChecked ? 'diaktifkan' : 'dinonaktifkan'}`, 'success');
        } catch (err) {
          showToast(err.message, 'error');
        } finally {
          setLoading(false);
        }
      });
    }
  }

  // --- Bootstrap on DOM Ready ---
  document.addEventListener('DOMContentLoaded', () => {
    setupEventListeners();
    initAuth();
  });
})();
