<script lang="ts">
	import Navigation from '../../components/Navigation.svelte';
	import { IntendedGameId } from '../../constants';
	import ApiService from '../../services/ApiService';
	import AuthService from '../../services/AuthService';
	import type { User } from '../../types';
	// TODO: define ./@src/ as root
	// https://stackoverflow.com/questions/73754777/svelte-import-by-absolute-path-does-not-work
	import ServerValidationError from '../../utils/ServerValidationError';

	let identifier: string = '';
	let password: string = '';
	let errors = new ServerValidationError();
	let message: string = '';
	let loading = false;

	async function login() {
		loading = true;

		const apiService = new ApiService();
		const response = await apiService.login(identifier, password);

		const data = await response.json();

		if (response.status === 422) {
			Object.keys(data.errors).forEach((key) => {
				errors.add(key, data.errors[key]);
			});

			errors = errors;
		}

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

	function clearError(event: Event) {
		const element = event.target as HTMLInputElement;
		const id = element.getAttribute('id');
		if (id) {
			errors.clear(id);
			errors = errors;
		}
	}
</script>

<svelte:head>
	<title>Login</title>
</svelte:head>

<div class="w-screen h-screen flex items-center justify-center bg-slate-900">
	<Navigation />
	<div class="w-80">
		<h1 class="font-bold text-white text-3xl mb-8 text-center">Login</h1>

		<form on:submit={login} on:keydown={clearError}>
			<div class="mb-3">
				<label for="identifier">Username / Email</label>
				<input
					type="text"
					id="identifier"
					bind:value={identifier}
					class:border-red-500={errors.has('username')}
					autofocus
					required
				/>
				{#if errors.has('username')}
					<small class="text-red-500">{errors.get('username')}</small>
				{/if}
			</div>

			<div class="mb-3">
				<label for="password">Password</label>
				<input
					type="password"
					id="password"
					bind:value={password}
					class:border-red-500={errors.has('password')}
					required
				/>
				{#if errors.has('password')}
					<small class="text-red-500">{errors.get('password')}</small>
				{/if}

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
