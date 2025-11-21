<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';

	let statistics = $state<any>(null);
	let achievements = $state<any[]>([]);
	let loading = $state(true);

	onMount(async () => {
		await loadData();
	});

	async function loadData() {
		loading = true;
		try {
			const [statsResult, achievementsResult] = await Promise.all([
				api.getStatistics(),
				api.getAchievements()
			]);

			if (statsResult.status === 'success') {
				statistics = statsResult.data;
			}

			if (achievementsResult.status === 'success') {
				achievements = achievementsResult.data || [];
			}
		} catch (error) {
			console.error('Failed to load data:', error);
		} finally {
			loading = false;
		}
	}

	let categoryStats = $derived(achievements.reduce((acc: Record<string, number>, achievement: any) => {
		const category = achievement.category || 'Other';
		acc[category] = (acc[category] || 0) + 1;
		return acc;
	}, {}));

	let verifiedByCategory = $derived(achievements
		.filter((a) => a.status === 'verified')
		.reduce((acc: Record<string, number>, achievement: any) => {
			const category = achievement.category || 'Other';
			acc[category] = (acc[category] || 0) + 1;
			return acc;
		}, {}));

	let totalPoints = $derived(achievements
		.filter((a) => a.status === 'verified')
		.reduce((sum, a) => sum + (a.points || 0), 0));

	let monthlyStats = $derived(achievements.reduce((acc: Record<string, number>, achievement: any) => {
		if (achievement.achievement_date) {
			const month = new Date(achievement.achievement_date).toLocaleDateString('en-US', {
				year: 'numeric',
				month: 'short'
			});
			acc[month] = (acc[month] || 0) + 1;
		}
		return acc;
	}, {}));
</script>

<div class="space-y-6">
	<div>
		<h1 class="text-3xl font-bold text-gray-800">Statistics & Reports</h1>
		<p class="text-gray-600 mt-1">Comprehensive overview of achievements and performance</p>
	</div>

	{#if loading}
		<div class="flex justify-center items-center py-12">
			<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
		</div>
	{:else if statistics}
		<!-- Summary Cards -->
		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
			<div class="card p-6 bg-gradient-to-br from-blue-500 to-blue-600 text-white">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm opacity-90">Total Achievements</p>
						<p class="text-3xl font-bold mt-1">{statistics.total_achievements || 0}</p>
					</div>
					<svg class="w-12 h-12 opacity-80" fill="currentColor" viewBox="0 0 20 20">
						<path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
					</svg>
				</div>
			</div>

			<div class="card p-6 bg-gradient-to-br from-green-500 to-green-600 text-white">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm opacity-90">Verified</p>
						<p class="text-3xl font-bold mt-1">{statistics.verified_achievements || 0}</p>
					</div>
					<svg class="w-12 h-12 opacity-80" fill="currentColor" viewBox="0 0 20 20">
						<path
							fill-rule="evenodd"
							d="M6.267 3.455a3.066 3.066 0 001.745-.723 3.066 3.066 0 013.976 0 3.066 3.066 0 001.745.723 3.066 3.066 0 012.812 2.812c.051.643.304 1.254.723 1.745a3.066 3.066 0 010 3.976 3.066 3.066 0 00-.723 1.745 3.066 3.066 0 01-2.812 2.812 3.066 3.066 0 00-1.745.723 3.066 3.066 0 01-3.976 0 3.066 3.066 0 00-1.745-.723 3.066 3.066 0 01-2.812-2.812 3.066 3.066 0 00-.723-1.745 3.066 3.066 0 010-3.976 3.066 3.066 0 00.723-1.745 3.066 3.066 0 012.812-2.812zm7.44 5.252a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
							clip-rule="evenodd"
						/>
					</svg>
				</div>
			</div>

			<div class="card p-6 bg-gradient-to-br from-yellow-500 to-yellow-600 text-white">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm opacity-90">Pending</p>
						<p class="text-3xl font-bold mt-1">{statistics.pending_verifications || 0}</p>
					</div>
					<svg class="w-12 h-12 opacity-80" fill="currentColor" viewBox="0 0 20 20">
						<path
							fill-rule="evenodd"
							d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z"
							clip-rule="evenodd"
						/>
					</svg>
				</div>
			</div>

			<div class="card p-6 bg-gradient-to-br from-purple-500 to-purple-600 text-white">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm opacity-90">Total Points</p>
						<p class="text-3xl font-bold mt-1">{totalPoints}</p>
					</div>
					<svg class="w-12 h-12 opacity-80" fill="currentColor" viewBox="0 0 20 20">
						<path d="M2 11a1 1 0 011-1h2a1 1 0 011 1v5a1 1 0 01-1 1H3a1 1 0 01-1-1v-5zM8 7a1 1 0 011-1h2a1 1 0 011 1v9a1 1 0 01-1 1H9a1 1 0 01-1-1V7zM14 4a1 1 0 011-1h2a1 1 0 011 1v12a1 1 0 01-1 1h-2a1 1 0 01-1-1V4z" />
					</svg>
				</div>
			</div>
		</div>

		<!-- Status Breakdown -->
		<div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
			<div class="card p-6">
				<h2 class="text-xl font-semibold text-gray-800 mb-4">Status Distribution</h2>
				<div class="space-y-4">
					<div>
						<div class="flex justify-between items-center mb-2">
							<span class="text-gray-700">Draft</span>
							<span class="font-semibold text-gray-900">{statistics.draft_achievements || 0}</span>
						</div>
						<div class="w-full bg-gray-200 rounded-full h-2">
							<div
								class="bg-gray-600 h-2 rounded-full"
								style="width: {(statistics.draft_achievements / statistics.total_achievements) * 100 || 0}%"
							></div>
						</div>
					</div>

					<div>
						<div class="flex justify-between items-center mb-2">
							<span class="text-gray-700">Submitted</span>
							<span class="font-semibold text-gray-900">{statistics.pending_verifications || 0}</span>
						</div>
						<div class="w-full bg-gray-200 rounded-full h-2">
							<div
								class="bg-yellow-500 h-2 rounded-full"
								style="width: {(statistics.pending_verifications / statistics.total_achievements) * 100 || 0}%"
							></div>
						</div>
					</div>

					<div>
						<div class="flex justify-between items-center mb-2">
							<span class="text-gray-700">Verified</span>
							<span class="font-semibold text-gray-900">{statistics.verified_achievements || 0}</span>
						</div>
						<div class="w-full bg-gray-200 rounded-full h-2">
							<div
								class="bg-green-500 h-2 rounded-full"
								style="width: {(statistics.verified_achievements / statistics.total_achievements) * 100 || 0}%"
							></div>
						</div>
					</div>

					<div>
						<div class="flex justify-between items-center mb-2">
							<span class="text-gray-700">Rejected</span>
							<span class="font-semibold text-gray-900">{statistics.rejected_achievements || 0}</span>
						</div>
						<div class="w-full bg-gray-200 rounded-full h-2">
							<div
								class="bg-red-500 h-2 rounded-full"
								style="width: {(statistics.rejected_achievements / statistics.total_achievements) * 100 || 0}%"
							></div>
						</div>
					</div>
				</div>
			</div>

			<div class="card p-6">
				<h2 class="text-xl font-semibold text-gray-800 mb-4">Category Distribution</h2>
				<div class="space-y-3">
					{#each Object.entries(categoryStats).sort((a, b) => b[1] - a[1]) as [category, count]}
						<div>
							<div class="flex justify-between items-center mb-1">
								<span class="text-sm text-gray-700">{category}</span>
								<span class="text-sm font-semibold text-gray-900">{count}</span>
							</div>
							<div class="w-full bg-gray-200 rounded-full h-2">
								<div
									class="bg-primary-600 h-2 rounded-full"
									style="width: {(Number(count) / achievements.length) * 100}%"
								></div>
							</div>
						</div>
					{/each}
				</div>
			</div>
		</div>

		<!-- User Statistics -->
		<div class="card p-6">
			<h2 class="text-xl font-semibold text-gray-800 mb-4">User Statistics</h2>
			<div class="grid grid-cols-1 md:grid-cols-3 gap-6">
				<div class="text-center p-4 bg-blue-50 rounded-lg">
					<p class="text-sm text-gray-600 mb-1">Mahasiswa</p>
					<p class="text-3xl font-bold text-blue-600">{statistics.students_count || 0}</p>
				</div>
				<div class="text-center p-4 bg-purple-50 rounded-lg">
					<p class="text-sm text-gray-600 mb-1">Dosen Wali</p>
					<p class="text-3xl font-bold text-purple-600">{statistics.lecturers_count || 0}</p>
				</div>
				<div class="text-center p-4 bg-indigo-50 rounded-lg">
					<p class="text-sm text-gray-600 mb-1">Admin</p>
					<p class="text-3xl font-bold text-indigo-600">{statistics.admins_count || 0}</p>
				</div>
			</div>
		</div>

		<!-- Monthly Trends -->
		{#if Object.keys(monthlyStats).length > 0}
			<div class="card p-6">
				<h2 class="text-xl font-semibold text-gray-800 mb-4">Monthly Achievement Trends</h2>
				<div class="space-y-3">
					{#each Object.entries(monthlyStats).sort((a, b) => new Date(a[0]).getTime() - new Date(b[0]).getTime()).slice(-6) as [month, count]}
						<div>
							<div class="flex justify-between items-center mb-1">
								<span class="text-sm text-gray-700">{month}</span>
								<span class="text-sm font-semibold text-gray-900">{count} achievements</span>
							</div>
							<div class="w-full bg-gray-200 rounded-full h-2">
								<div
									class="bg-primary-600 h-2 rounded-full"
									style="width: {(Number(count) / Math.max(...Object.values(monthlyStats).map(Number))) * 100}%"
								></div>
							</div>
						</div>
					{/each}
				</div>
			</div>
		{/if}
	{/if}
</div>
