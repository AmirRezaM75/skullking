<script lang="ts">
	import GameService from '../../../services/GameService.js';
	import { GameState, GameCommand } from '../../../constants.js';
	import User from '../../../components/User.svelte';

	export let data;

	let game = new GameService();

	const ws = new WebSocket(
		'ws://localhost:3000/games/join?gameId=' + data.gameId + '&token=' + data.token
	);

	ws.onopen = function (e) {
		console.log('OPEN');
	};

	ws.onmessage = async function (e) {
		let message = JSON.parse(e.data);
		console.log(message);

		game = await game.handle(message['command'], message['content'], message['senderId']);
	};

	function start() {
		ws.send(
			JSON.stringify({
				command: GameCommand.Start
			})
		);
	}
</script>

<div
	class="min-w-full min-h-screen bg-gray-900 flex items-center justify-center"
	style="background-color: #1B1B1B;"
>
	{#if game.state == GameState.Pending}
		<div class="flex-col">
			<div class="flex items-center justify-center gap-4 flex-wrap px-2 py-4 max-w-2xl">
				{#each game.players as player}
					<div class="">
						<div class="bg-gray-700 p-6 rounded-lg text-center">
							<div class="mb-3">
								<img
									src="/images/avatars/1.jpg"
									width="100"
									height="100"
									alt=""
									class="rounded-full"
								/>
							</div>
							<span class="text-gray-100 font-bold text-lg uppercase">{player.username}</span>
						</div>
					</div>
				{/each}
			</div>
			<div class="text-center mt-6">
				<button type="button" on:click={start} class="btn-secondary">
					<span>Start</span>
				</button>
			</div>
		</div>
	{:else}
		<div class="w-1/6 h-screen bg-gray-950">
			<div class="users-container">
				{#each game.players as player}
					<User {player} />
				{/each}
			</div>
		</div>
		<div class="flex-1 h-screen overflow-hidden relative">
			<!-- TODO: Implement later: -->
			<!-- {#if game.roundNotifier === true}
				<div class="w-full h-full bg-black flex items-center justify-center z-50 absolute">
					<div class="text-white font-bold text-5xl text-center back-in-left-animation">
						Round <span class="text-yellow-400 text-9xl">{game.round}</span>
					</div>
				</div>
			{/if} -->
			<div class="notifier">
				<p>Pick your card</p>
				<div class="countdown">
					<div class="progress-bar-container">
						<div class="progress-bar red countdown-card-animation" />
					</div>
					<div id="js-timer" class="timer">10</div>
					<div class="progress-bar-container rotate-y-180">
						<div class="progress-bar red countdown-card-animation" />
					</div>
				</div>
			</div>

			<!-- <div class="table">
				<div class="card" />
				<div class="card" />
				<div class="card" />
				<div class="card picked-card-animation" />
			</div> -->

			<div class="cards-container">
				{#each game.cards as card, index}
					<div class="card dealing-card-animation" style="animation-delay: {index}s;" />
				{/each}
			</div>
		</div>
	{/if}
</div>