<script lang="ts">
	import ApiService from '../services/ApiService.js';
	let loading = false;
	async function createGame() {
		if (loading == true) {
			return;
		}

		loading = true;
		const apiService = new ApiService();
		const response = await apiService.createGame();
		response.json().then((data) => {
			loading = false;
			window.location.href = `games/${data.id}`;
		});
	}

	export let data;
</script>

<svelte:head>
	<title>Skull King</title>
</svelte:head>

<div class="background">
	<div class="inner">
		<div>
			<p class="text-3xl text-white">Welcome to</p>
			<h1 class="title-primary">Skull King</h1>
			<div class="mt-10">
				{#if data.action == 'play'}
					<button type="button" on:click={createGame} class="btn-secondary" class:cursor-wait={loading}>
						{#if loading}
							<span class="circle-loader mr-2" />
						{/if}
						<span>Play</span>
					</button>
				{:else}
					<a href="/login" class="btn-secondary">Login</a>
				{/if}
			</div>
		</div>
	</div>
</div>
