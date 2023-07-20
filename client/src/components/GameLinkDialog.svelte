<script lang="ts">
	let copied = false;
	async function copyLink() {
		if (copied) return;
		try {
			navigator.clipboard.writeText(window.location.href).then(() => {
				copied = true;
				setTimeout(() => {
					copied = false;
				}, 4000);
			});
		} catch (err) {
			copied = false;
		}
	}
</script>

<div class="dialog-box">
	<p class="text-white">
		Invite at least one more player by sharing this link with the people you want to play.
	</p>
	<p class="text-blue-400 truncate">{window.location.href}</p>
	<div class="mt-3">
		{#if copied}
			<p class="text-lime-primary py-2.5">The link has been copied successfully</p>
		{:else}
			<button type="button" class="btn" on:click={copyLink}>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					width="20"
					height="20"
					class="mr-2"
					data-name="Group 27733"
					viewBox="0 0 230.451 235.429"
					><path
						fill="#2e2842"
						d="M86.832 63.594h115.407a28.247 28.247 0 0 1 28.212 28.216v115.4a28.247 28.247 0 0 1-28.212 28.216H86.832a28.247 28.247 0 0 1-28.211-28.216V91.81a28.247 28.247 0 0 1 28.211-28.216Zm115.407 157.431a13.826 13.826 0 0 0 13.808-13.812v-115.4a13.826 13.826 0 0 0-13.808-13.812H86.832a13.826 13.826 0 0 0-13.807 13.813v115.4a13.826 13.826 0 0 0 13.807 13.813Z"
						data-name="Path 25"
					/><path
						fill="#2e2842"
						d="M48.077 0h89.5a34.289 34.289 0 0 1 34.255 34.251v3.072a7.2 7.2 0 1 1-14.4 0v-3.072a19.875 19.875 0 0 0-19.853-19.848h-89.5a33.716 33.716 0 0 0-33.674 33.674v89.506a19.873 19.873 0 0 0 19.848 19.848 7.2 7.2 0 0 1 0 14.4A34.291 34.291 0 0 1 .001 137.58V48.078A48.131 48.131 0 0 1 48.077 0Z"
						data-name="Path 26"
					/></svg
				>
				Copy link
			</button>
		{/if}
	</div>
</div>

<style lang="scss">
	.dialog-box {
		@apply absolute bottom-5 left-5 right-5 p-5 max-w-md rounded-md bg-slate-700 text-lg;

		button {
			background: linear-gradient(45deg, aqua, var(--lime-primary));
		}
	}
</style>
