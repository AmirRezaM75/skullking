<script lang="ts">
	import GameService from '../../../services/GameService';
	import { GameState, GameCommand } from '../../../constants';
	import User from '../../../components/User.svelte';
	import Card from '../../../components/Card.svelte';
	import Countdown from '../../../components/Countdown.svelte';
	import QueueService from '../../../services/QueueService';
	import AudioService from '../../../services/AudioService';
	import Swiper from 'swiper';
	import 'swiper/css';
	import { onMount } from 'svelte';
	import ApiService from '../../../services/ApiService';
	import type { Countdown as CountdownType } from '../../../types';
	import AudioIcon from '../../../components/AudioIcon.svelte';
	import ConnectionErrorDialog from '../../../components/ConnectionErrorDialog.svelte';
	import ExceptionReporter from '../../../components/ExceptionReporter.svelte';
	import Loader from '../../../components/Loader.svelte';

	export let data;

	let isSidebarOpen = true;
	let isBackgroundAudioPlaying = false;
	let disconnected = false;

	let countdowns: CountdownType[] = [];

	let game = new GameService(data.cardService, data.authId);

	const queue = new QueueService();

	const apiService = new ApiService();

	const audioService = new AudioService();

	let deckSwiper: Swiper;
	let tableSwiper: Swiper;

	let bidFunc: (bid: number) => void;
	let pickFunc: (bid: number) => void;

	onMount(async () => {
		let ticketId = '';

		const response = await apiService.createTicket();

		if (response.status === 201) {
			const data = await response.json();
			ticketId = data.id;
		}

		const ws = apiService.joinGame(data.gameId, ticketId);

		ws.onopen = function () {
			toggleBackgroundAudio();
			disconnected = false;
		};

		ws.onclose = function () {
			disconnected = true;
		};

		ws.onmessage = async function (e) {
			try {
				let message = JSON.parse(e.data);
				queue.push(message);
				await run();
			} catch (error) {
				// If responses is empty or malformed we don't want to stop executing
			}
		};

		bidFunc = (bid: number) =>
			ws.send(
				JSON.stringify({
					command: GameCommand.Bid,
					content: bid.toString()
				})
			);

		pickFunc = (cardId: number) =>
			ws.send(
				JSON.stringify({
					command: GameCommand.Pick,
					content: cardId.toString()
				})
			);

		deckSwiper = new Swiper('.deck-swiper', {
			slidesPerView: 'auto'
		});

		tableSwiper = new Swiper('.table-swiper', {
			slidesPerView: 'auto'
		});

		isSidebarOpen = screen.width > 640;

		[
			'announceScores',
			'announceTrickWinner',
			'background',
			'countdown',
			'picked',
			'startPicking'
		].forEach((filename) => {
			const audio = new Audio(`/sounds/${filename}.mp3`);
			audio.preload = 'auto';
		});
	});

	async function toggleBackgroundAudio() {
		audioService.toggleBackgroundAudio().then(() => {
			isBackgroundAudioPlaying = audioService.isBackgroundAudioPlaying;
		});
	}

	function keyboardHandler(event: KeyboardEvent) {
		if (event.code === 'KeyM') {
			toggleBackgroundAudio();
			return;
		}

		if (event.code === 'KeyS') {
			toggleSidebar();
			return;
		}

		const isDigit = /^\d$/.test(event.key);

		if (isDigit && game.state === GameState.Bidding) {
			bid(parseInt(event.key));
			return;
		}

		if (isDigit && game.findPickingPlayerId() === data.authId) {
			const card = game.cards[parseInt(event.key) - 1];
			if (card) {
				pick(card.id);
			}
			return;
		}
	}

	async function run() {
		const message = queue.pop();

		if (message) {
			game = await game.handle(message.command, message.content);

			game.postHandler(message.command, deckSwiper, tableSwiper);

			if (GameCommand.Picked === message.command) {
				countdowns.forEach((countdown) => {
					clearInterval(countdown.id);
					countdown.audio.pause();
				});
				countdowns = [];
			}

			queue.isProcessing = false;
			run();
		}
	}

	function addCountdown(event: CustomEvent<CountdownType>) {
		countdowns.push(event.detail);
	}

	function toggleSidebar() {
		isSidebarOpen = !isSidebarOpen;
	}

	function bid(bid: number) {
		bidFunc(bid);
	}

	function pick(cardId: number) {
		pickFunc(cardId);
	}
</script>

<!-- We should not prevent default keydown, because user won't be able to use some useful keybindings like reload page. -->
<svelte:window on:keydown={keyboardHandler} />

<svelte:head>
	<title>Skull King</title>
	{#each ['chest', 'escape', 'king', 'kraken', 'map', 'mermaid', 'parrot', 'pirate', 'roger', 'whale'] as card}
		<!--
			This will cause the browser to preload the images as soon as the page loads,
			even if they are not actually present in the DOM yet!
			as="audio" is not supported in chrome that's why we didn't use the same technique for audio files.
		-->
		<link rel="preload" href="/images/cards/{card}.jpg" as="image" />
	{/each}
</svelte:head>

<div class="board">
	{#if disconnected}
		<ConnectionErrorDialog />
	{/if}
	{#if game.exceptionMessage !== ''}
		<ExceptionReporter message={game.exceptionMessage} errorCode={404} />
	{/if}
	{#if game.state == GameState.Pending}
		<Loader />
	{:else}
		<div
			class="sidebar relative {isSidebarOpen
				? 'open-sidebar-animation'
				: 'close-sidebar-animation'}"
		>
			<div class="users-container">
				{#each game.players as player}
					<User {player} />
				{/each}
			</div>
			<div
				class="absolute bottom-4 left-4 cursor-pointer"
				title="{isBackgroundAudioPlaying ? 'Mute' : 'Unmute'} background music"
				on:click={toggleBackgroundAudio}
				on:keydown={keyboardHandler}
			>
				<AudioIcon color={isBackgroundAudioPlaying ? 'white' : 'gray'} />
			</div>
		</div>
		<div class="flex-1 h-screen overflow-hidden relative">
			<div on:click={toggleSidebar} on:keydown={keyboardHandler} class="sidebar-button">
				<img
					width="20"
					src="/images/arrow-{isSidebarOpen ? 'left' : 'right'}.png"
					alt="Toggle Sidebar"
				/>
			</div>

			{#if game.state === GameState.EndGame}
				<div class="w-full h-full flex flex-col items-center justify-center">
					<a href="/lobbies/{game.lobbyId}" class="btn btn-secondary max-w-fit">Play Again</a>
					<a href="/lobbies" class="block text-white mt-4 uppercase text-sm">Lobbies</a>
				</div>
			{:else}
				<div class="notifier">
					<p>{game.notifierMessage}</p>
					{#if game.showCountdown}
						<Countdown
							number={game.timer}
							color={game.countdownColor}
							on:newCountdown={addCountdown}
						/>
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

				<div class="table-container {game.table.hasWinner ? 'has-winner' : ''}">
					<div class="table-swiper w-full">
						<div class="swiper-wrapper">
							{#each game.table.cards as card}
								<Card
									showCardOwner={true}
									{card}
									delay={0}
									class="picked-card-animation swiper-slide {card.isWinner ? 'winner' : ''}"
								/>
							{/each}
						</div>
					</div>
				</div>

				<div class="cards-container">
					<div class="deck-swiper w-full">
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
			{/if}
		</div>
	{/if}
</div>
