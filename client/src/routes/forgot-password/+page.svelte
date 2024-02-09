<script lang="ts">
	import Navigation from '../../components/Navigation.svelte';
	import ApiService from '../../services/ApiService';

	let email = '';
	let message = '';
	let sent = false;
	let loading = false;

	async function submit() {
		if (loading) return;

		loading = true;

		const apiService = new ApiService();

		const response = await apiService.forgotPassword(email);

		loading = false;

		if (response.status == 202) {
			sent = true;
		} else {
			response.json().then((data) => (message = data.message));
		}
	}
</script>

<svelte:head>
	<title>Forgot Password</title>
</svelte:head>

<div class="w-screen h-screen flex items-center justify-center bg-slate-900">
	<Navigation />
	{#if !sent}
		<div class="max-w-sm">
			<h1 class="font-bold text-white text-3xl mb-8 text-center">Forgot Password?</h1>
			<div>
				<p class="text-gray-100">
					Don't worry it happens. please enter the address associated with your account.
				</p>
				<form on:submit={submit}>
					<div class="mt-3">
						<label for="email">Email</label>
						<input type="email" id="email" bind:value={email} required />
					</div>

					{#if message != ''}
						<small class="text-red-500">{message}</small>
					{/if}

					<button type="submit" class="btn mt-4 {loading ? 'loading' : ''}">
						<span>Submit</span>
					</button>
				</form>
			</div>
		</div>
	{:else}
		<div class="max-w-md">
			<div class="plate">
				<img src="/images/checked.png" class="mx-auto" width="85" height="85" alt="" />

				<p class="text-gray-200 mt-4">
					We will send you an email with instructions on how to reset your password. Please check
					your email inbox and spam folder. If you don't receive an email within a few minutes,
					please try again or contact our support team for assistance.
				</p>
			</div>
		</div>
	{/if}
</div>
