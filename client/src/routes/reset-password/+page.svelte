<script lang="ts">
	import ApiService from '../../services/ApiService';
	import ServerValidationError from '../../utils/ServerValidationError';

	let password = '';
	let message = '';

	let errors = new ServerValidationError();

	async function submit(e: Event) {
		e.preventDefault();
		const params = new URLSearchParams(window.location.search);
		const email = params.get('email');
		const token = params.get('token');

		if (!email || !token) {
			return;
		}

		const apiService = new ApiService();
		const response = await apiService.resetPassword(email, token, password);

		if (response.status == 200) {
			window.location.href = '/login';
			return;
		}

		const data = await response.json();

		if (response.status === 422) {
			Object.keys(data.errors).forEach((key) => {
				errors.add(key, data.errors[key]);
			});

			errors = errors;
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

<svelte:head>
	<title>Reset Password</title>
</svelte:head>

<div class="w-screen h-screen flex items-center justify-center bg-slate-900">
	<div class="max-w-md">
		<h1 class="font-bold text-white text-3xl mb-8 text-center">Reset Password</h1>
		<div class="plate">
			<p class="text-gray-100">
				After you submit your new password, you will be redirected to the login page where you can
				log in with your new credentials.
			</p>
			<form on:submit={submit} on:keydown={clearError}>
				<div class="mt-3">
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

				{#if message != ''}
					<small class="text-red-500">{message}</small>
				{/if}

				<button type="submit" class="btn mt-4">Submit</button>
			</form>
		</div>
	</div>
</div>
