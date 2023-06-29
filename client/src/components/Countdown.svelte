<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte';
	import type { Countdown } from '../types';
	const dispatch = createEventDispatcher<{ newCountdown: Countdown }>();
	export let number: number;
	export let color: 'red' | 'blue';

	$: timer = number;

	let isAudioPlaying = false;

	onMount(() => {
		const audio = new Audio('/sounds/countdown.mp3');
		audio.volume = 0.5;

		let intervalId = setInterval(() => {
			timer -= 1;

			if (color === 'blue' && timer <= 5) {
				// The following code may be executed up to five times due to the condition timer <= 5.
				// While using timer == 5 could prevent this, it may not work if a player joins in the nick of time of picking after disconnecting,
				// causing the audio to not play. Additionally, there is a risk of the audio being paused when a player picks a card,
				// while it's still being played in the setInterval cycle.
				// This can result in an error message in Chrome, stating "The play() request was interrupted".
				if (!isAudioPlaying) {
					audio
						.play()
						.then(() => (isAudioPlaying = true))
						.catch(() => (isAudioPlaying = false));
				}
			}

			if (timer === 0) {
				if (color === 'blue') {
					audio.pause();
				}
				clearInterval(intervalId);
			}
		}, 1000);

		dispatch('newCountdown', { id: intervalId, audio: audio });
	});
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
