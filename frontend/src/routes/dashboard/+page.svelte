<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import { authStore } from '$lib/stores/auth';

	let statistics = $state<any>(null);
	let loading = $state(true);

	onMount(async () => {
		try {
			const result = await api.getStatistics();
			if (result.status === 'success') {
				statistics = result.data;
			}
		} catch (error) {
			console.error('Failed to load statistics:', error);
		} finally {
			loading = false;
		}
	});

	let user = $derived($authStore.user);
</script>

<div class="space-y-6">
	<div>
		<h1 class="text-3xl font-bold text-gray-800">Dashboard</h1>
		<p class="text-gray-600 mt-1">Welcome back, {user?.fullName}!</p>
	</div>

	{#if loading}
		<div class="flex justify-center items-center py-12">
			<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
		</div>
	{:else if statistics}
		<!-- Statistics Cards -->
		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-{user?.role.name === 'Mahasiswa' ? '3' : '4'} gap-6">
			<!-- Total Users - Hide for Mahasiswa -->
			{#if user?.role.name !== 'Mahasiswa'}
				<div class="card p-6 bg-gradient-to-br from-blue-500 to-blue-600 text-white">
					<div class="flex items-center justify-between">
						<div>
							<p class="text-sm opacity-90">Total Users</p>
							<p class="text-3xl font-bold mt-1">{statistics.total_users || 0}</p>
						</div>
						<svg class="w-12 h-12 opacity-80" fill="currentColor" viewBox="0 0 20 20">
							<path d="M9 6a3 3 0 11-6 0 3 3 0 016 0zM17 6a3 3 0 11-6 0 3 3 0 016 0zM12.93 17c.046-.327.07-.66.07-1a6.97 6.97 0 00-1.5-4.33A5 5 0 0119 16v1h-6.07zM6 11a5 5 0 015 5v1H1v-1a5 5 0 015-5z" />
						</svg>
					</div>
				</div>
			{/if}

			<!-- Total Achievements -->
			<div class="card p-6 bg-gradient-to-br from-green-500 to-green-600 text-white">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm opacity-90">{user?.role.name === 'Mahasiswa' ? 'My Achievements' : 'Total Achievements'}</p>
						<p class="text-3xl font-bold mt-1">{statistics.total_achievements || 0}</p>
					</div>
					<svg class="w-12 h-12 opacity-80" fill="currentColor" viewBox="0 0 20 20">
						<path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
					</svg>
				</div>
			</div>

			<!-- Pending Verifications -->
			<div class="card p-6 bg-gradient-to-br from-yellow-500 to-yellow-600 text-white">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm opacity-90">Pending Verification</p>
						<p class="text-3xl font-bold mt-1">{statistics.pending_verifications || 0}</p>
					</div>
					<svg class="w-12 h-12 opacity-80" fill="currentColor" viewBox="0 0 20 20">
						<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clip-rule="evenodd" />
					</svg>
				</div>
			</div>

			<!-- Verified Achievements -->
			<div class="card p-6 bg-gradient-to-br from-purple-500 to-purple-600 text-white">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm opacity-90">Verified</p>
						<p class="text-3xl font-bold mt-1">{statistics.verified_achievements || 0}</p>
					</div>
					<svg class="w-12 h-12 opacity-80" fill="currentColor" viewBox="0 0 20 20">
						<path fill-rule="evenodd" d="M6.267 3.455a3.066 3.066 0 001.745-.723 3.066 3.066 0 013.976 0 3.066 3.066 0 001.745.723 3.066 3.066 0 012.812 2.812c.051.643.304 1.254.723 1.745a3.066 3.066 0 010 3.976 3.066 3.066 0 00-.723 1.745 3.066 3.066 0 01-2.812 2.812 3.066 3.066 0 00-1.745.723 3.066 3.066 0 01-3.976 0 3.066 3.066 0 00-1.745-.723 3.066 3.066 0 01-2.812-2.812 3.066 3.066 0 00-.723-1.745 3.066 3.066 0 010-3.976 3.066 3.066 0 00.723-1.745 3.066 3.066 0 012.812-2.812zm7.44 5.252a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
					</svg>
				</div>
			</div>
		</div>

		<!-- Role-based Quick Actions -->
		<div class="card p-6">
			<h2 class="text-xl font-semibold text-gray-800 mb-4">Quick Actions</h2>
			<div class="grid grid-cols-1 md:grid-cols-3 gap-4">
				{#if user?.role.name === 'Mahasiswa'}
					<a href="/achievements" class="btn btn-primary flex items-center justify-center gap-2">
						<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
							<path d="M10 3.5a1.5 1.5 0 013 0V4a1 1 0 001 1h3a1 1 0 011 1v3a1 1 0 01-1 1h-.5a1.5 1.5 0 000 3h.5a1 1 0 011 1v3a1 1 0 01-1 1h-3a1 1 0 01-1-1v-.5a1.5 1.5 0 00-3 0v.5a1 1 0 01-1 1H6a1 1 0 01-1-1v-3a1 1 0 00-1-1h-.5a1.5 1.5 0 010-3H4a1 1 0 001-1V6a1 1 0 011-1h3a1 1 0 001-1v-.5z" />
						</svg>
						View My Achievements
					</a>
				{/if}

				{#if user?.role.name === 'Dosen Wali'}
					<a href="/verification" class="btn btn-primary flex items-center justify-center gap-2">
						<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
							<path fill-rule="evenodd" d="M6.267 3.455a3.066 3.066 0 001.745-.723 3.066 3.066 0 013.976 0 3.066 3.066 0 001.745.723 3.066 3.066 0 012.812 2.812c.051.643.304 1.254.723 1.745a3.066 3.066 0 010 3.976 3.066 3.066 0 00-.723 1.745 3.066 3.066 0 01-2.812 2.812 3.066 3.066 0 00-1.745.723 3.066 3.066 0 01-3.976 0 3.066 3.066 0 00-1.745-.723 3.066 3.066 0 01-2.812-2.812 3.066 3.066 0 00-.723-1.745 3.066 3.066 0 010-3.976 3.066 3.066 0 00.723-1.745 3.066 3.066 0 012.812-2.812zm7.44 5.252a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
						</svg>
						Verify Achievements
					</a>
				{/if}

				{#if user?.role.name === 'Admin'}
					<a href="/achievements" class="btn btn-primary flex items-center justify-center gap-2">
						<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
							<path d="M10 3.5a1.5 1.5 0 013 0V4a1 1 0 001 1h3a1 1 0 011 1v3a1 1 0 01-1 1h-.5a1.5 1.5 0 000 3h.5a1 1 0 011 1v3a1 1 0 01-1 1h-3a1 1 0 01-1-1v-.5a1.5 1.5 0 00-3 0v.5a1 1 0 01-1 1H6a1 1 0 01-1-1v-3a1 1 0 00-1-1h-.5a1.5 1.5 0 010-3H4a1 1 0 001-1V6a1 1 0 011-1h3a1 1 0 001-1v-.5z" />
						</svg>
						Manage Achievements
					</a>
					<a href="/users" class="btn btn-secondary flex items-center justify-center gap-2">
						<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
							<path d="M9 6a3 3 0 11-6 0 3 3 0 016 0zM17 6a3 3 0 11-6 0 3 3 0 016 0zM12.93 17c.046-.327.07-.66.07-1a6.97 6.97 0 00-1.5-4.33A5 5 0 0119 16v1h-6.07zM6 11a5 5 0 015 5v1H1v-1a5 5 0 015-5z" />
						</svg>
						Manage Users
					</a>
				{/if}

				{#if user?.role.name !== 'Mahasiswa'}
					<a href="/statistics" class="btn btn-secondary flex items-center justify-center gap-2">
						<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
							<path d="M2 11a1 1 0 011-1h2a1 1 0 011 1v5a1 1 0 01-1 1H3a1 1 0 01-1-1v-5zM8 7a1 1 0 011-1h2a1 1 0 011 1v9a1 1 0 01-1 1H9a1 1 0 01-1-1V7zM14 4a1 1 0 011-1h2a1 1 0 011 1v12a1 1 0 01-1 1h-2a1 1 0 01-1-1V4z" />
						</svg>
						View Statistics
					</a>
				{/if}
			</div>
		</div>

		<!-- System Overview - Only for Admin and Dosen Wali -->
		{#if user?.role.name === 'Admin' || user?.role.name === 'Dosen Wali'}
			<div class="card p-6">
				<h2 class="text-xl font-semibold text-gray-800 mb-4">System Overview</h2>
				<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
					<div>
						<h3 class="font-medium text-gray-700 mb-2">Achievement Status</h3>
						<div class="space-y-2">
							<div class="flex justify-between items-center">
								<span class="text-gray-600">Draft</span>
								<span class="px-3 py-1 bg-gray-100 text-gray-700 rounded-full text-sm font-medium">
									{statistics.draft_achievements || 0}
								</span>
							</div>
							<div class="flex justify-between items-center">
								<span class="text-gray-600">Submitted</span>
								<span class="px-3 py-1 bg-yellow-100 text-yellow-700 rounded-full text-sm font-medium">
									{statistics.pending_verifications || 0}
								</span>
							</div>
							<div class="flex justify-between items-center">
								<span class="text-gray-600">Verified</span>
								<span class="px-3 py-1 bg-green-100 text-green-700 rounded-full text-sm font-medium">
									{statistics.verified_achievements || 0}
								</span>
							</div>
							<div class="flex justify-between items-center">
								<span class="text-gray-600">Rejected</span>
								<span class="px-3 py-1 bg-red-100 text-red-700 rounded-full text-sm font-medium">
									{statistics.rejected_achievements || 0}
								</span>
							</div>
						</div>
					</div>

					<div>
						<h3 class="font-medium text-gray-700 mb-2">User Roles</h3>
						<div class="space-y-2">
							<div class="flex justify-between items-center">
								<span class="text-gray-600">Mahasiswa</span>
								<span class="px-3 py-1 bg-blue-100 text-blue-700 rounded-full text-sm font-medium">
									{statistics.byRole?.['Mahasiswa'] || 0}
								</span>
							</div>
							<div class="flex justify-between items-center">
								<span class="text-gray-600">Dosen Wali</span>
								<span class="px-3 py-1 bg-purple-100 text-purple-700 rounded-full text-sm font-medium">
									{statistics.byRole?.['Dosen Wali'] || 0}
								</span>
							</div>
							<div class="flex justify-between items-center">
								<span class="text-gray-600">Admin</span>
								<span class="px-3 py-1 bg-indigo-100 text-indigo-700 rounded-full text-sm font-medium">
									{statistics.byRole?.['Admin'] || 0}
								</span>
							</div>
						</div>
					</div>
				</div>
			</div>
		{/if}
	{/if}
</div>
