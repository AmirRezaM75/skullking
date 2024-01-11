<script lang="ts">
	import AvatarModel from '../../../components/AvatarModel.svelte';
	import GameLinkDialog from '../../../components/GameLinkDialog.svelte';
	import LobbySidebar from '../../../components/LobbySidebar.svelte';
	import Pencil from '../../../components/icons/PencilIcon.svelte';
	let isAvatarModalOpen = false;

	var data = {
		authId: '123'
	};

	var game = {
		players: [
			{
				id: '123',
				avatar: '/images/avatars/1.jpg',
				username: 'amirrezam'
			},
			{
				id: '124',
				avatar: '/images/avatars/2.jpg',
				username: 'skullking'
			},
			{
				id: '125',
				avatar: '/images/avatars/3.jpg',
				username: 'darkman'
			}
		],
		creator: {
			id: '123',
			username: 'AmirReza'
		}
	};

	function openAvatarModal(playerId: string) {
		if (data.authId === playerId) {
			isAvatarModalOpen = true;
		}
	}

	function closeAvatarModal() {
		isAvatarModalOpen = false;
	}
</script>

<div class="min-w-full min-h-screen flex items-center justify-center bg-slate-700">
	<!-- <LobbySidebar /> -->
	{#if isAvatarModalOpen}
		<AvatarModel on:closeModal={closeAvatarModal} />
	{/if}
	<div class="flex-col">
		<div class="flex items-center justify-center gap-4 flex-wrap px-2 py-4 max-w-2xl">
			{#each game.players as player}
				<div
					on:click={() => openAvatarModal(player.id)}
					on:keypress={() => openAvatarModal(player.id)}
					class="relative bg-slate-900 p-6 rounded-lg text-center border border-slate-900
					{player.id === data.authId ? 'cursor-pointer hover:border-lime-primary' : ''}"
				>
					{#if data.authId === player.id}
						<div class="absolute top-3 left-3 w-4">
							<Pencil />
						</div>
					{/if}
					<div
						class="mb-3 rounded-full w-24 h-24 bg-top bg-no-repeat"
						style="background-image: url({player.avatar}); background-size: 110px"
					/>
					<span class="text-gray-300 font-bold text-lg uppercase">{player.username}</span>
				</div>
			{/each}
		</div>
		{#if game.creator.id !== data.authId}
			<p class="text-yellow-500 text-center">
				Wait for {game.creator.username} to start the game.
			</p>
		{:else if game.players.length != 1}
			<div class="text-center mt-6 mb-6 sm:mb-0">
				<button type="button" class="btn-secondary">
					<span>Start</span>
				</button>
			</div>
		{/if}
	</div>
	{#if game.players.length === 1}
		<GameLinkDialog />
	{/if}
</div>
