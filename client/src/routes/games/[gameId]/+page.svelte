<script lang="ts">
	import GameService from '../../../services/GameService.js';
	import { GameState, GameCommand } from '../../../constants.js';
	import User from '../../../components/User.svelte';
	import Card from '../../../components/Card.svelte';
	import Countdown from '../../../components/Countdown.svelte';
	import QueueService from '../../../services/QueueService.js';
	import Swiper from 'swiper';
	import 'swiper/css';
	import { onMount } from 'svelte';

	export let data;

	let game = new GameService(data.cardService, data.authId);
	const queue = new QueueService();
	game.test([1, 2, 3, 4,6,5,7,8,9,10]);
	const ws = new WebSocket(
		'ws://localhost:3000/games/join?gameId=' + data.gameId + '&token=' + data.token
	);

	onMount(() => {
		new Swiper('.swiper', {
			slidesPerView: 'auto',
			spaceBetween: 10,
		});
	});

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
	{#if game.state != GameState.Pending}
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
			{#if game.creator.id !== data.authId}
				<p class="text-yellow-500 text-center">
					Wait for {game.creator.username} to start the game.
				</p>
			{:else if game.players.length === 1}
				<p class="text-yellow-500 text-center">Invite at least one more player to start the game</p>
			{:else}
				<div class="text-center mt-6 mb-6 sm:mb-0">
					<button type="button" on:click={start} class="btn-secondary">
						<span>Start</span>
					</button>
				</div>
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
				<div class="swiper w-full">
					<div class="swiper-wrapper">
						{#each game.cards as card, index}
							<Card
								{card}
								delay={index}
								class="dealing-card-animation swiper-slide"
								on:click={() => pick(card.id)}
							/>
						{/each}
					</div>
				</div>
			</div>
		</div>
	{/if}
</div>
