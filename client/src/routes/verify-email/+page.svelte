<script lang="ts">
	import Navigation from '../../components/Navigation.svelte';
	import ApiService from '../../services/ApiService.js';
	let loading = false;

	let response: {
		success: boolean;
		message: string;
	} | null;
	response = null;

	async function resend() {
		if (loading) return;

		loading = true;

		const apiService = new ApiService();

		const r = await apiService.sendEmailVerificationNotification();

		loading = false;

		if (r.status == 202) {
			response = {
				success: true,
				message: 'Email verification notification has been sent.'
			};
		} else {
			response = {
				success: false,
				message: 'Something goes wrong!'
			};
		}
	}
	export let data;
</script>

<div class="w-screen h-screen flex items-center justify-center bg-slate-900">
	<Navigation />
	<div class="max-w-md">
		<h1 class="font-bold text-white text-3xl mb-8 text-center">Email Verification</h1>
		<div class="plate">
			<p class="text-gray-100">
				Avast, me hearty! Ye be one step closer to joining our crew. We've just sent a message in a
				bottle to <span class="text-fuchsia-500">{data.email}</span>. If ye can't find it, check yer
				spam or junk folder. Thank ye for choosing to sail with us on the high seas of the internet.
			</p>
			{#if response}
				<p class="mt-3 {response.success ? 'text-lime-primary' : 'text-red-500'}">
					{response.message}
				</p>
			{/if}

			{#if !response || !response.success}
				<button type="button" class="btn mt-4 {loading ? 'loading' : ''}" on:click={resend}>
					<span>Resend</span>
				</button>
			{/if}
			<p class="text-gray-400 mt-3">
				If you have already confirmed your email on other devices or browsers, please log out and
				log back in once more.
			</p>
		</div>
	</div>
</div>
