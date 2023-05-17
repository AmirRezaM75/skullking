<script lang="ts">
	async function createGame() {
		const response = await fetch('http://localhost:3000/games', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			}
		});

		const data = await response.json();
		const gameId = data.id;
		window.location.href = `games/${gameId}`;
	}

	export let data;
</script>

<div class="background">
	<div class="inner">
		<div>
			<p class="text-3xl text-white">Welcome to</p>
			<h1 class="title-primary">Skull King</h1>
			<div class="mt-10">
				{#if data.action == 'play'}
					<button type="button" on:click={createGame} class="btn-secondary">
						<span>Play</span>
					</button>
				{:else}
					<a href="/login" class="btn-secondary">Login</a>
				{/if}
			</div>
		</div>
	</div>
</div>
