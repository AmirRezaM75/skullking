<script lang="ts">
	let username: string = '';
	let email: string = '';
	let password: string = '';

	async function register() {
		const response = await fetch('http://localhost:3000/register', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify({
				username,
				email,
				password
			})
		});
	}
</script>

<!-- #6DFA84 -->

<div class="w-screen h-screen flex items-center justify-center bg-gray-900">
	<div class="w-80 max-w-lg">
		<h1 class="font-bold text-white text-3xl mb-8 text-center">Register</h1>
		<form on:submit={register}>
			<div class="mb-3">
				<label for="email" class="block mb-2 text-sm font-medium text-gray-300">Email</label>
				<input type="email" id="email" bind:value={email} autofocus required />
			</div>
			<div class="mb-3">
				<label for="username" class="block mb-2 text-sm font-medium text-gray-300">Username</label>
				<input type="text" id="username" bind:value={username} required />
			</div>

			<div class="mb-3">
				<label for="password" class="block mb-2 text-sm font-medium text-gray-300">Password</label>
				<input type="password" id="password" bind:value={password} required />
			</div>

			<button type="submit">Register</button>
		</form>
	</div>
</div>
