<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';

	let lecturers = $state<any[]>([]);
	let users = $state<any[]>([]);
	let loading = $state(true);
	let selectedLecturer = $state<any>(null);
	let advisees = $state<any[]>([]);
	let showAdviseesModal = $state(false);
	let showLecturerModal = $state(false);
	let loadingAdvisees = $state(false);
	let isEditMode = $state(false);

	// Form fields
	let formData = $state({
		user_id: '',
		lecturer_id: '',
		department: ''
	});

	onMount(async () => {
		await loadData();
	});

	async function loadData() {
		loading = true;
		try {
			const [lecturersResult, usersResult] = await Promise.all([
				api.getLecturers(),
				api.getUsers()
			]);

			if (lecturersResult.status === 'success') {
				lecturers = lecturersResult.data || [];
			}

			if (usersResult.status === 'success') {
				// Filter only users with Dosen Wali role
				users = (usersResult.data || []).filter((u: any) => u.role?.name === 'Dosen Wali');
			}
		} catch (error) {
			console.error('Failed to load data:', error);
		} finally {
			loading = false;
		}
	}

	function openCreateModal() {
		isEditMode = false;
		selectedLecturer = null;
		formData = {
			user_id: '',
			lecturer_id: '',
			department: ''
		};
		showLecturerModal = true;
	}

	function openEditModal(lecturer: any) {
		isEditMode = true;
		selectedLecturer = lecturer;
		formData = {
			user_id: lecturer.user_id || '',
			lecturer_id: lecturer.lecturer_id || lecturer.nip || '',
			department: lecturer.department || ''
		};
		showLecturerModal = true;
	}

	async function saveLecturer() {
		try {
			if (isEditMode && selectedLecturer) {
				await api.updateLecturer(selectedLecturer.id, formData);
			} else {
				await api.createLecturer(formData);
			}
			showLecturerModal = false;
			await loadData();
		} catch (error) {
			console.error('Failed to save lecturer:', error);
			alert('Failed to save lecturer');
		}
	}

	async function deleteLecturer(lecturer: any) {
		if (!confirm(`Are you sure you want to delete ${lecturer.user?.full_name || 'this lecturer'}?`)) {
			return;
		}

		try {
			await api.deleteLecturer(lecturer.id);
			await loadData();
		} catch (error) {
			console.error('Failed to delete lecturer:', error);
			alert('Failed to delete lecturer');
		}
	}

	async function viewAdvisees(lecturer: any) {
		selectedLecturer = lecturer;
		loadingAdvisees = true;
		showAdviseesModal = true;

		try {
			const result = await api.getLecturerAdvisees(lecturer.id);
			if (result.status === 'success') {
				advisees = result.data || [];
			}
		} catch (error) {
			console.error('Failed to load advisees:', error);
			advisees = [];
		} finally {
			loadingAdvisees = false;
		}
	}
</script>

<div class="space-y-6">
	<div class="flex justify-between items-center">
		<div>
			<h1 class="text-3xl font-bold text-gray-800">Lecturers Management</h1>
			<p class="text-gray-600 mt-1">Manage lecturers (Dosen Wali) and view their advisees</p>
		</div>
		<button onclick={openCreateModal} class="btn btn-primary px-6 py-3">
			<svg class="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 20 20">
				<path fill-rule="evenodd" d="M10 3a1 1 0 011 1v5h5a1 1 0 110 2h-5v5a1 1 0 11-2 0v-5H4a1 1 0 110-2h5V4a1 1 0 011-1z" clip-rule="evenodd" />
			</svg>
			Add Lecturer
		</button>
	</div>

	{#if loading}
		<div class="flex justify-center items-center py-12">
			<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
		</div>
	{:else if lecturers.length === 0}
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
					d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"
				/>
			</svg>
			<h3 class="text-lg font-medium text-gray-900 mb-2">No lecturers found</h3>
			<p class="text-gray-600">Lecturers will appear here once they are registered</p>
		</div>
	{:else}
		<!-- Summary Stats -->
		<div class="grid grid-cols-1 md:grid-cols-3 gap-6">
			<div class="card p-6 bg-gradient-to-br from-purple-500 to-purple-600 text-white">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm opacity-90">Total Lecturers</p>
						<p class="text-3xl font-bold mt-1">{lecturers.length}</p>
					</div>
					<svg class="w-12 h-12 opacity-80" fill="currentColor" viewBox="0 0 20 20">
						<path d="M9 6a3 3 0 11-6 0 3 3 0 016 0zM17 6a3 3 0 11-6 0 3 3 0 016 0zM12.93 17c.046-.327.07-.66.07-1a6.97 6.97 0 00-1.5-4.33A5 5 0 0119 16v1h-6.07zM6 11a5 5 0 015 5v1H1v-1a5 5 0 015-5z" />
					</svg>
				</div>
			</div>

			<div class="card p-6 bg-gradient-to-br from-blue-500 to-blue-600 text-white">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm opacity-90">Active Advisors</p>
						<p class="text-3xl font-bold mt-1">{lecturers.filter(l => l.advisees_count > 0).length}</p>
					</div>
					<svg class="w-12 h-12 opacity-80" fill="currentColor" viewBox="0 0 20 20">
						<path fill-rule="evenodd" d="M6.267 3.455a3.066 3.066 0 001.745-.723 3.066 3.066 0 013.976 0 3.066 3.066 0 001.745.723 3.066 3.066 0 012.812 2.812c.051.643.304 1.254.723 1.745a3.066 3.066 0 010 3.976 3.066 3.066 0 00-.723 1.745 3.066 3.066 0 01-2.812 2.812 3.066 3.066 0 00-1.745.723 3.066 3.066 0 01-3.976 0 3.066 3.066 0 00-1.745-.723 3.066 3.066 0 01-2.812-2.812 3.066 3.066 0 00-.723-1.745 3.066 3.066 0 010-3.976 3.066 3.066 0 00.723-1.745 3.066 3.066 0 012.812-2.812zm7.44 5.252a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
					</svg>
				</div>
			</div>

			<div class="card p-6 bg-gradient-to-br from-indigo-500 to-indigo-600 text-white">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm opacity-90">Departments</p>
						<p class="text-3xl font-bold mt-1">
							{new Set(lecturers.map(l => l.department).filter(Boolean)).size}
						</p>
					</div>
					<svg class="w-12 h-12 opacity-80" fill="currentColor" viewBox="0 0 20 20">
						<path d="M10.707 2.293a1 1 0 00-1.414 0l-7 7a1 1 0 001.414 1.414L4 10.414V17a1 1 0 001 1h2a1 1 0 001-1v-2a1 1 0 011-1h2a1 1 0 011 1v2a1 1 0 001 1h2a1 1 0 001-1v-6.586l.293.293a1 1 0 001.414-1.414l-7-7z" />
					</svg>
				</div>
			</div>
		</div>

		<!-- Lecturers Table -->
		<div class="card overflow-hidden">
			<div class="overflow-x-auto">
				<table class="min-w-full divide-y divide-gray-200">
					<thead class="bg-gray-50">
						<tr>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
								Lecturer
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
								Lecturer ID / NIP
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
								Department
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
								Advisees Count
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
								Actions
							</th>
						</tr>
					</thead>
					<tbody class="bg-white divide-y divide-gray-200">
						{#each lecturers as lecturer}
							<tr class="hover:bg-gray-50 transition-colors">
								<td class="px-6 py-4 whitespace-nowrap">
									<div class="flex items-center">
										<div class="flex-shrink-0 h-10 w-10">
											<div
												class="h-10 w-10 rounded-full bg-gradient-to-br from-purple-400 to-purple-600 flex items-center justify-center text-white font-semibold"
											>
												{lecturer.user?.full_name?.charAt(0).toUpperCase() || lecturer.full_name?.charAt(0).toUpperCase() || 'L'}
											</div>
										</div>
										<div class="ml-4">
											<div class="text-sm font-medium text-gray-900">
												{lecturer.user?.full_name || lecturer.full_name || 'N/A'}
											</div>
											<div class="text-sm text-gray-500">
												{lecturer.user?.email || 'N/A'}
											</div>
										</div>
									</div>
								</td>
								<td class="px-6 py-4 whitespace-nowrap">
									<div class="text-sm text-gray-900">{lecturer.lecturer_id || lecturer.nip || 'N/A'}</div>
								</td>
								<td class="px-6 py-4 whitespace-nowrap">
									<div class="text-sm text-gray-900">{lecturer.department || 'N/A'}</div>
								</td>
								<td class="px-6 py-4 whitespace-nowrap">
									<span class="px-3 py-1 inline-flex text-xs leading-5 font-semibold rounded-full bg-blue-100 text-blue-800">
										{lecturer.advisees_count || 0} students
									</span>
								</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm font-medium space-x-2">
								<button
									onclick={() => openEditModal(lecturer)}
									class="text-blue-600 hover:text-blue-900"
								>
									Edit
								</button>
								<button
									onclick={() => viewAdvisees(lecturer)}
									class="text-primary-600 hover:text-primary-900"
								>
									View Advisees
								</button>
								<button
									onclick={() => deleteLecturer(lecturer)}
									class="text-red-600 hover:text-red-900"
								>
									Delete
								</button>
							</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		</div>

		<!-- Department Summary -->
		<div class="card p-6">
			<h2 class="text-xl font-semibold text-gray-800 mb-4">Lecturers by Department</h2>
			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
				{#each [...new Set(lecturers.map(l => l.department).filter(Boolean))] as dept}
					{@const count = lecturers.filter(l => l.department === dept).length}
					<div class="p-4 bg-gray-50 rounded-lg">
						<div class="text-sm text-gray-600">{dept}</div>
						<div class="text-2xl font-bold text-gray-800 mt-1">{count} lecturer{count !== 1 ? 's' : ''}</div>
					</div>
				{/each}
			</div>
		</div>
	{/if}
</div>

<!-- Advisees Modal -->
{#if showAdviseesModal && selectedLecturer}
	<div class="modal-overlay" onclick={() => (showAdviseesModal = false)}>
		<div class="modal-content" onclick={(e) => e.stopPropagation()}>
			<div class="modal-header">
				<h2 class="text-2xl font-bold text-gray-800">
					Advisees of {selectedLecturer.user?.full_name || selectedLecturer.full_name}
				</h2>
				<button onclick={() => (showAdviseesModal = false)} class="text-gray-400 hover:text-gray-600">
					<svg class="w-6 h-6" fill="currentColor" viewBox="0 0 20 20">
						<path
							fill-rule="evenodd"
							d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
							clip-rule="evenodd"
						/>
					</svg>
				</button>
			</div>

			<div class="modal-body">
				{#if loadingAdvisees}
					<div class="flex justify-center items-center py-8">
						<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
					</div>
				{:else if advisees.length === 0}
					<div class="text-center py-8">
						<svg
							class="w-12 h-12 mx-auto text-gray-400 mb-3"
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
						<p class="text-gray-600">No advisees assigned yet</p>
					</div>
				{:else}
					<div class="space-y-3">
						{#each advisees as student}
							<div class="flex items-center p-3 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors">
								<div class="flex-shrink-0 h-10 w-10">
									<div
										class="h-10 w-10 rounded-full bg-gradient-to-br from-blue-400 to-blue-600 flex items-center justify-center text-white font-semibold"
									>
										{student.user?.full_name?.charAt(0).toUpperCase() || student.full_name?.charAt(0).toUpperCase() || 'S'}
									</div>
								</div>
								<div class="ml-4 flex-1">
									<div class="text-sm font-medium text-gray-900">
										{student.user?.full_name || student.full_name || 'N/A'}
									</div>
									<div class="text-sm text-gray-500">
										{student.student_id || student.nim || 'N/A'} • {student.program_study || student.major || 'N/A'}
									</div>
								</div>
							</div>
						{/each}
					</div>

					<div class="mt-4 p-4 bg-blue-50 border border-blue-200 rounded-lg">
						<div class="flex items-center gap-2">
							<svg class="w-5 h-5 text-blue-600" fill="currentColor" viewBox="0 0 20 20">
								<path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" />
							</svg>
							<span class="text-sm text-blue-800">
								Total advisees: <strong>{advisees.length}</strong>
							</span>
						</div>
					</div>
				{/if}
			</div>

			<div class="modal-footer">
				<button onclick={() => (showAdviseesModal = false)} class="btn btn-secondary">
					Close
				</button>
			</div>
		</div>
	</div>
{/if}

<!-- Create/Edit Lecturer Modal -->
{#if showLecturerModal}
	<div class="modal-overlay" onclick={() => (showLecturerModal = false)}>
		<div class="modal-content max-w-md" onclick={(e) => e.stopPropagation()}>
			<div class="modal-header">
				<h2 class="text-2xl font-bold text-gray-800">
					{isEditMode ? 'Edit Lecturer' : 'Create Lecturer'}
				</h2>
				<button onclick={() => (showLecturerModal = false)} class="text-gray-400 hover:text-gray-600">
					<svg class="w-6 h-6" fill="currentColor" viewBox="0 0 20 20">
						<path
							fill-rule="evenodd"
							d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
							clip-rule="evenodd"
						/>
					</svg>
				</button>
			</div>

			<div class="modal-body">
				<div class="space-y-4">
					<div>
						<label class="form-label">Lecturer User Account</label>
						<select bind:value={formData.user_id} class="input w-full" disabled={isEditMode}>
							<option value="">Select User</option>
							{#each users as usr}
								<option value={usr.id}>{usr.full_name} ({usr.email})</option>
							{/each}
						</select>
						<p class="text-sm text-gray-500 mt-1">Select the user account for this lecturer</p>
					</div>

					<div>
						<label class="form-label">Lecturer ID (NIP)</label>
						<input
							type="text"
							bind:value={formData.lecturer_id}
							class="input w-full"
							placeholder="e.g., 199001012020031001"
						/>
					</div>

					<div>
						<label class="form-label">Department / Faculty</label>
						<input
							type="text"
							bind:value={formData.department}
							class="input w-full"
							placeholder="e.g., Computer Science"
						/>
					</div>
				</div>
			</div>

			<div class="modal-footer">
				<button onclick={() => (showLecturerModal = false)} class="btn btn-secondary">
					Cancel
				</button>
				<button onclick={saveLecturer} class="btn btn-primary">
					{isEditMode ? 'Update' : 'Create'} Lecturer
				</button>
			</div>
		</div>
	</div>
{/if}

<style>
	.modal-overlay {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background-color: rgba(0, 0, 0, 0.5);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 50;
		padding: 1rem;
	}

	.modal-content {
		background: white;
		border-radius: 0.5rem;
		max-width: 48rem;
		width: 100%;
		max-height: 90vh;
		overflow-y: auto;
		box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
	}

	.modal-content.max-w-md {
		max-width: 28rem;
	}

	.modal-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 1.5rem;
		border-bottom: 1px solid #e5e7eb;
	}

	.modal-body {
		padding: 1.5rem;
	}

	.modal-footer {
		display: flex;
		justify-content: flex-end;
		gap: 0.75rem;
		padding: 1.5rem;
		border-top: 1px solid #e5e7eb;
	}
</style>
