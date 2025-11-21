import { authStore } from '$lib/stores/auth';
import { get } from 'svelte/store';

// Use relative path for Vercel deployment
const API_URL = '/api/v1';

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

	const response = await fetch(`${API_URL}${endpoint}`, {
		...options,
		headers
	});

	return response.json();
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
			method: 'PATCH'
		}),

	// Verification
	getPendingVerifications: () => fetchApi('/verification/pending'),

	verifyAchievement: (id: string) =>
		fetchApi(`/verification/${id}/verify`, {
			method: 'PATCH'
		}),

	rejectAchievement: (id: string, note: string) =>
		fetchApi(`/verification/${id}/reject`, {
			method: 'PATCH',
			body: JSON.stringify({ note })
		}),

	// Users
	getUsers: () => fetchApi('/users'),

	toggleUserStatus: (id: number) =>
		fetchApi(`/users/${id}/toggle-status`, {
			method: 'PATCH'
		}),

	// Statistics
	getStatistics: () => fetchApi('/reports/statistics')
};
