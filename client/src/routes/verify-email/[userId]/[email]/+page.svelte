<script lang="ts">
	import { onMount } from 'svelte';
	import ApiService from '../../../../services/ApiService';
	import AuthService from '../../../../services/AuthService';

	let status = 'loading';

	onMount(async () => {
		const apiService = new ApiService
		const response = await apiService.verifyEmail(window.location.pathname + window.location.search)
        
        if (response.status === 200) {
            status = 'succeeded'
			const authService = new AuthService
			authService.markEmailAsVerified()
			setTimeout(() => window.location.href = '/', 1500)
        } else {
            status = 'failed'
        }
	});
</script>

<div class="w-screen h-screen flex items-center justify-center bg-slate-900">
	<div
		class="circle-loader"
		class:succeeded={status == 'succeeded'}
		class:failed={status == 'failed'}
	>
		<div class="checkmark" />
		<div class="cross" />
	</div>
</div>
