<script lang="ts">
	import AuthService from '../services/AuthService';

	const authService = new AuthService();

	let isLoggedIn = authService.isLoggedIn();

	function logout() {
		authService.logout();
		window.location.reload()
	}
</script>

<!-- z-index is neccessary in combination of animated background -->
<nav style="z-index: 2;">
	<div class="inner">
		<div class="left">
			<a href="/">Home</a>
			<a href="https://github.com/AmirRezaM75/skull-king" target="_blank">Github</a>
			<a href="https://youtu.be/MOC6pa87pis?t=126" target="_blank">Tutorial</a>
			<a href="https://zarinp.al/@amirrezam75" class="" target="_blank">Donate</a>
		</div>
		<div class="flex items-center">
			{#if isLoggedIn}
				<button type="button" class="btn danger-gradient" on:click={logout}>Logout</button>
			{:else}
				<a href="/register" class="mr-4 hidden sm:inline-block">Register</a>
				<a href="/login" class="btn">Login</a>
			{/if}
		</div>
	</div>
</nav>

<style lang="scss">
	nav {
		@apply bg-slate-900 w-full h-16 px-4 absolute top-0 left-0;

		.inner {
			@apply max-w-7xl flex items-center justify-between h-full mx-auto;
		}

		a:not(.btn) {
			@apply text-gray-200 sm:text-lg hover:text-blue-300 pr-4 sm:px-4;
		}
	}
</style>
