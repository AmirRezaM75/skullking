<script lang="ts">
	import { onMount } from 'svelte';

	let status = 'loading';

	onMount(async () => {
        const URL = 'http://localhost:3000' + window.location.pathname + window.location.search
		const response = await fetch(URL, {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json'
			}
		});

        
        if (response.status === 200) {
            status = 'succeeded'
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
