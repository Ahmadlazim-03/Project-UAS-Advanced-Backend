const API_URL = '/api/v1';
let token = localStorage.getItem('token');
let user = JSON.parse(localStorage.getItem('user'));

// DOM Elements
const loginView = document.getElementById('loginView');
const dashboardView = document.getElementById('dashboardView');
const loginForm = document.getElementById('loginForm');
const userInfo = document.getElementById('userInfo');
const logoutBtn = document.getElementById('logoutBtn');
const achievementsTable = document.getElementById('achievementsTable').querySelector('tbody');
const addAchievementBtn = document.getElementById('addAchievementBtn');
const addModal = new bootstrap.Modal(document.getElementById('addModal'));
const addForm = document.getElementById('addForm');

// Init
function init() {
    if (token && user) {
        showDashboard();
    } else {
        showLogin();
    }
}

function showLogin() {
    loginView.style.display = 'block';
    dashboardView.style.display = 'none';
    userInfo.textContent = '';
    logoutBtn.style.display = 'none';
}

function showDashboard() {
    loginView.style.display = 'none';
    dashboardView.style.display = 'block';
    userInfo.textContent = `Welcome, ${user.fullName} (${user.role})`;
    logoutBtn.style.display = 'block';
    loadAchievements();
}

// Login
loginForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;

    try {
        const res = await fetch(`${API_URL}/auth/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password })
        });
        const data = await res.json();

        if (data.status === 'success') {
            token = data.data.token;
            user = data.data.user;
            localStorage.setItem('token', token);
            localStorage.setItem('user', JSON.stringify(user));
            showDashboard();
        } else {
            alert(data.message);
        }
    } catch (err) {
        console.error(err);
        alert('Login failed');
    }
});

// Logout
logoutBtn.addEventListener('click', () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    token = null;
    user = null;
    showLogin();
});

// Load Achievements
async function loadAchievements() {
    if (!token) return;

    try {
        const res = await fetch(`${API_URL}/achievements`, {
            headers: { 'Authorization': `Bearer ${token}` }
        });
        const data = await res.json();

        if (data.status === 'success') {
            renderAchievements(data.data);
        }
    } catch (err) {
        console.error(err);
    }
}

function renderAchievements(achievements) {
    achievementsTable.innerHTML = '';
    if (!achievements) return;
    
    achievements.forEach(item => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${item.data.title}</td>
            <td>${item.data.achievementType}</td>
            <td><span class="badge bg-${getStatusColor(item.status)}">${item.status}</span></td>
            <td>${item.data.points}</td>
            <td>
                <button class="btn btn-sm btn-info" onclick="viewDetail('${item.id}')">View</button>
                ${item.status === 'draft' ? `<button class="btn btn-sm btn-primary" onclick="submitAchievement('${item.id}')">Submit</button>` : ''}
            </td>
        `;
        achievementsTable.appendChild(row);
    });
}

function getStatusColor(status) {
    switch(status) {
        case 'draft': return 'secondary';
        case 'submitted': return 'warning';
        case 'verified': return 'success';
        case 'rejected': return 'danger';
        default: return 'primary';
    }
}

// Add Achievement
addAchievementBtn.addEventListener('click', () => {
    addModal.show();
});

addForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    const title = document.getElementById('achTitle').value;
    const type = document.getElementById('achType').value;
    const desc = document.getElementById('achDesc').value;
    const points = parseInt(document.getElementById('achPoints').value);

    const payload = {
        title,
        achievementType: type,
        description: desc,
        points,
        details: {}, // Simplified for demo
        tags: []
    };

    try {
        const res = await fetch(`${API_URL}/achievements`, {
            method: 'POST',
            headers: { 
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify(payload)
        });
        const data = await res.json();

        if (data.status === 'success') {
            addModal.hide();
            addForm.reset();
            loadAchievements();
        } else {
            alert(data.message);
        }
    } catch (err) {
        console.error(err);
        alert('Failed to add achievement');
    }
});

// Global functions for onclick
window.submitAchievement = async (id) => {
    if (!confirm('Are you sure you want to submit this achievement for verification?')) return;

    try {
        const res = await fetch(`${API_URL}/achievements/${id}/submit`, {
            method: 'POST',
            headers: { 'Authorization': `Bearer ${token}` }
        });
        const data = await res.json();
        if (data.status === 'success') {
            loadAchievements();
        } else {
            alert(data.message);
        }
    } catch (err) {
        console.error(err);
    }
};

window.viewDetail = (id) => {
    alert('Detail view not implemented in this demo');
};

init();
