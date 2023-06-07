<script lang="ts">
	import type { Card } from '../types';
	export let card: Card;
	export let index: number;
</script>

<div
	class="card dealing-card-animation {card.borderColor}"
	style="background-image: url({card.imageURL}); animation-delay: {index}s;"
>
	<div
		class="number {card.backgroundColor}"
		style="{card.number === 0 ? 'display: none;' : ''}"
	>
		{card.number}
	</div>
</div>
