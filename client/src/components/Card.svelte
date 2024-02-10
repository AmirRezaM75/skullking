<script lang="ts">
	import type { Card } from '../types';
	let clazz = '';

	export let card: Card;
	export let delay: number;
	export let showCardOwner = false;
	export { clazz as class };
</script>

<!-- TODO: https://stackoverflow.com/questions/74974066 -->
<!-- svelte-ignore a11y-click-events-have-key-events -->
<div
	on:click
	class="card {card.borderColor} {clazz} {card.backgroundColor} {card.disabled ? 'disabled' : ''}"
	style="animation-delay: {delay / 2}s;"
>
	<div
		class="w-full bg-cover"
		style="height: {showCardOwner ? '90%' : '100%'}; background-image: url({card.imageURL});"
	/>
	<div
		class="number {card.backgroundColor} {card.textColor}"
		style={card.number === 0 ? 'display: none;' : ''}
	>
		{card.number}
	</div>
	{#if showCardOwner}
		<div class="card-owner {card.backgroundColor} {card.textColor}" style="height: 10%;">
			{card.ownerUsername}
		</div>
	{/if}
</div>
