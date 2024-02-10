<script lang="ts">
	import Navigation from '../../components/Navigation.svelte';
	import { IntendedUrl } from '../../constants';
	import ApiService from '../../services/ApiService';
	import AuthService from '../../services/AuthService';
	import type { User } from '../../types';
	import { goto } from '$app/navigation';
	// TODO: define ./@src/ as root
	// https://stackoverflow.com/questions/73754777/svelte-import-by-absolute-path-does-not-work

	let identifier = '';
	let password = '';
	let message = '';
	let loading = false;

	async function login() {
		if (loading) {
			return;
		}

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
				avatarId: data.user.avatarId,
				token: data.token
			};
			authService.save(user);

			if (user.verified) {
				const intendedUrl = sessionStorage.getItem(IntendedUrl);
				sessionStorage.removeItem(IntendedUrl);
				goto(intendedUrl ?? '/lobbies');
			} else {
				goto('verify-email');
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
				<input type="text" id="identifier" bind:value={identifier} autofocus required />
			</div>

			<div class="mb-3">
				<label for="password">Password</label>
				<input type="password" id="password" bind:value={password} required />
				{#if message}
					<small class="text-red-500">{message}</small>
				{/if}
			</div>

			<button type="submit" class="btn {loading ? 'loading' : ''}"> Login </button>

			<a href="/register" class="text-blue-400 hover:text-blue-300 block mt-3">Not registered?</a>

			<a href="/forgot-password" class="text-blue-400 hover:text-blue-300 inline-block">
				Forgot your password?
			</a>
		</form>
	</div>
</div>
