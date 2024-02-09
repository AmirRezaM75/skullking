<script lang="ts">
	export let message: string;
	export let label = 'Error';
	export let errorCode = 0;
</script>

<div class="exception-container">
	<div class="dialog-box">
		<p class="text-xl text-red-400 font-semibold uppercase mb-2">{label}</p>
		<p class="text-gray-100">{message}</p>
		<p>
			<slot />
			{#if errorCode !== 0}
				<a href="/" class="btn light-gradient mt-4">
					<img
						width="18"
						class="mr-2"
						src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAADAAAAAwCAYAAABXAvmHAAAAAXNSR0IArs4c6QAAAhxJREFUaEPtms0uBEEQx38OPk6IE+I5cJNIJG7OEryCxMVHEB8RHIRnQDyAE04SFxfeAieEGxfyl91kjOnZ6Zmu3ZFMH3d7q/6/6qpKd9GGzeoE9oEZ4As4BZaAj9Du2kIbBLqBc2AsZvsamALeQ/oMDSDxl8CoQ+QdMAG8hoIICdADXKSIr2u+BSaBtwhEB7ANzAIDkc8fgRNgHfhMgg4FkFW8C2IPWEw5FX2/bAXgK76uI5pOinQ08nGtT8CgBUBe8XGIlww1kZgtRVJI4lWwIxmcp21RTbiKPvq7oABFI5+HORhAqMj7QgQBaJV4wRYGaKX4wgCtFl8IoAzicwOURXwugDKJ9waQ+Ctg2LfXGe7vjV0Af1wltaayRb4ek/vaVfzXtSMOUFbxTogoQNnFJ0LUAf6L+D8QAmivvaTGDQvQwvSNakIAR8C8hYcm2DwUwDPQ1wRnFi6eBaDnWr+F9SbYfBTAAbDQBGcWLg4E0AXsANMNHtYWAvLaVNacAatpb2KNBMu0vB80FUDg4zM/AY3+1mqTaA2qkpaGUxofbgEaJ/oscwCNzzVSz7K0dzfLxsgecwBFV90hy9IY0XVKrt+bA/hO+XybRAXQKDWqE2gQoSqFqhSKRaDqQgpIyNto1YWqLmRRZClRNS/iB9ffZh2iLGtAWoaS/KY51XVX196syxJAWlZ8AfTg2ATmMp6EBYAifwxsuP5X4hsd0nMrxs6iQAAAAABJRU5ErkJggg=="
						alt=""
					/>
					<span>Home</span>
				</a>
			{/if}
		</p>
	</div>
</div>

<style lang="scss">
	.dialog-box {
		@apply p-5 max-w-md rounded-md bg-slate-700 text-lg;
	}

	.exception-container {
		@apply w-full h-full absolute bg-black bg-opacity-70 flex items-center justify-center;
		z-index: 3;
	}
</style>
