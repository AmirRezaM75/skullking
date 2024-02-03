<script lang="ts">
	import { onMount } from 'svelte';
	import AvatarModel from '../../../components/AvatarModel.svelte';
	import LobbyLinkDialog from '../../../components/LobbyLinkDialog.svelte';
	import PencilIcon from '../../../components/icons/PencilIcon.svelte';
	import ApiService from '../../../services/ApiService';
	import LobbyService from '../../../services/LobbyService';
	import { goto } from '$app/navigation';
	import ConnectionErrorDialog from '../../../components/ConnectionErrorDialog.svelte';

	export let data;

	// Avatar
	let currentAvatarId = data.auth.avatarId;

	let isAvatarModalOpen = false;

	function openAvatarModal(playerId: string) {
		if (data.auth.id === playerId) {
			isAvatarModalOpen = true;
		}
	}

	function closeAvatarModal() {
		isAvatarModalOpen = false;
	}

	// SSE
	var isOpen = false;

	let disconnected = false;

	const apiService = new ApiService();

	let lobbyService = new LobbyService(data.auth.id);

	async function start() {
		const audio = new Audio('/sounds/start.mp3');
		audio.play();
		const response = await apiService.createGame(data.lobbyId);
		if (response.status === 201) {
			response.json().then((data) => {
				goto(`/games/${data.id}`);
			});
		}
	}

	onMount(async () => {
		let ticketId = '';

		const response = await apiService.createTicket();

		if (response.status === 201) {
			const data = await response.json();
			ticketId = data.id;
		}

		const sse = apiService.joinLobby(data.lobbyId, ticketId);

		sse.onopen = () => {
			if (isOpen) {
				sse.close();
			} else {
				console.log('Connection to server opened.');
				isOpen = true;
			}
		};

		sse.onerror = function () {
			// In case of timeout and opening duplicate tabs
			sse.close();
			disconnected = true;
		};

		sse.onmessage = (message) => {
			const m = JSON.parse(message.data);
			lobbyService = lobbyService.handle(m.type, JSON.parse(m.content));
		};
	});
</script>

<svelte:head>
	<title>Kenopsia - Lobby</title>
	{#each [...Array(30).keys()] as number}
		<link rel="preload" href="/images/avatars/{number + 1}.jpg" as="image" />
	{/each}
</svelte:head>

<div class="min-w-full min-h-screen flex items-center justify-center bg-slate-700">
	<!-- <LobbySidebar /> -->

	{#if disconnected}
		<ConnectionErrorDialog />
	{/if}

	{#if isAvatarModalOpen}
		<AvatarModel
			on:closeModal={closeAvatarModal}
			on:avatarIdUpdated={(e) => (currentAvatarId = e.detail.avatarId)}
			{currentAvatarId}
		/>
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
							style="background-image: url({`/images/avatars/${
								player.avatarId + 1
							}.jpg`}); background-size: 110px"
						/>
						<span class="text-gray-300 font-bold text-lg uppercase">{player.username}</span>
					</div>
				{/each}
			</div>
			{#if lobbyService.lobby.managerId !== data.auth.id}
				<p class="text-yellow-500 text-center">
					Wait for {lobbyService.lobby.getManager()?.username} to start the game.
				</p>
			{:else if lobbyService.lobby.players.length != 1}
				<div class="text-center mt-6 mb-6 sm:mb-0">
					<button type="button" class="btn-secondary" on:click={start}>
						<span>Start</span>
					</button>
				</div>
			{/if}
		</div>
		{#if lobbyService.lobby.players.length === 1}
			<LobbyLinkDialog />
		{/if}
	{/if}
</div>
