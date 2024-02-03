<script lang="ts">
	import { goto } from '$app/navigation';
	import AnimatedBackground from '../../components/AnimatedBackground.svelte';
	import ConnectionErrorDialog from '../../components/ConnectionErrorDialog.svelte';
	import Navigation from '../../components/Navigation.svelte';
	import ApiService from '../../services/ApiService';
	import LobbiesService from '../../services/LobbiesService';

	export let data;

	const apiService = new ApiService();

	let sse = apiService.getLobbies(data.ticketId);

	let lobbiesService = new LobbiesService();

	let disconnected = false;

	var isOpen = false;

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
		lobbiesService = lobbiesService.handle(m.type, JSON.parse(m.content));
	};

	async function create() {
		const response = await apiService.createLobby();
		if (response.status === 201) {
			response.json().then((data) => {
				goto(`/lobbies/${data.id}`);
			});
		}
	}
</script>

<svelte:head>
	<title>Lobbies - Kenopsia</title>
</svelte:head>

<section class="min-h-screen flex justify-center bg-slate-900">
	{#if disconnected}
		<ConnectionErrorDialog />
	{/if}
	<Navigation />
	<AnimatedBackground />
	<div class="w-full max-w-2xl px-4 2xl:px-0 mt-24 pb-4" style="z-index: 2;">
		<div class="bg-slate-800 shadow-lg rounded-lg">
			<div class="px-4 py-3">
				<div class="flex flex-wrap items-center">
					<div class="flex-1">
						<h3 class="font-semibold text-base text-gray-100">Lobbies</h3>
					</div>
					<div class="flex-1 text-right">
						<button
							on:click={create}
							class="bg-indigo-500 text-white active:bg-indigo-600 font-bold uppercase px-3 py-1 rounded outline-none focus:outline-none ease-linear transition-all duration-150"
							type="button"
						>
							Create
						</button>
					</div>
				</div>
			</div>

			<div class="block w-full overflow-x-auto">
				<table class="items-center bg-slate-800 w-full border-collapse rounded-b-lg">
					<thead>
						<tr
							class="align-middle text-xs uppercase whitespace-nowrap bg-slate-900 border border-solid border-gray-600 border-l-0 border-r-0 text-gray-400 font-semibold"
						>
							<th class="px-4 py-3 text-left">Lobby Name</th>
							<th class="px-4 py-3 text-center">Players</th>
							<th class="px-4 py-3 text-center">Created At</th>
							<th class="px-4 py-3" />
						</tr>
					</thead>

					<tbody>
						{#if lobbiesService.lobbies.length === 0}
							<tr class="text-center">
								<td colspan="4" class="py-4 text-red-300"
									>There are currently no public lobbies available.</td
								>
							</tr>
						{/if}

						{#each lobbiesService.lobbies as lobby}
							<tr class="text-gray-200">
								<td class="align-middle whitespace-nowrap p-4 text-left">
									{lobby.name}
								</td>
								<td class="align-middle whitespace-nowrap p-4 text-center">
									{lobby.players.length}
								</td>
								<td class="align-center whitespace-nowrap p-4 text-center">
									{lobby.getCreatedAt()}
								</td>
								<td class="align-middle whitespace-nowrap p-4 text-right">
									<a
										href="/lobbies/{lobby.id}"
										class="border border-lime-primary text-lime-primary hover:text-slate-800 hover:bg-lime-primary active:bg-lime-primary font-bold uppercase px-3 py-1 rounded outline-none focus:outline-none ease-linear transition-all duration-150"
										type="button"
									>
										Join
									</a>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		</div>
	</div>
</section>
