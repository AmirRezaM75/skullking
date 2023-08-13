<script lang="ts">
	export let message: string
</script>

<div class="fluid-container">
	<div class="dialog-box">
		<p class="text-xl text-red-400 font-semibold uppercase mb-2">Error</p>
		<p class="text-gray-100">{message}<p>
	</div>
</div>

<style lang="scss">
	.dialog-box {
		@apply p-5 max-w-md rounded-md bg-slate-700 text-lg;
	}
</style>
