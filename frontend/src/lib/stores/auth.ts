import { writable } from 'svelte/store';
import { browser } from '$app/environment';

// Types
export interface User {
	id: number;
	username: string;
	email: string;
	fullName: string;
	role: {
		id: number;
		name: string;
	};
}

export interface AuthState {
	token: string | null;
	user: User | null;
	isAuthenticated: boolean;
}

// Initialize from localStorage if in browser
const initialState: AuthState = browser
	? {
			token: localStorage.getItem('token'),
			user: JSON.parse(localStorage.getItem('user') || 'null'),
			isAuthenticated: !!localStorage.getItem('token')
	  }
	: {
			token: null,
			user: null,
			isAuthenticated: false
	  };

function createAuthStore() {
	const { subscribe, set, update } = writable<AuthState>(initialState);

	return {
		subscribe,
		login: (token: string, user: User) => {
			if (browser) {
				localStorage.setItem('token', token);
				localStorage.setItem('user', JSON.stringify(user));
			}
			set({ token, user, isAuthenticated: true });
		},
		logout: () => {
			if (browser) {
				localStorage.removeItem('token');
				localStorage.removeItem('user');
			}
			set({ token: null, user: null, isAuthenticated: false });
		},
		updateUser: (user: User) => {
			if (browser) {
				localStorage.setItem('user', JSON.stringify(user));
			}
			update((state) => ({ ...state, user }));
		}
	};
}

export const authStore = createAuthStore();
