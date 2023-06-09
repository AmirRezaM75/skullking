<script lang="ts">
	import ApiService from '../../services/ApiService';
	import AuthService from '../../services/AuthService';
	import type { User } from '../../types';
	// TODO: define ./@src/ as root
	// https://stackoverflow.com/questions/73754777/svelte-import-by-absolute-path-does-not-work
	import ServerValidationError from '../../utils/ServerValidationError';

	let username: string = '';
	let password: string = '';
	let errors = new ServerValidationError();
	let message: string = '';

	async function login() {
		const apiService = new ApiService();
		const response = await apiService.login(username, password);

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
				window.location.href = '/';
			} else {
				window.location.href = 'verify-email';
			}
		}
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

<div class="w-screen h-screen flex items-center justify-center bg-slate-900">
	<div class="w-80 max-w-lg">
		<h1 class="font-bold text-white text-3xl mb-8 text-center">Login</h1>
		<form on:submit={login} on:keydown={clearError}>
			<div class="mb-3">
				<label for="username">Username</label>
				<input
					type="text"
					id="username"
					bind:value={username}
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

			<button type="submit" class="btn">Login</button>

			<a href="/forgot-password" class="text-gray-200 inline-block mt-3">Forgot your password?</a>
		</form>
	</div>
</div>
