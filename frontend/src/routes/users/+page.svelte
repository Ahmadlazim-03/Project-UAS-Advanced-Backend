<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';

	let users = $state<any[]>([]);
	let loading = $state(true);
	let filter = $state('all');

	onMount(async () => {
		await loadUsers();
	});

	async function loadUsers() {
		loading = true;
		try {
			const result = await api.getUsers();
			if (result.status === 'success') {
				users = result.data || [];
			}
		} catch (error) {
			console.error('Failed to load users:', error);
		} finally {
			loading = false;
		}
	}

	async function toggleUserStatus(userId: string) {
		try {
			await api.toggleUserStatus(userId);
			await loadUsers();
		} catch (error) {
			console.error('Failed to toggle user status:', error);
			alert('Failed to update user status');
		}
	}

	let filteredUsers = $derived(users.filter((user) => {
		if (filter === 'all') return true;
		if (filter === 'active') return user.is_active;
		if (filter === 'inactive') return !user.is_active;
		return user.role?.name === filter;
	}));

	function getRoleBadge(roleName: string) {
		const badges: Record<string, { class: string; text: string }> = {
			Admin: { class: 'bg-indigo-100 text-indigo-700', text: 'Admin' },
			'Dosen Wali': { class: 'bg-purple-100 text-purple-700', text: 'Dosen Wali' },
			Mahasiswa: { class: 'bg-blue-100 text-blue-700', text: 'Mahasiswa' }
		};
		return badges[roleName] || { class: 'bg-gray-100 text-gray-700', text: roleName };
	}
</script>

<div class="space-y-6">
	<div>
		<h1 class="text-3xl font-bold text-gray-800">User Management</h1>
		<p class="text-gray-600 mt-1">Manage users and their access to the system</p>
	</div>

	<!-- Filters -->
	<div class="card p-4">
		<div class="flex flex-wrap gap-2">
			<button
				onclick={() => (filter = 'all')}
				class="px-4 py-2 rounded-lg transition-colors {filter === 'all'
					? 'bg-primary-600 text-white'
					: 'bg-gray-100 text-gray-700 hover:bg-gray-200'}"
			>
				All Users ({users.length})
			</button>
			<button
				onclick={() => (filter = 'Mahasiswa')}
				class="px-4 py-2 rounded-lg transition-colors {filter === 'Mahasiswa'
					? 'bg-primary-600 text-white'
					: 'bg-gray-100 text-gray-700 hover:bg-gray-200'}"
			>
				Mahasiswa ({users.filter((u) => u.role?.name === 'Mahasiswa').length})
			</button>
			<button
				onclick={() => (filter = 'Dosen Wali')}
				class="px-4 py-2 rounded-lg transition-colors {filter === 'Dosen Wali'
					? 'bg-primary-600 text-white'
					: 'bg-gray-100 text-gray-700 hover:bg-gray-200'}"
			>
				Dosen Wali ({users.filter((u) => u.role?.name === 'Dosen Wali').length})
			</button>
			<button
				onclick={() => (filter = 'Admin')}
				class="px-4 py-2 rounded-lg transition-colors {filter === 'Admin'
					? 'bg-primary-600 text-white'
					: 'bg-gray-100 text-gray-700 hover:bg-gray-200'}"
			>
				Admin ({users.filter((u) => u.role?.name === 'Admin').length})
			</button>
			<button
				onclick={() => (filter = 'active')}
				class="px-4 py-2 rounded-lg transition-colors {filter === 'active'
					? 'bg-primary-600 text-white'
					: 'bg-gray-100 text-gray-700 hover:bg-gray-200'}"
			>
				Active ({users.filter((u) => u.is_active).length})
			</button>
			<button
				onclick={() => (filter = 'inactive')}
				class="px-4 py-2 rounded-lg transition-colors {filter === 'inactive'
					? 'bg-primary-600 text-white'
					: 'bg-gray-100 text-gray-700 hover:bg-gray-200'}"
			>
				Inactive ({users.filter((u) => !u.is_active).length})
			</button>
		</div>
	</div>

	{#if loading}
		<div class="flex justify-center items-center py-12">
			<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
		</div>
	{:else if filteredUsers.length === 0}
		<div class="card p-12 text-center">
			<svg
				class="w-16 h-16 mx-auto text-gray-400 mb-4"
				fill="none"
				stroke="currentColor"
				viewBox="0 0 24 24"
			>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					stroke-width="2"
					d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"
				/>
			</svg>
			<h3 class="text-lg font-medium text-gray-900 mb-2">No users found</h3>
			<p class="text-gray-600">Try adjusting your filters</p>
		</div>
	{:else}
		<!-- Users Table -->
		<div class="card overflow-hidden">
			<div class="overflow-x-auto">
				<table class="min-w-full divide-y divide-gray-200">
					<thead class="bg-gray-50">
						<tr>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
								User
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
								Role
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
								Details
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
								Status
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
								Actions
							</th>
						</tr>
					</thead>
					<tbody class="bg-white divide-y divide-gray-200">
						{#each filteredUsers as user}
							<tr class="hover:bg-gray-50 transition-colors">
								<td class="px-6 py-4 whitespace-nowrap">
									<div class="flex items-center">
										<div class="flex-shrink-0 h-10 w-10">
											<div
												class="h-10 w-10 rounded-full bg-gradient-to-br from-primary-400 to-primary-600 flex items-center justify-center text-white font-semibold"
											>
												{user.full_name?.charAt(0).toUpperCase() || user.username?.charAt(0).toUpperCase() || 'U'}
											</div>
										</div>
										<div class="ml-4">
											<div class="text-sm font-medium text-gray-900">{user.full_name || 'N/A'}</div>
											<div class="text-sm text-gray-500">{user.email || user.username}</div>
										</div>
								</div>
							</td>
							<td class="px-6 py-4 whitespace-nowrap">
								{#snippet roleBadge()}
									{@const badge = getRoleBadge(user.role?.name || 'Unknown')}
									<span class="px-3 py-1 inline-flex text-xs leading-5 font-semibold rounded-full {badge.class}">
										{badge.text}
									</span>
								{/snippet}
								{@render roleBadge()}
							</td>
							<td class="px-6 py-4">
									<div class="text-sm text-gray-900">
										{#if user.role?.name === 'Mahasiswa' && user.student}
											<div>NIM: {user.student.nim || 'N/A'}</div>
											{#if user.student.major}
												<div class="text-gray-500">Major: {user.student.major}</div>
											{/if}
										{:else if user.role?.name === 'Dosen Wali' && user.lecturer}
											<div>NIP: {user.lecturer.nip || 'N/A'}</div>
											{#if user.lecturer.department}
												<div class="text-gray-500">Dept: {user.lecturer.department}</div>
											{/if}
										{:else}
											<div class="text-gray-500">-</div>
										{/if}
									</div>
								</td>
								<td class="px-6 py-4 whitespace-nowrap">
									<span
										class="px-3 py-1 inline-flex text-xs leading-5 font-semibold rounded-full {user.is_active
											? 'bg-green-100 text-green-800'
											: 'bg-red-100 text-red-800'}"
									>
										{user.is_active ? 'Active' : 'Inactive'}
									</span>
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
									<button
										onclick={() => toggleUserStatus(user.id)}
										class="text-primary-600 hover:text-primary-900 mr-3"
									>
										{user.is_active ? 'Deactivate' : 'Activate'}
									</button>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		</div>

		<!-- Summary Stats -->
		<div class="grid grid-cols-1 md:grid-cols-4 gap-4">
			<div class="card p-4">
				<div class="text-sm text-gray-600">Total Users</div>
				<div class="text-2xl font-bold text-gray-800 mt-1">{users.length}</div>
			</div>
			<div class="card p-4">
				<div class="text-sm text-gray-600">Active Users</div>
				<div class="text-2xl font-bold text-green-600 mt-1">
					{users.filter((u) => u.is_active).length}
				</div>
			</div>
			<div class="card p-4">
				<div class="text-sm text-gray-600">Inactive Users</div>
				<div class="text-2xl font-bold text-red-600 mt-1">
					{users.filter((u) => !u.is_active).length}
				</div>
			</div>
			<div class="card p-4">
				<div class="text-sm text-gray-600">Students</div>
				<div class="text-2xl font-bold text-blue-600 mt-1">
					{users.filter((u) => u.role?.name === 'Mahasiswa').length}
				</div>
			</div>
		</div>
	{/if}
</div>
