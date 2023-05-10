<script lang="ts">
	import AuthService from '../../services/AuthService';
	import type { User } from '../../types';
	// TODO: define ./@src/ as root
	// https://stackoverflow.com/questions/73754777/svelte-import-by-absolute-path-does-not-work
	import ServerValidationError from '../../utils/ServerValidationError';

	let username: string = '';
	let email: string = '';
	let password: string = '';
	let errors = new ServerValidationError();

	async function register() {
		const response = await fetch('http://localhost:3000/register', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({
				username,
				email,
				password
			})
		});

		// TODO: Refactor and use this format:
		// .then(response => response.json())
		// .then(data => {})
		const data = await response.json();

		if (response.status === 422) {
			Object.keys(data.errors).forEach((key) => {
				errors.add(key, data.errors[key]);
			});

			errors = errors;
		}

		if (response.status === 201) {
			const authService = new AuthService();
			const user: User = {
				username: data.user.username,
				email: data.user.email,
				verified: data.user.verified,
				token: data.token
			};
			authService.save(user);

			window.location.href = '/verify-email';
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
		<h1 class="font-bold text-white text-3xl mb-8 text-center">Register</h1>
		<form on:submit={register} on:keydown={clearError}>
			<div class="mb-3">
				<label for="email">Email</label>
				<input
					type="email"
					id="email"
					bind:value={email}
					class:border-red-500={errors.has('email')}
					autofocus
					required
				/>
				{#if errors.has('email')}
					<small class="text-red-500">{errors.get('email')}</small>
				{/if}
			</div>

			<div class="mb-3">
				<label for="username">Username</label>
				<input
					type="text"
					id="username"
					bind:value={username}
					class:border-red-500={errors.has('username')}
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
			</div>

			<button type="submit" class="btn">Join the crew</button>
		</form>
	</div>
</div>
