<script lang="ts">
	export let data;
</script>

<div class="w-screen h-screen flex items-center justify-center bg-slate-900">
	<div class="max-w-md">
		<h1 class="font-bold text-white text-3xl mb-8 text-center">Email Verification</h1>
		<div class="card">
			<p class="text-gray-100">
				Avast, me hearty! Ye be one step closer to joining our crew. We've just sent a message in a
				bottle to <span class="text-fuchsia-500">{data.email}</span>. If ye can't find it, check yer
				spam or junk folder. Thank ye for choosing to sail with us on the high seas of the internet.
			</p>
			<button type="button" class="btn mt-4">Resend</button>
		</div>
	</div>
</div>
