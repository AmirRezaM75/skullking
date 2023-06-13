<script lang="ts">
	import type { Card } from '../types';
	let clazz = '';

	export let card: Card;
	export let delay: number;
	export {clazz as class}; 
</script>

<!-- TODO: https://stackoverflow.com/questions/74974066 -->
<!-- svelte-ignore a11y-click-events-have-key-events -->
<div
	on:click
	class="card {card.borderColor} {clazz} {card.disabled ? 'disabled' : ''}"
	style="background-image: url({card.imageURL}); animation-delay: {delay}s;"
>
	<div
		class="number {card.backgroundColor} {card.textColor}"
		style={card.number === 0 ? 'display: none;' : ''}
	>
		{card.number}
	</div>
</div>
