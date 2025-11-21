<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import { authStore } from '$lib/stores/auth';

	let students = $state<any[]>([]);
	let lecturers = $state<any[]>([]);
	let users = $state<any[]>([]);
	let loading = $state(true);
	let showAssignAdvisorModal = $state(false);
	let showStudentModal = $state(false);
	let selectedStudent = $state<any>(null);
	let selectedAdvisorId = $state('');
	let isEditMode = $state(false);

	// Form fields
	let formData = $state({
		user_id: '',
		student_id: '',
		program_study: '',
		academic_year: new Date().getFullYear().toString()
	});

	let user = $derived($authStore.user);

	onMount(async () => {
		await loadData();
	});

	async function loadData() {
		loading = true;
		try {
			const [studentsResult, lecturersResult, usersResult] = await Promise.all([
				api.getStudents(),
				api.getLecturers(),
				api.getUsers()
			]);

			if (studentsResult.status === 'success') {
				students = studentsResult.data || [];
			}

			if (lecturersResult.status === 'success') {
				lecturers = lecturersResult.data || [];
			}

			if (usersResult.status === 'success') {
				// Filter only users with Mahasiswa role
				users = (usersResult.data || []).filter((u: any) => u.role?.name === 'Mahasiswa');
			}
		} catch (error) {
			console.error('Failed to load data:', error);
		} finally {
			loading = false;
		}
	}

	function openCreateModal() {
		isEditMode = false;
		selectedStudent = null;
		formData = {
			user_id: '',
			student_id: '',
			program_study: '',
			academic_year: new Date().getFullYear().toString()
		};
		showStudentModal = true;
	}

	function openEditModal(student: any) {
		isEditMode = true;
		selectedStudent = student;
		formData = {
			user_id: student.user_id || '',
			student_id: student.student_id || student.nim || '',
			program_study: student.program_study || student.major || '',
			academic_year: student.academic_year || new Date().getFullYear().toString()
		};
		showStudentModal = true;
	}

	async function saveStudent() {
		try {
			if (isEditMode && selectedStudent) {
				await api.updateStudent(selectedStudent.id, formData);
			} else {
				await api.createStudent(formData);
			}
			showStudentModal = false;
			await loadData();
		} catch (error) {
			console.error('Failed to save student:', error);
			alert('Failed to save student');
		}
	}

	async function deleteStudent(student: any) {
		if (!confirm(`Are you sure you want to delete ${student.user?.full_name || 'this student'}?`)) {
			return;
		}

		try {
			await api.deleteStudent(student.id);
			await loadData();
		} catch (error) {
			console.error('Failed to delete student:', error);
			alert('Failed to delete student');
		}
	}

	function openAssignAdvisorModal(student: any) {
		selectedStudent = student;
		selectedAdvisorId = student.advisor_id || student.lecturer_id || '';
		showAssignAdvisorModal = true;
	}

	async function assignAdvisor() {
		if (!selectedAdvisorId) {
			alert('Please select an advisor');
			return;
		}

		try {
			await api.updateStudentAdvisor(selectedStudent.id, selectedAdvisorId);
			showAssignAdvisorModal = false;
			await loadData();
		} catch (error) {
			console.error('Failed to assign advisor:', error);
			alert('Failed to assign advisor');
		}
	}

	async function viewStudentAchievements(studentId: string) {
		try {
			const result = await api.getStudentAchievements(studentId);
			if (result.status === 'success') {
				const count = result.data?.length || 0;
				alert(`This student has ${count} achievement(s)`);
			}
		} catch (error) {
			console.error('Failed to load achievements:', error);
		}
	}
</script>

<div class="space-y-6">
	<div class="flex justify-between items-center">
		<div>
			<h1 class="text-3xl font-bold text-gray-800">Students Management</h1>
			<p class="text-gray-600 mt-1">Manage students and their advisors</p>
		</div>
		<button onclick={openCreateModal} class="btn bg-primary-600 hover:bg-primary-700 text-white px-6 py-3 rounded-lg font-medium transition flex items-center gap-2">
			<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
				<path fill-rule="evenodd" d="M10 3a1 1 0 011 1v5h5a1 1 0 110 2h-5v5a1 1 0 11-2 0v-5H4a1 1 0 110-2h5V4a1 1 0 011-1z" clip-rule="evenodd" />
			</svg>
			Add Student
		</button>
	</div>

	{#if loading}
		<div class="flex justify-center items-center py-12">
			<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
		</div>
	{:else if students.length === 0}
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
			<h3 class="text-lg font-medium text-gray-900 mb-2">No students found</h3>
			<p class="text-gray-600">Students will appear here once they are registered</p>
		</div>
	{:else}
		<!-- Summary Stats -->
		<div class="grid grid-cols-1 md:grid-cols-3 gap-6">
			<div class="card p-6 bg-gradient-to-br from-blue-500 to-blue-600 text-white">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm opacity-90">Total Students</p>
						<p class="text-3xl font-bold mt-1">{students.length}</p>
					</div>
					<svg class="w-12 h-12 opacity-80" fill="currentColor" viewBox="0 0 20 20">
						<path d="M9 6a3 3 0 11-6 0 3 3 0 016 0zM17 6a3 3 0 11-6 0 3 3 0 016 0zM12.93 17c.046-.327.07-.66.07-1a6.97 6.97 0 00-1.5-4.33A5 5 0 0119 16v1h-6.07zM6 11a5 5 0 015 5v1H1v-1a5 5 0 015-5z" />
					</svg>
				</div>
			</div>

			<div class="card p-6 bg-gradient-to-br from-purple-500 to-purple-600 text-white">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm opacity-90">With Advisor</p>
						<p class="text-3xl font-bold mt-1">
							{students.filter(s => s.advisor_id || s.advisor).length}
						</p>
					</div>
					<svg class="w-12 h-12 opacity-80" fill="currentColor" viewBox="0 0 20 20">
						<path fill-rule="evenodd" d="M6.267 3.455a3.066 3.066 0 001.745-.723 3.066 3.066 0 013.976 0 3.066 3.066 0 001.745.723 3.066 3.066 0 012.812 2.812c.051.643.304 1.254.723 1.745a3.066 3.066 0 010 3.976 3.066 3.066 0 00-.723 1.745 3.066 3.066 0 01-2.812 2.812 3.066 3.066 0 00-1.745.723 3.066 3.066 0 01-3.976 0 3.066 3.066 0 00-1.745-.723 3.066 3.066 0 01-2.812-2.812 3.066 3.066 0 00-.723-1.745 3.066 3.066 0 010-3.976 3.066 3.066 0 00.723-1.745 3.066 3.066 0 012.812-2.812zm7.44 5.252a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
					</svg>
				</div>
			</div>

			<div class="card p-6 bg-gradient-to-br from-green-500 to-green-600 text-white">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm opacity-90">Without Advisor</p>
						<p class="text-3xl font-bold mt-1">
							{students.filter(s => !s.advisor_id && !s.advisor).length}
						</p>
					</div>
					<svg class="w-12 h-12 opacity-80" fill="currentColor" viewBox="0 0 20 20">
						<path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
					</svg>
				</div>
			</div>
		</div>

		<!-- Students Table -->
		<div class="card overflow-hidden">
			<div class="overflow-x-auto">
				<table class="min-w-full divide-y divide-gray-200">
					<thead class="bg-gray-50">
						<tr>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
								Student
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
								Student ID / NIM
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
								Program
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
								Advisor
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
								Actions
							</th>
						</tr>
					</thead>
					<tbody class="bg-white divide-y divide-gray-200">
						{#each students as student}
							<tr class="hover:bg-gray-50 transition-colors">
								<td class="px-6 py-4 whitespace-nowrap">
									<div class="flex items-center">
										<div class="flex-shrink-0 h-10 w-10">
											<div
												class="h-10 w-10 rounded-full bg-gradient-to-br from-blue-400 to-blue-600 flex items-center justify-center text-white font-semibold"
											>
												{student.user?.full_name?.charAt(0).toUpperCase() || student.full_name?.charAt(0).toUpperCase() || 'S'}
											</div>
										</div>
										<div class="ml-4">
											<div class="text-sm font-medium text-gray-900">
												{student.user?.full_name || student.full_name || 'N/A'}
											</div>
											<div class="text-sm text-gray-500">
												{student.user?.email || 'N/A'}
											</div>
										</div>
									</div>
								</td>
								<td class="px-6 py-4 whitespace-nowrap">
									<div class="text-sm text-gray-900">{student.student_id || student.nim || 'N/A'}</div>
									<div class="text-sm text-gray-500">{student.academic_year || 'N/A'}</div>
								</td>
								<td class="px-6 py-4 whitespace-nowrap">
									<div class="text-sm text-gray-900">{student.program_study || student.major || 'N/A'}</div>
								</td>
								<td class="px-6 py-4 whitespace-nowrap">
									{#if student.advisor}
										<div class="text-sm text-gray-900">
											{student.advisor.user?.full_name || student.advisor.full_name || 'N/A'}
										</div>
										<div class="text-sm text-gray-500">
											{student.advisor.lecturer_id || student.advisor.nip || ''}
										</div>
									{:else}
										<span class="px-2 py-1 text-xs font-semibold rounded-full bg-yellow-100 text-yellow-800">
											No Advisor
										</span>
									{/if}
								</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm font-medium space-x-2">
								<button
									onclick={() => openEditModal(student)}
									class="text-blue-600 hover:text-blue-900"
								>
									Edit
								</button>
								<button
									onclick={() => openAssignAdvisorModal(student)}
									class="text-primary-600 hover:text-primary-900"
								>
									{student.advisor ? 'Change' : 'Assign'} Advisor
								</button>
								<button
									onclick={() => viewStudentAchievements(student.id)}
									class="text-green-600 hover:text-green-900"
								>
									View Achievements
								</button>
								<button
									onclick={() => deleteStudent(student)}
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
	{/if}
</div>

<!-- Assign Advisor Modal -->
{#if showAssignAdvisorModal && selectedStudent}
	<div class="modal-overlay" onclick={() => (showAssignAdvisorModal = false)}>
		<div class="modal-content max-w-md" onclick={(e) => e.stopPropagation()}>
			<div class="modal-header">
				<h2 class="text-2xl font-bold text-gray-800">Assign Advisor</h2>
				<button onclick={() => (showAssignAdvisorModal = false)} class="text-gray-400 hover:text-gray-600">
					<svg class="w-6 h-6" fill="currentColor" viewBox="0 0 20 20">
						<path
							fill-rule="evenodd"
							d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
							clip-rule="evenodd"
						/>
					</svg>
				</button>
			</div>

			<form onsubmit={(e) => (e.preventDefault(), assignAdvisor())}>
				<div class="modal-body">
					<div class="mb-4">
						<p class="text-sm text-gray-600 mb-1">Student:</p>
						<p class="font-medium text-gray-900">
							{selectedStudent.user?.full_name || selectedStudent.full_name} ({selectedStudent.student_id || selectedStudent.nim})
						</p>
					</div>

					<div>
						<label class="form-label">Select Advisor (Dosen Wali)</label>
						<select bind:value={selectedAdvisorId} class="form-input" required>
							<option value="">-- Select Advisor --</option>
							{#each lecturers as lecturer}
								<option value={lecturer.id}>
									{lecturer.user?.full_name || lecturer.full_name} - {lecturer.lecturer_id || lecturer.nip}
								</option>
							{/each}
						</select>
					</div>
				</div>

				<div class="modal-footer">
					<button type="button" onclick={() => (showAssignAdvisorModal = false)} class="btn btn-secondary">
						Cancel
					</button>
					<button type="submit" class="btn btn-primary">
						Assign Advisor
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}

<!-- Create/Edit Student Modal -->
{#if showStudentModal}
	<div class="modal-overlay" onclick={() => (showStudentModal = false)}>
		<div class="modal-content max-w-md" onclick={(e) => e.stopPropagation()}>
			<div class="modal-header">
				<h2 class="text-2xl font-bold text-gray-800">
					{isEditMode ? 'Edit Student' : 'Create Student'}
				</h2>
				<button onclick={() => (showStudentModal = false)} class="text-gray-400 hover:text-gray-600">
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
						<label class="form-label">Student User Account</label>
						<select bind:value={formData.user_id} class="input w-full" disabled={isEditMode}>
							<option value="">Select User</option>
							{#each users as usr}
								<option value={usr.id}>{usr.full_name} ({usr.email})</option>
							{/each}
						</select>
						<p class="text-sm text-gray-500 mt-1">Select the user account for this student</p>
					</div>

					<div>
						<label class="form-label">Student ID (NIM)</label>
						<input
							type="text"
							bind:value={formData.student_id}
							class="input w-full"
							placeholder="e.g., 2024001"
						/>
					</div>

					<div>
						<label class="form-label">Program Study / Major</label>
						<input
							type="text"
							bind:value={formData.program_study}
							class="input w-full"
							placeholder="e.g., Computer Science"
						/>
					</div>

					<div>
						<label class="form-label">Academic Year</label>
						<input
							type="text"
							bind:value={formData.academic_year}
							class="input w-full"
							placeholder="e.g., 2024"
						/>
					</div>
				</div>
			</div>

			<div class="modal-footer">
				<button onclick={() => (showStudentModal = false)} class="btn btn-secondary">
					Cancel
				</button>
				<button onclick={saveStudent} class="btn btn-primary">
					{isEditMode ? 'Update' : 'Create'} Student
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
