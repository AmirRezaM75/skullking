<script lang="ts">
	import type { Player } from '../types';
	export let player: Player;
</script>

<div class="user">
	<div class="image">
		{#if player.picking}
			<div class="loader">
				<div class="circle" />
			</div>
		{/if}
		<img src={player.avatar} alt="" />
	</div>
	<div class="info">
		<div class="username">{player.username}</div>
		<div class="action">
			<span title="Bid" class="mr-2">B: {player.bid}</span>
			<span title="Won tricks count">W: {player.wonTricksCount}</span>
		</div>
	</div>
	<div class="score">{player.score}</div>
</div>
