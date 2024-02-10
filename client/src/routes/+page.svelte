<script lang="ts">
	import Navigation from '../components/Navigation.svelte';

	export let data;
</script>

<svelte:head>
	<title>Skull King Online</title>
	<meta
		name="description"
		content="Play Skull King online for free! Join the fun with this fast-paced card game that's perfect for game night with friends and family."
	/>
	<meta name="keywords" content="skull king, SkullKing, skull king online" />
</svelte:head>

<Navigation />

<div class="background">
	<div class="inner">
		<div>
			<p class="text-3xl text-white">Welcome to</p>
			<h1 class="title-primary">Skull King</h1>
			<div class="mt-10">
				{#if data.isUserVerified}
					<a class="btn-secondary" href="/lobbies">Play</a>
				{:else}
					<a href="/login" class="btn-secondary">Login</a>
				{/if}
			</div>
		</div>
	</div>
</div>
