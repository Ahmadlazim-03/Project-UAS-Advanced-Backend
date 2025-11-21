<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import { authStore } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';

	let { children } = $props();

	let sidebarOpen = $state(true);

	function handleLogout() {
		authStore.logout();
		goto('/login');
	}

	function toggleSidebar() {
		sidebarOpen = !sidebarOpen;
	}

	let user = $derived($authStore.user);
	let isAuthenticated = $derived($authStore.isAuthenticated);
	let currentPath = $derived($page.url.pathname);

	function isActive(path: string) {
		return currentPath === path;
	}
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<div class="min-h-screen bg-gray-50">
	{#if isAuthenticated && user}
		<!-- Top Header -->
		<header class="bg-white shadow-sm fixed top-0 left-0 right-0 z-30">
			<div class="flex items-center justify-between h-16 px-4">
				<div class="flex items-center gap-4">
					<button onclick={toggleSidebar} class="p-2 rounded-lg hover:bg-gray-100 lg:hidden">
						<svg class="w-6 h-6" fill="currentColor" viewBox="0 0 20 20">
							<path fill-rule="evenodd" d="M3 5a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zM3 10a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zM3 15a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z" clip-rule="evenodd" />
						</svg>
					</button>
					<div class="flex items-center gap-2">
						<svg class="w-8 h-8 text-blue-600" fill="currentColor" viewBox="0 0 20 20">
							<path d="M10 3.5a1.5 1.5 0 013 0V4a1 1 0 001 1h3a1 1 0 011 1v3a1 1 0 01-1 1h-.5a1.5 1.5 0 000 3h.5a1 1 0 011 1v3a1 1 0 01-1 1h-3a1 1 0 01-1-1v-.5a1.5 1.5 0 00-3 0v.5a1 1 0 01-1 1H6a1 1 0 01-1-1v-3a1 1 0 00-1-1h-.5a1.5 1.5 0 010-3H4a1 1 0 001-1V6a1 1 0 011-1h3a1 1 0 001-1v-.5z" />
						</svg>
						<h1 class="text-xl font-bold text-gray-800">Achievement System</h1>
					</div>
				</div>

				<div class="flex items-center gap-4">
					<div class="text-right">
						<p class="text-sm font-medium text-gray-800">{user.fullName}</p>
						<p class="text-xs text-gray-500">{user.role.name}</p>
					</div>
					<button onclick={handleLogout} class="bg-red-500 hover:bg-red-600 text-white px-4 py-2 rounded-lg text-sm font-medium transition">
						Logout
					</button>
				</div>
			</div>
		</header>

		<!-- Sidebar -->
		<aside class="fixed left-0 top-16 bottom-0 w-64 bg-gradient-to-b from-blue-600 to-blue-800 text-white shadow-xl z-20 transition-transform duration-300 {sidebarOpen ? 'translate-x-0' : '-translate-x-full'} lg:translate-x-0">
			<nav class="p-4 space-y-2">
				<!-- Role Debug Info -->
				<div class="text-xs text-blue-200 mb-4 p-2 bg-blue-900 rounded">
					Role: {user?.role?.name || 'Unknown'}
				</div>

				<!-- Dashboard - Always visible -->
				<a href="/dashboard" class="nav-item {isActive('/dashboard') ? 'active' : ''}">
					<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
						<path d="M3 4a1 1 0 011-1h12a1 1 0 011 1v2a1 1 0 01-1 1H4a1 1 0 01-1-1V4zM3 10a1 1 0 011-1h6a1 1 0 011 1v6a1 1 0 01-1 1H4a1 1 0 01-1-1v-6zM14 9a1 1 0 00-1 1v6a1 1 0 001 1h2a1 1 0 001-1v-6a1 1 0 00-1-1h-2z" />
					</svg>
					<span>Dashboard</span>
				</a>

				{#if user && user.role && user.role.name === 'Mahasiswa'}
					<!-- Mahasiswa Menu -->
					<a href="/achievements" class="nav-item {isActive('/achievements') ? 'active' : ''}">
						<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
							<path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
						</svg>
						<span>My Achievements</span>
					</a>
				{/if}

				{#if user && user.role && user.role.name === 'Dosen Wali'}
					<!-- Dosen Wali Menu -->
					<a href="/verification" class="nav-item {isActive('/verification') ? 'active' : ''}">
						<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
							<path fill-rule="evenodd" d="M6.267 3.455a3.066 3.066 0 001.745-.723 3.066 3.066 0 013.976 0 3.066 3.066 0 001.745.723 3.066 3.066 0 012.812 2.812c.051.643.304 1.254.723 1.745a3.066 3.066 0 010 3.976 3.066 3.066 0 00-.723 1.745 3.066 3.066 0 01-2.812 2.812 3.066 3.066 0 00-1.745.723 3.066 3.066 0 01-3.976 0 3.066 3.066 0 00-1.745-.723 3.066 3.066 0 01-2.812-2.812 3.066 3.066 0 00-.723-1.745 3.066 3.066 0 010-3.976 3.066 3.066 0 00.723-1.745 3.066 3.066 0 012.812-2.812zm7.44 5.252a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
						</svg>
						<span>Verification</span>
					</a>
					<a href="/statistics" class="nav-item {isActive('/statistics') ? 'active' : ''}">
						<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
							<path d="M2 11a1 1 0 011-1h2a1 1 0 011 1v5a1 1 0 01-1 1H3a1 1 0 01-1-1v-5zM8 7a1 1 0 011-1h2a1 1 0 011 1v9a1 1 0 01-1 1H9a1 1 0 01-1-1V7zM14 4a1 1 0 011-1h2a1 1 0 011 1v12a1 1 0 01-1 1h-2a1 1 0 01-1-1V4z" />
						</svg>
						<span>Reports</span>
					</a>
				{/if}

				{#if user && user.role && user.role.name === 'Admin'}
					<!-- Admin Menu -->
					<a href="/verification" class="nav-item {isActive('/verification') ? 'active' : ''}">
						<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
							<path fill-rule="evenodd" d="M6.267 3.455a3.066 3.066 0 001.745-.723 3.066 3.066 0 013.976 0 3.066 3.066 0 001.745.723 3.066 3.066 0 012.812 2.812c.051.643.304 1.254.723 1.745a3.066 3.066 0 010 3.976 3.066 3.066 0 00-.723 1.745 3.066 3.066 0 01-2.812 2.812 3.066 3.066 0 00-1.745.723 3.066 3.066 0 01-3.976 0 3.066 3.066 0 00-1.745-.723 3.066 3.066 0 01-2.812-2.812 3.066 3.066 0 00-.723-1.745 3.066 3.066 0 010-3.976 3.066 3.066 0 00.723-1.745 3.066 3.066 0 012.812-2.812zm7.44 5.252a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
						</svg>
						<span>Verification</span>
					</a>
					<a href="/users" class="nav-item {isActive('/users') ? 'active' : ''}">
						<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
							<path d="M9 6a3 3 0 11-6 0 3 3 0 016 0zM17 6a3 3 0 11-6 0 3 3 0 016 0zM12.93 17c.046-.327.07-.66.07-1a6.97 6.97 0 00-1.5-4.33A5 5 0 0119 16v1h-6.07zM6 11a5 5 0 015 5v1H1v-1a5 5 0 015-5z" />
						</svg>
						<span>Users</span>
					</a>
					<a href="/students" class="nav-item {isActive('/students') ? 'active' : ''}">
						<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
							<path d="M10.394 2.08a1 1 0 00-.788 0l-7 3a1 1 0 000 1.84L5.25 8.051a.999.999 0 01.356-.257l4-1.714a1 1 0 11.788 1.838L7.667 9.088l1.94.831a1 1 0 00.787 0l7-3a1 1 0 000-1.838l-7-3zM3.31 9.397L5 10.12v4.102a8.969 8.969 0 00-1.05-.174 1 1 0 01-.89-.89 11.115 11.115 0 01.25-3.762zM9.3 16.573A9.026 9.026 0 007 14.935v-3.957l1.818.78a3 3 0 002.364 0l5.508-2.361a11.026 11.026 0 01.25 3.762 1 1 0 01-.89.89 8.968 8.968 0 00-5.35 2.524 1 1 0 01-1.4 0zM6 18a1 1 0 001-1v-2.065a8.935 8.935 0 00-2-.712V17a1 1 0 001 1z" />
						</svg>
						<span>Students</span>
					</a>
					<a href="/lecturers" class="nav-item {isActive('/lecturers') ? 'active' : ''}">
						<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
							<path d="M13 6a3 3 0 11-6 0 3 3 0 016 0zM18 8a2 2 0 11-4 0 2 2 0 014 0zM14 15a4 4 0 00-8 0v3h8v-3zM6 8a2 2 0 11-4 0 2 2 0 014 0zM16 18v-3a5.972 5.972 0 00-.75-2.906A3.005 3.005 0 0119 15v3h-3zM4.75 12.094A5.973 5.973 0 004 15v3H1v-3a3 3 0 013.75-2.906z" />
						</svg>
						<span>Lecturers</span>
					</a>
					<a href="/statistics" class="nav-item {isActive('/statistics') ? 'active' : ''}">
						<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
							<path d="M2 11a1 1 0 011-1h2a1 1 0 011 1v5a1 1 0 01-1 1H3a1 1 0 01-1-1v-5zM8 7a1 1 0 011-1h2a1 1 0 011 1v9a1 1 0 01-1 1H9a1 1 0 01-1-1V7zM14 4a1 1 0 011-1h2a1 1 0 011 1v12a1 1 0 01-1 1h-2a1 1 0 01-1-1V4z" />
						</svg>
						<span>Reports</span>
					</a>
				{/if}
			</nav>
		</aside>

		<!-- Main Content -->
		<main class="pt-16 lg:pl-64 min-h-screen">
			<div class="p-6">
				{@render children()}
			</div>
		</main>

		<!-- Mobile Overlay -->
		{#if sidebarOpen}
			<div class="fixed inset-0 bg-black bg-opacity-50 z-10 lg:hidden" onclick={toggleSidebar}></div>
		{/if}
	{:else}
		{@render children()}
	{/if}
</div>

<style>
	.nav-item {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 0.75rem 1rem;
		border-radius: 0.5rem;
		font-weight: 500;
		transition: all 0.2s;
		color: rgba(255, 255, 255, 0.9);
	}

	.nav-item:hover {
		background-color: rgba(255, 255, 255, 0.1);
		color: white;
	}

	.nav-item.active {
		background-color: rgba(255, 255, 255, 0.2);
		color: white;
		font-weight: 600;
	}
</style>
