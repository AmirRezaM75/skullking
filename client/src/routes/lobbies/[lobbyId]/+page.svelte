<script lang="ts">
	import AvatarModel from '../../../components/AvatarModel.svelte';
	import GameLinkDialog from '../../../components/GameLinkDialog.svelte';
	import LobbySidebar from '../../../components/LobbySidebar.svelte';
	import PencilIcon from '../../../components/icons/PencilIcon.svelte';
	import ApiService from '../../../services/ApiService';
	import LobbyService from '../../../services/LobbyService';
	let isAvatarModalOpen = false;

	export let data;

	let currentAvatarId = data.auth.avatarId

	const apiService = new ApiService();
	const sse = apiService.joinLobby(data.lobbyId, data.ticketId);

	var isOpen = false;
	sse.onopen = (...args) => {
		if (isOpen) {
			sse.close();
		} else {
			console.log('Connection to server opened.');
			isOpen = true;
		}
	};
	sse.onerror = function (e) {
		console.log('error', e);
	};

	let lobbyService = new LobbyService();

	sse.onmessage = (message) => {
		const m = JSON.parse(message.data);
		lobbyService = lobbyService.handle(m.type, JSON.parse(m.content));
	};

	function openAvatarModal(playerId: string) {
		if (data.auth.id === playerId) {
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
		<AvatarModel
		on:closeModal={closeAvatarModal}
		on:avatarIdUpdated={(e) => currentAvatarId = e.detail.avatarId}
		currentAvatarId={currentAvatarId} />
	{/if}
	{#if lobbyService.lobby != null}
		<div class="flex-col">
			<div class="flex items-center justify-center gap-4 flex-wrap px-2 py-4 max-w-2xl">
				{#each lobbyService.lobby.players as player}
					<div
						on:click={() => openAvatarModal(player.id)}
						on:keypress={() => openAvatarModal(player.id)}
						class="relative bg-slate-900 p-6 rounded-lg text-center border border-slate-900
					{player.id === data.auth.id ? 'cursor-pointer hover:border-lime-primary' : ''}"
					>
						{#if data.auth.id === player.id}
							<div class="absolute top-3 left-3 w-4">
								<PencilIcon />
							</div>
						{/if}
						<div
							class="mb-3 rounded-full w-24 h-24 bg-top bg-no-repeat"
							style="background-image: url({`/images/avatars/${player.avatarId+1}.jpg`}); background-size: 110px"
						/>
						<span class="text-gray-300 font-bold text-lg uppercase">{player.username}</span>
					</div>
				{/each}
			</div>
			{#if lobbyService.lobby.getCreator()?.id !== data.auth.id}
				<p class="text-yellow-500 text-center">
					Wait for {lobbyService.lobby.getCreator()?.username} to start the game.
				</p>
			{:else if lobbyService.lobby.players.length != 1}
				<div class="text-center mt-6 mb-6 sm:mb-0">
					<button type="button" class="btn-secondary">
						<span>Start</span>
					</button>
				</div>
			{/if}
		</div>
		{#if lobbyService.lobby.players.length === 1}
			<GameLinkDialog />
		{/if}
	{/if}
</div>
