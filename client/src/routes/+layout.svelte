<script>
	import '../app.scss';
</script>

<slot />

<svelte:head>
	{#if import.meta.env.VITE_ENV === 'production'}
		<script async src="https://www.googletagmanager.com/gtag/js?id=G-SMBNL31WB7"></script>
		<script>
			window.dataLayer = window.dataLayer || [];
			function gtag() {
				dataLayer.push(arguments);
			}
			gtag('js', new Date());

			gtag('config', 'G-SMBNL31WB7');
		</script>
	{/if}
</svelte:head>
