<script lang="ts">
	import { api } from '$lib/api';
	import { authStore } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';

	let isLogin = $state(true);
	let loading = $state(false);
	let error = $state('');

	// Login form
	let loginUsername = $state('');
	let loginPassword = $state('');

	// Register form
	let regUsername = $state('');
	let regEmail = $state('');
	let regFullName = $state('');
	let regPassword = $state('');
	let regRole = $state('Mahasiswa');

	let isAuthenticated = $derived($authStore.isAuthenticated);

	onMount(() => {
		if (isAuthenticated) {
			goto('/dashboard');
		}
	});

	async function handleLogin(e: Event) {
		e.preventDefault();
		loading = true;
		error = '';

		try {
			const result = await api.login(loginUsername, loginPassword);
			
			if (result.status === 'success' && result.data) {
				authStore.login(result.data.token, result.data.user);
				goto('/dashboard');
			} else {
				error = result.message || 'Login failed';
			}
		} catch (err) {
			error = 'Network error. Please try again.';
		} finally {
			loading = false;
		}
	}

	async function handleRegister(e: Event) {
		e.preventDefault();
		loading = true;
		error = '';

		try {
			const result = await api.register({
				username: regUsername,
				email: regEmail,
				fullName: regFullName,
				password: regPassword,
				roleName: regRole
			});

			if (result.status === 'success') {
				alert('Registration successful! Please login.');
				isLogin = true;
				// Reset form
				regUsername = '';
				regEmail = '';
				regFullName = '';
				regPassword = '';
				regRole = 'Mahasiswa';
			} else {
				error = result.message || 'Registration failed';
			}
		} catch (err) {
			error = 'Network error. Please try again.';
		} finally {
			loading = false;
		}
	}
</script>

<div class="min-h-screen flex items-center justify-center px-4 bg-gradient-to-br from-primary-50 to-primary-100">
	<div class="card max-w-md w-full p-8 shadow-xl">
		<div class="text-center mb-8">
			<div class="w-16 h-16 bg-gradient-to-br from-primary-500 to-primary-600 rounded-full flex items-center justify-center mx-auto mb-4">
				<svg class="w-10 h-10 text-white" fill="currentColor" viewBox="0 0 20 20">
					<path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
				</svg>
			</div>
			<h1 class="text-3xl font-bold text-gray-800">Achievement System</h1>
			<p class="text-gray-600 mt-2">Manage your academic achievements</p>
		</div>

		<!-- Tab Switcher -->
		<div class="flex border-b border-gray-200 mb-6">
			<button
				class="flex-1 py-2 text-center font-medium transition {isLogin ? 'border-b-2 border-primary-500 text-primary-600' : 'text-gray-500'}"
				onclick={() => { isLogin = true; error = ''; }}
			>
				Login
			</button>
			<button
				class="flex-1 py-2 text-center font-medium transition {!isLogin ? 'border-b-2 border-primary-500 text-primary-600' : 'text-gray-500'}"
				onclick={() => { isLogin = false; error = ''; }}
			>
				Register
			</button>
		</div>

		{#if error}
			<div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
				{error}
			</div>
		{/if}

		{#if isLogin}
			<!-- Login Form -->
			<form onsubmit={handleLogin} class="space-y-4">
				<div>
					<label class="block text-sm font-medium text-gray-700 mb-1">Username</label>
					<input
						type="text"
						bind:value={loginUsername}
						required
						class="input"
						placeholder="Enter username"
					/>
				</div>
				<div>
					<label class="block text-sm font-medium text-gray-700 mb-1">Password</label>
					<input
						type="password"
						bind:value={loginPassword}
						required
						class="input"
						placeholder="Enter password"
					/>
				</div>
				<button
					type="submit"
					disabled={loading}
					class="w-full btn btn-primary"
				>
					{loading ? 'Loading...' : 'Login'}
				</button>
			</form>
		{:else}
			<!-- Register Form -->
			<form onsubmit={handleRegister} class="space-y-4">
				<div>
					<label class="block text-sm font-medium text-gray-700 mb-1">Username</label>
					<input
						type="text"
						bind:value={regUsername}
						required
						class="input"
						placeholder="Choose username"
					/>
				</div>
				<div>
					<label class="block text-sm font-medium text-gray-700 mb-1">Email</label>
					<input
						type="email"
						bind:value={regEmail}
						required
						class="input"
						placeholder="your.email@example.com"
					/>
				</div>
				<div>
					<label class="block text-sm font-medium text-gray-700 mb-1">Full Name</label>
					<input
						type="text"
						bind:value={regFullName}
						required
						class="input"
						placeholder="Your full name"
					/>
				</div>
				<div>
					<label class="block text-sm font-medium text-gray-700 mb-1">Password</label>
					<input
						type="password"
						bind:value={regPassword}
						required
						class="input"
						placeholder="Choose password"
					/>
				</div>
				<div>
					<label class="block text-sm font-medium text-gray-700 mb-1">Role</label>
					<select bind:value={regRole} class="input">
						<option value="Mahasiswa">Mahasiswa</option>
						<option value="Dosen Wali">Dosen Wali</option>
						<option value="Admin">Admin</option>
					</select>
				</div>
				<button
					type="submit"
					disabled={loading}
					class="w-full btn btn-primary"
				>
					{loading ? 'Loading...' : 'Register'}
				</button>
			</form>
		{/if}
	</div>
</div>
