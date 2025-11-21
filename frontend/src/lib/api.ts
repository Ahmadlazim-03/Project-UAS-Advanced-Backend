import { authStore } from '$lib/stores/auth';
import { get } from 'svelte/store';

// Detect API URL based on environment
function getApiUrl(): string {
	// Check for explicit environment variable first
	if (import.meta.env.VITE_API_URL) {
		return import.meta.env.VITE_API_URL;
	}

	// If in browser, detect from window.location
	if (typeof window !== 'undefined') {
		const hostname = window.location.hostname;
		
		// GitHub Codespaces detection
		if (hostname.includes('.app.github.dev')) {
			// Extract codespace name and construct backend URL
			const match = hostname.match(/([^.]+)-\d+\.app\.github\.dev/);
			if (match) {
				const codespaceName = match[1];
				return `https://${codespaceName}-3000.app.github.dev/api/v1`;
			}
		}
		
		// Gitpod detection
		if (hostname.includes('.gitpod.io')) {
			const match = hostname.match(/^(\d+)-(.+)\.gitpod\.io$/);
			if (match) {
				const workspaceId = match[2];
				return `https://3000-${workspaceId}.gitpod.io/api/v1`;
			}
		}
	}

	// Local development fallback
	if (import.meta.env.DEV) {
		return 'http://localhost:3000/api/v1';
	}

	// Production (Vercel)
	return '/api/v1';
}

const API_URL = getApiUrl();
console.log('🌐 API_URL:', API_URL, 'DEV mode:', import.meta.env.DEV, 'Window location:', typeof window !== 'undefined' ? window.location.href : 'SSR');

interface ApiResponse<T = any> {
	status: string;
	message?: string;
	data?: T;
}

async function fetchApi<T = any>(
	endpoint: string,
	options: RequestInit = {}
): Promise<ApiResponse<T>> {
	const auth = get(authStore);
	const headers: Record<string, string> = {
		'Content-Type': 'application/json'
	};

	// Merge additional headers if provided
	if (options.headers) {
		const customHeaders = options.headers as Record<string, string>;
		Object.assign(headers, customHeaders);
	}

	// Add authorization header if token exists
	if (auth.token) {
		headers['Authorization'] = `Bearer ${auth.token}`;
	}

	try {
		const url = `${API_URL}${endpoint}`;
		console.log('Fetching:', url, 'Method:', options.method || 'GET');
		
		const response = await fetch(url, {
			...options,
			headers
		});

		const data = await response.json();
		console.log('Response:', data);
		return data;
	} catch (error) {
		console.error('API Error:', error);
		throw error;
	}
}

export const api = {
	// Auth
	login: (username: string, password: string) =>
		fetchApi('/auth/login', {
			method: 'POST',
			body: JSON.stringify({ username, password })
		}),

	register: (data: {
		username: string;
		email: string;
		fullName: string;
		password: string;
		roleName: string;
	}) =>
		fetchApi('/auth/register', {
			method: 'POST',
			body: JSON.stringify(data)
		}),

	// Achievements
	getAchievements: () => fetchApi('/achievements'),

	getAchievement: (id: string) => fetchApi(`/achievements/${id}`),

	createAchievement: (data: any) =>
		fetchApi('/achievements', {
			method: 'POST',
			body: JSON.stringify(data)
		}),

	updateAchievement: (id: string, data: any) =>
		fetchApi(`/achievements/${id}`, {
			method: 'PUT',
			body: JSON.stringify(data)
		}),

	deleteAchievement: (id: string) =>
		fetchApi(`/achievements/${id}`, {
			method: 'DELETE'
		}),

	submitAchievement: (id: string) =>
		fetchApi(`/achievements/${id}/submit`, {
			method: 'POST'
		}),

	// Verification
	getPendingVerifications: () => fetchApi('/verification/pending'),

	verifyAchievement: (id: string) =>
		fetchApi(`/achievements/${id}/verify`, {
			method: 'POST'
		}),

	rejectAchievement: (id: string, note: string) =>
		fetchApi(`/achievements/${id}/reject`, {
			method: 'POST',
			body: JSON.stringify({ note })
		}),

	// Achievement History & Attachments
	getAchievementHistory: (id: string) => fetchApi(`/achievements/${id}/history`),

	uploadAttachment: (id: string, file: File) => {
		const formData = new FormData();
		formData.append('file', file);
		const auth = get(authStore);
		return fetch(`${API_URL}/achievements/${id}/attachments`, {
			method: 'POST',
			headers: {
				'Authorization': `Bearer ${auth.token}`
			},
			body: formData
		}).then(res => res.json());
	},

	// Users
	getUsers: () => fetchApi('/users'),

	getUser: (id: string) => fetchApi(`/users/${id}`),

	createUser: (data: any) =>
		fetchApi('/users', {
			method: 'POST',
			body: JSON.stringify(data)
		}),

	updateUser: (id: string, data: any) =>
		fetchApi(`/users/${id}`, {
			method: 'PUT',
			body: JSON.stringify(data)
		}),

	deleteUser: (id: string) =>
		fetchApi(`/users/${id}`, {
			method: 'DELETE'
		}),

	updateUserRole: (id: string, roleName: string) =>
		fetchApi(`/users/${id}/role`, {
			method: 'PUT',
			body: JSON.stringify({ roleName })
		}),

	toggleUserStatus: (id: string) =>
		fetchApi(`/users/${id}/toggle-status`, {
			method: 'PATCH'
		}),

	// Students
	getStudents: () => fetchApi('/students'),

	getStudent: (id: string) => fetchApi(`/students/${id}`),

	createStudent: (data: any) =>
		fetchApi('/students', {
			method: 'POST',
			body: JSON.stringify(data)
		}),

	updateStudent: (id: string, data: any) =>
		fetchApi(`/students/${id}`, {
			method: 'PUT',
			body: JSON.stringify(data)
		}),

	deleteStudent: (id: string) =>
		fetchApi(`/students/${id}`, {
			method: 'DELETE'
		}),

	getStudentAchievements: (id: string) => fetchApi(`/students/${id}/achievements`),

	updateStudentAdvisor: (id: string, advisorId: string) =>
		fetchApi(`/students/${id}/advisor`, {
			method: 'PUT',
			body: JSON.stringify({ advisorId })
		}),

	// Lecturers
	getLecturers: () => fetchApi('/lecturers'),

	getLecturer: (id: string) => fetchApi(`/lecturers/${id}`),

	createLecturer: (data: any) =>
		fetchApi('/lecturers', {
			method: 'POST',
			body: JSON.stringify(data)
		}),

	updateLecturer: (id: string, data: any) =>
		fetchApi(`/lecturers/${id}`, {
			method: 'PUT',
			body: JSON.stringify(data)
		}),

	deleteLecturer: (id: string) =>
		fetchApi(`/lecturers/${id}`, {
			method: 'DELETE'
		}),

	getLecturerAdvisees: (id: string) => fetchApi(`/lecturers/${id}/advisees`),

	// Statistics & Reports
	getStatistics: () => fetchApi('/reports/statistics'),

	getStudentReport: (id: string) => fetchApi(`/reports/student/${id}`)
};
