<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import { authStore } from '$lib/stores/auth';
	import { goto } from '$app/navigation';

	let { children } = $props();

	function handleLogout() {
		authStore.logout();
		goto('/login');
	}

	let user = $derived($authStore.user);
	let isAuthenticated = $derived($authStore.isAuthenticated);
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<div class="min-h-screen">
	{#if isAuthenticated && user}
		<!-- Navbar -->
		<nav class="bg-primary-600 shadow-lg">
			<div class="container mx-auto px-4">
				<div class="flex items-center justify-between h-16">
					<div class="flex items-center space-x-8">
						<a href="/dashboard" class="text-white text-xl font-bold flex items-center gap-2">
							<svg class="w-6 h-6" fill="currentColor" viewBox="0 0 20 20">
								<path d="M10 3.5a1.5 1.5 0 013 0V4a1 1 0 001 1h3a1 1 0 011 1v3a1 1 0 01-1 1h-.5a1.5 1.5 0 000 3h.5a1 1 0 011 1v3a1 1 0 01-1 1h-3a1 1 0 01-1-1v-.5a1.5 1.5 0 00-3 0v.5a1 1 0 01-1 1H6a1 1 0 01-1-1v-3a1 1 0 00-1-1h-.5a1.5 1.5 0 010-3H4a1 1 0 001-1V6a1 1 0 011-1h3a1 1 0 001-1v-.5z" />
							</svg>
							Achievement System
						</a>

						<div class="hidden md:flex space-x-4">
							<a href="/dashboard" class="text-white hover:bg-primary-700 px-3 py-2 rounded-lg transition">
								Dashboard
							</a>

							{#if user.role.name === 'Mahasiswa'}
								<a href="/achievements" class="text-white hover:bg-primary-700 px-3 py-2 rounded-lg transition">
									My Achievements
								</a>
							{/if}

							{#if user.role.name === 'Dosen Wali'}
								<a href="/verification" class="text-white hover:bg-primary-700 px-3 py-2 rounded-lg transition">
									Verify
								</a>
							{/if}

							{#if user.role.name === 'Admin'}
								<a href="/achievements" class="text-white hover:bg-primary-700 px-3 py-2 rounded-lg transition">
									Achievements
								</a>
								<a href="/verification" class="text-white hover:bg-primary-700 px-3 py-2 rounded-lg transition">
									Verify
								</a>
								<a href="/users" class="text-white hover:bg-primary-700 px-3 py-2 rounded-lg transition">
									Users
								</a>
							{/if}

							<a href="/statistics" class="text-white hover:bg-primary-700 px-3 py-2 rounded-lg transition">
								Statistics
							</a>
						</div>
					</div>

					<div class="flex items-center space-x-4">
						<span class="text-white text-sm">
							{user.fullName} ({user.role.name})
						</span>
						<button onclick={handleLogout} class="btn bg-red-500 hover:bg-red-600 text-white text-sm">
							Logout
						</button>
					</div>
				</div>
			</div>
		</nav>

		<!-- Main Content -->
		<main class="container mx-auto px-4 py-8">
			{@render children()}
		</main>
	{:else}
		{@render children()}
	{/if}
</div>
