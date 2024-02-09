<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import ApiService from '../services/ApiService';
	import Modal from './Modal.svelte';
	import AuthService from '../services/AuthService';

	const apiService = new ApiService();

	export let currentAvatarId: number;

	let avatarId = currentAvatarId;

	function choose(id: number) {
		avatarId = id;
	}

	const dispatch = createEventDispatcher<{
		closeModal: boolean;
		avatarIdUpdated: { avatarId: number };
	}>();

	async function update() {
		const response = await apiService.updateAvatarId(avatarId);
		if (response.status === 204) {
			currentAvatarId = avatarId;
			const authService = new AuthService();
			authService.updateAvatarId(avatarId);
			dispatch('avatarIdUpdated', { avatarId: avatarId });
			closeModal();
		}
	}

	function closeModal() {
		dispatch('closeModal', true);
	}
</script>

<Modal header="Update your avatar" on:closeModal={closeModal}>
	<div
		slot="body"
		class="w-full max-h-96 grid gap-2 flex-wrap py-6 px-4 overflow-y-auto"
		style="grid-template-columns: repeat(auto-fit, minmax(100px, 1fr));"
	>
		{#each { length: 30 } as _, i}
			<!-- svelte-ignore a11y-click-events-have-key-events -->
			<img
				on:click={() => choose(i)}
				src="/images/avatars/{i + 1}.jpg"
				alt=""
				class="border-4 hover:opacity-100 cursor-pointer {avatarId == i
					? 'border-indigo-500'
					: 'border-slate-700 opacity-50'}"
			/>
		{/each}
	</div>
	<slot slot="footer">
		<div class="text-gray-400">
			Images by <a
				class="text-indigo-500"
				target="_blank"
				href="https://www.freepik.com/serie/27470333">Freepick</a
			>
		</div>
		<button
			on:click={() => update()}
			class="px-2 text-slate-800 bg-lime-primary py-2 rounded outline-none focus:outline-none"
			type="button"
		>
			Update
		</button>
	</slot>
</Modal>
