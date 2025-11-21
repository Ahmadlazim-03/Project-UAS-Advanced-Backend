<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import { authStore } from '$lib/stores/auth';

	let achievements = $state<any[]>([]);
	let loading = $state(true);
	let showModal = $state(false);
	let isEditing = $state(false);
	let currentAchievement = $state<any>(null);

	let form = $state({
		title: '',
		description: '',
		category: '',
		points: 0,
		achievement_date: '',
		certificate_number: '',
		organizer: ''
	});

	let user = $derived($authStore.user);
	let isAdmin = $derived(user?.role.name === 'Admin');
	let isMahasiswa = $derived(user?.role.name === 'Mahasiswa');

	onMount(async () => {
		await loadAchievements();
	});

	async function loadAchievements() {
		loading = true;
		try {
			const result = await api.getAchievements();
			if (result.status === 'success') {
				// Backend returns array of {id, status, submittedAt, verifiedAt, data}
				achievements = (result.data || []).map((item: any) => ({
					id: item.id,
					status: item.status,
					submittedAt: item.submittedAt,
					verifiedAt: item.verifiedAt,
					rejectionNote: item.rejectionNote,
					// Flatten the MongoDB data
					title: item.data?.title || '',
					description: item.data?.description || '',
					category: item.data?.category || '',
					points: item.data?.points || 0,
					achievement_date: item.data?.achievementDate || item.data?.achievement_date || '',
					certificate_number: item.data?.certificateNo || item.data?.certificate_number || '',
					organizer: item.data?.organizer || ''
				}));
			}
		} catch (error) {
			console.error('Failed to load achievements:', error);
		} finally {
			loading = false;
		}
	}

	function openCreateModal() {
		isEditing = false;
		form = {
			title: '',
			description: '',
			category: '',
			points: 0,
			achievement_date: '',
			certificate_number: '',
			organizer: ''
		};
		showModal = true;
	}

	function openEditModal(achievement: any) {
		isEditing = true;
		currentAchievement = achievement;
		const achievementDate = achievement.achievement_date || achievement.achievementDate;
		form = {
			title: achievement.title || '',
			description: achievement.description || '',
			category: achievement.category || '',
			points: achievement.points || 0,
			achievement_date: achievementDate ? (typeof achievementDate === 'string' ? achievementDate.split('T')[0] : new Date(achievementDate).toISOString().split('T')[0]) : '',
			certificate_number: achievement.certificate_number || achievement.certificateNo || '',
			organizer: achievement.organizer || ''
		};
		showModal = true;
	}

	async function handleSubmit() {
		try {
			if (isEditing && currentAchievement) {
				await api.updateAchievement(currentAchievement.id, form);
			} else {
				await api.createAchievement(form);
			}
			showModal = false;
			await loadAchievements();
		} catch (error) {
			console.error('Failed to save achievement:', error);
			alert('Failed to save achievement');
		}
	}

	async function handleDelete(id: string) {
		if (!confirm('Are you sure you want to delete this achievement?')) return;

		try {
			await api.deleteAchievement(id);
			await loadAchievements();
		} catch (error) {
			console.error('Failed to delete achievement:', error);
			alert('Failed to delete achievement');
		}
	}

	async function handleSubmitForVerification(id: string) {
		try {
			await api.submitAchievement(id);
			await loadAchievements();
		} catch (error) {
			console.error('Failed to submit achievement:', error);
			alert('Failed to submit achievement');
		}
	}

	function getStatusBadge(status: string) {
		const badges: Record<string, { class: string; text: string }> = {
			draft: { class: 'bg-gray-100 text-gray-700', text: 'Draft' },
			submitted: { class: 'bg-yellow-100 text-yellow-700', text: 'Submitted' },
			verified: { class: 'bg-green-100 text-green-700', text: 'Verified' },
			rejected: { class: 'bg-red-100 text-red-700', text: 'Rejected' }
		};
		return badges[status] || badges.draft;
	}
</script>

<div class="space-y-6">
	<div class="flex justify-between items-center">
		<div>
			<h1 class="text-3xl font-bold text-gray-800">Achievements</h1>
			<p class="text-gray-600 mt-1">Manage your achievements and track your progress</p>
		</div>
		{#if isMahasiswa}
			<button onclick={() => openCreateModal()} class="btn btn-primary">
				<svg class="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 20 20">
					<path
						fill-rule="evenodd"
						d="M10 3a1 1 0 011 1v5h5a1 1 0 110 2h-5v5a1 1 0 11-2 0v-5H4a1 1 0 110-2h5V4a1 1 0 011-1z"
						clip-rule="evenodd"
					/>
				</svg>
				Add Achievement
			</button>
		{/if}
	</div>

	{#if loading}
		<div class="flex justify-center items-center py-12">
			<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
		</div>
	{:else if achievements.length === 0}
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
					d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
				/>
			</svg>
			<h3 class="text-lg font-medium text-gray-900 mb-2">No achievements yet</h3>
			<p class="text-gray-600 mb-4">Get started by adding your first achievement</p>
			{#if isMahasiswa}
				<button onclick={() => openCreateModal()} class="btn btn-primary">
					Add Your First Achievement
				</button>
			{/if}
		</div>
	{:else}
		<div class="grid grid-cols-1 gap-6">
			{#each achievements as achievement}
				{@const badge = getStatusBadge(achievement.status)}
				<div class="card p-6 hover:shadow-lg transition-shadow">
					<div class="flex justify-between items-start mb-4">
						<div class="flex-1">
							<div class="flex items-center gap-3 mb-2">
								<h3 class="text-xl font-semibold text-gray-800">{achievement.title}</h3>
								<span class="px-3 py-1 rounded-full text-sm font-medium {badge.class}">
									{badge.text}
								</span>
							</div>
							<p class="text-gray-600 mb-3">{achievement.description}</p>
							<div class="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
								<div>
									<span class="text-gray-500">Category:</span>
									<span class="font-medium text-gray-800 ml-1">{achievement.category}</span>
								</div>
								<div>
									<span class="text-gray-500">Points:</span>
									<span class="font-medium text-gray-800 ml-1">{achievement.points}</span>
								</div>
								<div>
									<span class="text-gray-500">Date:</span>
									<span class="font-medium text-gray-800 ml-1">
										{new Date(achievement.achievement_date).toLocaleDateString()}
									</span>
								</div>
								<div>
									<span class="text-gray-500">Organizer:</span>
									<span class="font-medium text-gray-800 ml-1">{achievement.organizer}</span>
								</div>
							</div>
							{#if achievement.certificate_number}
								<div class="mt-2 text-sm">
									<span class="text-gray-500">Certificate:</span>
									<span class="font-medium text-gray-800 ml-1">
										{achievement.certificate_number}
									</span>
								</div>
							{/if}
							{#if achievement.rejectionNote || achievement.rejection_note}
								<div class="mt-3 p-3 bg-red-50 border border-red-200 rounded-lg">
									<p class="text-sm text-red-800">
										<strong>Rejection Reason:</strong>
										{achievement.rejectionNote || achievement.rejection_note}
									</p>
								</div>
							{/if}
						</div>
					</div>

					<div class="flex gap-2 flex-wrap">
						{#if isMahasiswa && achievement.status === 'draft'}
							<button
								onclick={() => handleSubmitForVerification(achievement.id)}
								class="btn btn-primary btn-sm"
							>
								Submit for Verification
							</button>
							<button onclick={() => openEditModal(achievement)} class="btn btn-secondary btn-sm">
								Edit
							</button>
							<button onclick={() => handleDelete(achievement.id)} class="btn btn-danger btn-sm">
								Delete
							</button>
						{/if}

						{#if isAdmin}
							<button onclick={() => openEditModal(achievement)} class="btn btn-secondary btn-sm">
								Edit
							</button>
							<button onclick={() => handleDelete(achievement.id)} class="btn btn-danger btn-sm">
								Delete
							</button>
						{/if}
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>

<!-- Modal -->
{#if showModal}
	<div class="modal-overlay" onclick={() => (showModal = false)}>
		<div class="modal-content" onclick={(e) => e.stopPropagation()}>
			<div class="modal-header">
				<h2 class="text-2xl font-bold text-gray-800">
					{isEditing ? 'Edit Achievement' : 'Add New Achievement'}
				</h2>
				<button onclick={() => (showModal = false)} class="text-gray-400 hover:text-gray-600">
					<svg class="w-6 h-6" fill="currentColor" viewBox="0 0 20 20">
						<path
							fill-rule="evenodd"
							d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
							clip-rule="evenodd"
						/>
					</svg>
				</button>
			</div>

			<form onsubmit={(e) => (e.preventDefault(), handleSubmit())}>
				<div class="modal-body">
					<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
						<div class="md:col-span-2">
							<label class="form-label">Title</label>
							<input
								type="text"
								bind:value={form.title}
								class="form-input"
								placeholder="Achievement title"
								required
							/>
						</div>

						<div class="md:col-span-2">
							<label class="form-label">Description</label>
							<textarea
								bind:value={form.description}
								class="form-input"
								rows="3"
								placeholder="Describe your achievement"
								required
							></textarea>
						</div>

						<div>
							<label class="form-label">Category</label>
							<select bind:value={form.category} class="form-input" required>
								<option value="">Select category</option>
								<option value="Competition">Competition</option>
								<option value="Research">Research</option>
								<option value="Publication">Publication</option>
								<option value="Community Service">Community Service</option>
								<option value="Organization">Organization</option>
								<option value="Certification">Certification</option>
								<option value="Other">Other</option>
							</select>
						</div>

						<div>
							<label class="form-label">Points</label>
							<input
								type="number"
								bind:value={form.points}
								class="form-input"
								placeholder="0"
								min="0"
								required
							/>
						</div>

						<div>
							<label class="form-label">Achievement Date</label>
							<input type="date" bind:value={form.achievement_date} class="form-input" required />
						</div>

						<div>
							<label class="form-label">Certificate Number</label>
							<input
								type="text"
								bind:value={form.certificate_number}
								class="form-input"
								placeholder="Certificate number (optional)"
							/>
						</div>

						<div class="md:col-span-2">
							<label class="form-label">Organizer</label>
							<input
								type="text"
								bind:value={form.organizer}
								class="form-input"
								placeholder="Organization or institution"
								required
							/>
						</div>
					</div>
				</div>

				<div class="modal-footer">
					<button type="button" onclick={() => (showModal = false)} class="btn btn-secondary">
						Cancel
					</button>
					<button type="submit" class="btn btn-primary">
						{isEditing ? 'Update' : 'Create'} Achievement
					</button>
				</div>
			</form>
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
