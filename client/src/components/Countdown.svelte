<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	const dispatch = createEventDispatcher();
	export let number: number;
	export let color: 'red' | 'blue';

	$: timer = number;

	let x = setInterval(() => {
		timer -= 1;
		if (timer == 0) {
			clearInterval(x);
		}
	}, 1000);
</script>

<div class="countdown">
	<div class="progress-bar-container">
		<div class="progress-bar {color} countdown-animation" style="animation-duration: {number}s;" />
	</div>
	<div id="js-timer" class="timer">{timer}</div>
	<div class="progress-bar-container rotate-y-180">
		<div class="progress-bar {color} countdown-animation" style="animation-duration: {number}s;" />
	</div>
</div>
