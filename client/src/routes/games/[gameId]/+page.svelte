<script lang="ts">
	import GameService from '../../../services/GameService.js';
	import { GameState, GameCommand } from '../../../constants.js';
	import User from '../../../components/User.svelte';
	import Card from '../../../components/Card.svelte';
	import Countdown from '../../../components/Countdown.svelte';
	import QueueService from '../../../services/QueueService.js';

	export let data;

	let game = new GameService(data.cardService, data.authId);
	const queue = new QueueService();

	const ws = new WebSocket(
		'ws://localhost:3000/games/join?gameId=' + data.gameId + '&token=' + data.token
	);

	ws.onopen = function (e) {
		console.log('OPEN');
	};

	ws.onmessage = async function (e) {
		let message = JSON.parse(e.data);
		queue.push(message);
		await run();
	};

	async function run() {
		const message = queue.pop();

		if (message) {
			console.log(message); // TODO: Remove
			game = await game.handle(message.command, message.content);
			queue.isProcessing = false;
			run();
		}
	}

	function start() {
		ws.send(
			JSON.stringify({
				command: GameCommand.Start
			})
		);
	}

	function bid(bid: number) {
		ws.send(
			JSON.stringify({
				command: GameCommand.Bid,
				content: bid.toString()
			})
		);
	}

	function pick(cardId: number) {
		ws.send(
			JSON.stringify({
				command: GameCommand.Pick,
				content: cardId.toString()
			})
		);
	}
</script>

<svelte:head>
	<title>Skull King</title>
	{#each ['chest', 'escape', 'king', 'kraken', 'map', 'mermaid', 'parrot', 'pirate', 'roger', 'whale'] as card}
		<!-- This will cause the browser to preload the images as soon as the page loads,
			even if they are not actually present in the DOM yet! -->
		<link rel="preload" href="/images/cards/{card}.jpg" as="image" />
	{/each}
</svelte:head>

<div
	class="min-w-full min-h-screen flex items-center justify-center"
	style="background-color: #1B1B1B;"
>
	{#if game.state == GameState.Pending}
		<div class="flex-col">
			<div class="flex items-center justify-center gap-4 flex-wrap px-2 py-4 max-w-2xl">
				{#each game.players as player}
					<div class="">
						<div class="bg-gray-700 p-6 rounded-lg text-center">
							<div class="mb-3">
								<img src={player.avatar} width="100" height="100" alt="" class="rounded-full" />
							</div>
							<span class="text-gray-100 font-bold text-lg uppercase">{player.username}</span>
						</div>
					</div>
				{/each}
			</div>
			{#if game.creator.id === data.authId}
				<div class="text-center mt-6">
					<button type="button" on:click={start} class="btn-secondary">
						<span>Start</span>
					</button>
				</div>
			{:else}
				<p class="text-yellow-500 text-center">
					Wait for {game.creator.username} to start the game.
				</p>
			{/if}
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
				<p>{game.notifierMessage}</p>
				{#if game.showCountdown}
					<Countdown number={game.timer} color={game.countdownColor} />
				{/if}
			</div>

			<div class="bids-container">
				{#each game.bids as bidNumber, index}
					<div
						class="bid
						{game.bid === bidNumber ? 'active' : ''}
						{game.state === GameState.EndBidding ? 'fade-in-down-animation animation-duration-500 active' : ''}"
						on:click={() => bid(bidNumber)}
						on:keydown={() => bid(bidNumber)}
					>
						{bidNumber}
					</div>
				{/each}
			</div>

			<div class="table">
				{#each game.tableCards as card}
					<Card {card} delay={0} class="picked-card-animation {card.isWinner ? 'winner' : ''}" />
				{/each}
			</div>

			<div class="cards-container">
				{#each game.cards as card, index}
					<Card
						{card}
						delay={index}
						class="dealing-card-animation"
						on:click={() => pick(card.id)}
					/>
				{/each}
			</div>
		</div>
	{/if}
</div>
