<script lang="ts">
	import Navigation from '../../components/Navigation.svelte';
	import { IntendedGameId } from '../../constants';
	import ApiService from '../../services/ApiService';
	import AuthService from '../../services/AuthService';
	import type { User } from '../../types';
	// TODO: define ./@src/ as root
	// https://stackoverflow.com/questions/73754777/svelte-import-by-absolute-path-does-not-work

	let identifier: string = '';
	let password: string = '';
	let message: string = '';
	let loading = false;

	async function login() {
		loading = true;

		const apiService = new ApiService();
		const response = await apiService.login(identifier, password);

		const data = await response.json();

		if (response.status === 400) {
			message = data.message;
		}

		if (response.status === 200) {
			const authService = new AuthService();
			const user: User = {
				id: data.user.id,
				username: data.user.username,
				email: data.user.email,
				verified: data.user.verified,
				token: data.token
			};
			authService.save(user);

			if (user.verified) {
				const gameId = sessionStorage.getItem(IntendedGameId);
				sessionStorage.removeItem(IntendedGameId);
				window.location.href = gameId ? `/games/${gameId}` : '/';
			} else {
				window.location.href = 'verify-email';
			}
		}

		loading = false;
	}
</script>

<svelte:head>
	<title>Login</title>
</svelte:head>

<div class="w-screen h-screen flex items-center justify-center bg-slate-900">
	<Navigation />
	<div class="w-80">
		<h1 class="font-bold text-white text-3xl mb-8 text-center">Login</h1>

		<form on:submit={login}>
			<div class="mb-3">
				<label for="identifier">Username / Email</label>
				<input
					type="text"
					id="identifier"
					bind:value={identifier}
					autofocus
					required
				/>
			</div>

			<div class="mb-3">
				<label for="password">Password</label>
				<input
					type="password"
					id="password"
					bind:value={password}
					required
				/>
				{#if message}
					<small class="text-red-500">{message}</small>
				{/if}
			</div>

			<button type="submit" class="btn">
				{#if loading}
					<span class="circle-loader mr-2" />
				{/if}
				Login
			</button>

			<a href="/register" class="text-blue-400 hover:text-blue-300 block mt-3">Not registered?</a>

			<a href="/forgot-password" class="text-blue-400 hover:text-blue-300 inline-block">
				Forgot your password?
			</a>
		</form>
	</div>
</div>
