<script lang="ts">
	import type { Player } from '../types';
	export let player: Player;
</script>

<div class="user">
	{#if player.roundStarter}
		<div class="roundStarterIndicator"></div>
	{/if}
	<div class="image">
		{#if player.picking}
			<div class="loader">
				<div class="circle" />
			</div>
		{/if}
		<img src={player.avatar} alt="" />
	</div>
	<div class="info">
		<div class="username-container">
			<div
				class="connection-status {player.isConnected ? 'bg-green-500' : 'bg-red-500'}"
				title={player.isConnected ? 'Online' : 'Offline'}
			/>
			<div class="username">{player.username}</div>
		</div>
		<div class="action">
			<span title="Bid" class="mr-2">B: {player.bid}</span>
			<span title="Won tricks count">W: {player.wonTricksCount}</span>
		</div>
	</div>
	<div class="score">{player.score}</div>
</div>

<style lang="scss">
	.roundStarterIndicator {
		@apply w-full h-full absolute;
		left: -1.25rem;
		background: linear-gradient(271deg, rgba(255,255,255,0) 0%, rgb(75 255 255 / 40%) 95%, rgb(75 255 255 / 40%) 100%);
	}
</style>