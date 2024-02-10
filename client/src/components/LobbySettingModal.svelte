<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import ApiService from '../services/ApiService';
	import Modal from './Modal.svelte';
	import CopyIcon from './icons/CopyIcon.svelte';
	import CheckIcon from './icons/CheckIcon.svelte';

	export let lobbyId: string;
	export let link: string;
	export let name: string;
	export let editable = false;

	let errorMessage: string;
	let loading = false;

	const apiService = new ApiService();

	const dispatch = createEventDispatcher<{
		closeModal: boolean;
	}>();

	async function update() {
		if (loading) {
			return;
		}

		loading = true;

		const response = await apiService.updateLobby(lobbyId, { name: name });

		loading = false;

		if (response.status === 204) {
			closeModal();
		} else {
			const data = await response.json();
			errorMessage = data.message;
		}
	}

	function closeModal() {
		dispatch('closeModal', true);
	}

	// Copy lobby link

	let copied = false;

	async function copyLink() {
		if (copied) return;
		try {
			navigator.clipboard.writeText(link).then(() => {
				copied = true;
				setTimeout(() => {
					copied = false;
				}, 4000);
			});
		} catch (err) {
			copied = false;
		}
	}
</script>

<Modal header="Lobby Settings" on:closeModal={closeModal}>
	<div slot="body" class="lobby-settings-body">
		<form action="">
			<div class="mb-3">
				<label for="name">Name</label>
				<input type="text" id="name" disabled={!editable} bind:value={name} required />
				{#if errorMessage}
					<span class="text-red-500 text-sm">{errorMessage}</span>
				{/if}
			</div>
			<div>
				<label for="link">Lobby Link</label>
				<div class="relative">
					<input type="text" id="link" bind:value={link} disabled />
					<div class="absolute top-2 right-1 w-10 bg-slate-800">
						<div class="w-full h-full flex justify-end items-center">
							{#if copied}
								<CheckIcon />
							{:else}
								<!-- svelte-ignore a11y-click-events-have-key-events -->
								<div on:click={() => copyLink()}>
									<CopyIcon />
								</div>
							{/if}
						</div>
					</div>
				</div>
			</div>
		</form>
	</div>
	<slot slot="footer">
		{#if editable}
			<button on:click={() => update()} class="btn {loading ? 'loading' : ''}" type="button">
				Update
			</button>
		{:else}
			<p class="text-red-400">You don't have permission to update lobby setting.</p>
		{/if}
	</slot>
</Modal>

<style lang="scss">
	.lobby-settings-body {
		@apply w-full max-h-96 py-6 px-4 overflow-y-auto;
		min-width: 24rem;
	}

	@media (max-width: 640px) {
		.lobby-settings-body {
			@apply min-w-full;
		}
	}
</style>
