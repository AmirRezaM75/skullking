<script lang="ts">
	let email = '';
	let message = '';
	let sent = false;

	async function submit() {
		const response = await fetch('http://localhost:3000/forgot-password', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({
				email
			})
		});

		if (response.status == 202) {
			sent = true;
		} else {
			const data = await response.json();
			message = data.message;
		}
	}
</script>

<div class="w-screen h-screen flex items-center justify-center bg-slate-900">
	{#if !sent}
		<div class="max-w-md">
			<h1 class="font-bold text-white text-3xl mb-8 text-center">Forgot Password?</h1>
			<div class="card">
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

					<button type="submit" class="btn mt-4">Submit</button>
				</form>
			</div>
		</div>
	{:else}
		<div class="max-w-md">
			<div class="card">
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
