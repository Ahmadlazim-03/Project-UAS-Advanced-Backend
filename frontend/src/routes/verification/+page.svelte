<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import { authStore } from '$lib/stores/auth';

	let achievements = $state<any[]>([]);
	let loading = $state(true);
	let showDetailModal = $state(false);
	let showRejectModal = $state(false);
	let selectedAchievement = $state<any>(null);
	let rejectNote = $state('');

	let user = $derived($authStore.user);

	onMount(async () => {
		await loadPendingVerifications();
	});

	async function loadPendingVerifications() {
		loading = true;
		try {
			const result = await api.getPendingVerifications();
			if (result.status === 'success') {
				// Backend returns array of {id, status, submittedAt, data, student}
				achievements = (result.data || []).map((item: any) => ({
					id: item.id,
					status: item.status,
					submittedAt: item.submittedAt,
					student: item.student,
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
			console.error('Failed to load pending verifications:', error);
		} finally {
			loading = false;
		}
	}

	function openDetailModal(achievement: any) {
		selectedAchievement = achievement;
		showDetailModal = true;
	}

	function openRejectModal(achievement: any) {
		selectedAchievement = achievement;
		rejectNote = '';
		showRejectModal = true;
	}

	async function handleVerify(id: string) {
		if (!confirm('Are you sure you want to verify this achievement?')) return;

		try {
			await api.verifyAchievement(id);
			showDetailModal = false;
			await loadPendingVerifications();
		} catch (error) {
			console.error('Failed to verify achievement:', error);
			alert('Failed to verify achievement');
		}
	}

	async function handleReject() {
		if (!rejectNote.trim()) {
			alert('Please provide a reason for rejection');
			return;
		}

		try {
			if (selectedAchievement) {
				await api.rejectAchievement(selectedAchievement.id, rejectNote);
				showRejectModal = false;
				showDetailModal = false;
				await loadPendingVerifications();
			}
		} catch (error) {
			console.error('Failed to reject achievement:', error);
			alert('Failed to reject achievement');
		}
	}
</script>

<div class="space-y-6">
	<div>
		<h1 class="text-3xl font-bold text-gray-800">Verification</h1>
		<p class="text-gray-600 mt-1">Review and verify student achievements</p>
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
					d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
				/>
			</svg>
			<h3 class="text-lg font-medium text-gray-900 mb-2">No pending verifications</h3>
			<p class="text-gray-600">All achievements have been reviewed</p>
		</div>
	{:else}
		<div class="grid grid-cols-1 gap-6">
			{#each achievements as achievement}
				<div class="card p-6 hover:shadow-lg transition-shadow">
					<div class="flex justify-between items-start">
						<div class="flex-1">
							<div class="flex items-center gap-3 mb-2">
								<h3 class="text-xl font-semibold text-gray-800">{achievement.title}</h3>
								<span class="px-3 py-1 bg-yellow-100 text-yellow-700 rounded-full text-sm font-medium">
									Pending Review
								</span>
							</div>

							{#if achievement.student}
								<div class="mb-3 flex items-center gap-2 text-sm text-gray-600">
									<svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
										<path
											fill-rule="evenodd"
											d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z"
											clip-rule="evenodd"
										/>
									</svg>
									<span>
										{achievement.student.full_name || achievement.student.user?.full_name || 'N/A'}
										{#if achievement.student.nim}
											<span class="text-gray-500">({achievement.student.nim})</span>
										{/if}
									</span>
								</div>
							{/if}

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
						</div>
					</div>

					<div class="flex gap-2 mt-4">
						<button onclick={() => openDetailModal(achievement)} class="btn btn-secondary btn-sm">
							<svg class="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
								<path d="M10 12a2 2 0 100-4 2 2 0 000 4z" />
								<path
									fill-rule="evenodd"
									d="M.458 10C1.732 5.943 5.522 3 10 3s8.268 2.943 9.542 7c-1.274 4.057-5.064 7-9.542 7S1.732 14.057.458 10zM14 10a4 4 0 11-8 0 4 4 0 018 0z"
									clip-rule="evenodd"
								/>
							</svg>
							View Details
						</button>
						<button
							onclick={() => handleVerify(achievement.id)}
							class="btn btn-sm bg-green-600 hover:bg-green-700 text-white"
						>
							<svg class="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
								<path
									fill-rule="evenodd"
									d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
									clip-rule="evenodd"
								/>
							</svg>
							Verify
						</button>
						<button
							onclick={() => openRejectModal(achievement)}
							class="btn btn-sm bg-red-600 hover:bg-red-700 text-white"
						>
							<svg class="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
								<path
									fill-rule="evenodd"
									d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
									clip-rule="evenodd"
								/>
							</svg>
							Reject
						</button>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>

<!-- Detail Modal -->
{#if showDetailModal && selectedAchievement}
	<div class="modal-overlay" onclick={() => (showDetailModal = false)}>
		<div class="modal-content" onclick={(e) => e.stopPropagation()}>
			<div class="modal-header">
				<h2 class="text-2xl font-bold text-gray-800">Achievement Details</h2>
				<button onclick={() => (showDetailModal = false)} class="text-gray-400 hover:text-gray-600">
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
					{#if selectedAchievement.student}
						<div class="p-4 bg-blue-50 border border-blue-200 rounded-lg">
							<h3 class="font-semibold text-gray-800 mb-2">Student Information</h3>
							<div class="grid grid-cols-2 gap-3 text-sm">
								<div>
									<span class="text-gray-600">Name:</span>
									<span class="font-medium text-gray-800 ml-1">
										{selectedAchievement.student.full_name ||
											selectedAchievement.student.user?.full_name ||
											'N/A'}
									</span>
								</div>
								{#if selectedAchievement.student.nim}
									<div>
										<span class="text-gray-600">NIM:</span>
										<span class="font-medium text-gray-800 ml-1">
											{selectedAchievement.student.nim}
										</span>
									</div>
								{/if}
								{#if selectedAchievement.student.major}
									<div>
										<span class="text-gray-600">Major:</span>
										<span class="font-medium text-gray-800 ml-1">
											{selectedAchievement.student.major}
										</span>
									</div>
								{/if}
							</div>
						</div>
					{/if}

					<div>
						<h3 class="font-semibold text-gray-800 mb-2">Achievement Information</h3>
						<div class="space-y-3">
							<div>
								<label class="text-sm text-gray-600">Title</label>
								<p class="font-medium text-gray-800">{selectedAchievement.title}</p>
							</div>
							<div>
								<label class="text-sm text-gray-600">Description</label>
								<p class="text-gray-800">{selectedAchievement.description}</p>
							</div>
							<div class="grid grid-cols-2 gap-4">
								<div>
									<label class="text-sm text-gray-600">Category</label>
									<p class="font-medium text-gray-800">{selectedAchievement.category}</p>
								</div>
								<div>
									<label class="text-sm text-gray-600">Points</label>
									<p class="font-medium text-gray-800">{selectedAchievement.points}</p>
								</div>
							</div>
							<div class="grid grid-cols-2 gap-4">
								<div>
									<label class="text-sm text-gray-600">Date</label>
									<p class="font-medium text-gray-800">
										{new Date(selectedAchievement.achievement_date).toLocaleDateString()}
									</p>
								</div>
								<div>
									<label class="text-sm text-gray-600">Organizer</label>
									<p class="font-medium text-gray-800">{selectedAchievement.organizer}</p>
								</div>
							</div>
							{#if selectedAchievement.certificate_number}
								<div>
									<label class="text-sm text-gray-600">Certificate Number</label>
									<p class="font-medium text-gray-800">{selectedAchievement.certificate_number}</p>
								</div>
							{/if}
						</div>
					</div>
				</div>
			</div>

			<div class="modal-footer">
				<button onclick={() => (showDetailModal = false)} class="btn btn-secondary">Close</button>
				<button
					onclick={() => openRejectModal(selectedAchievement)}
					class="btn bg-red-600 hover:bg-red-700 text-white"
				>
					Reject
				</button>
				<button
					onclick={() => handleVerify(selectedAchievement.id)}
					class="btn bg-green-600 hover:bg-green-700 text-white"
				>
					Verify
				</button>
			</div>
		</div>
	</div>
{/if}

<!-- Reject Modal -->
{#if showRejectModal && selectedAchievement}
	<div class="modal-overlay" onclick={() => (showRejectModal = false)}>
		<div class="modal-content max-w-md" onclick={(e) => e.stopPropagation()}>
			<div class="modal-header">
				<h2 class="text-2xl font-bold text-gray-800">Reject Achievement</h2>
				<button onclick={() => (showRejectModal = false)} class="text-gray-400 hover:text-gray-600">
					<svg class="w-6 h-6" fill="currentColor" viewBox="0 0 20 20">
						<path
							fill-rule="evenodd"
							d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
							clip-rule="evenodd"
						/>
					</svg>
				</button>
			</div>

			<form onsubmit={(e) => (e.preventDefault(), handleReject())}>
				<div class="modal-body">
					<p class="text-gray-600 mb-4">
						Please provide a reason for rejecting this achievement:
					</p>
					<textarea
						bind:value={rejectNote}
						class="form-input"
						rows="4"
						placeholder="Explain why this achievement is being rejected..."
						required
					></textarea>
				</div>

				<div class="modal-footer">
					<button type="button" onclick={() => (showRejectModal = false)} class="btn btn-secondary">
						Cancel
					</button>
					<button type="submit" class="btn bg-red-600 hover:bg-red-700 text-white">
						Reject Achievement
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
