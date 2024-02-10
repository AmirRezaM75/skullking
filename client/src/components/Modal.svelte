<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import CloseIcon from './icons/CloseIcon.svelte';

	export let header: string;

	const dispatch = createEventDispatcher<{ closeModal: boolean }>();

	function close() {
		dispatch('closeModal', true);
	}
</script>

<div
	class="absolute top-0 w-full min-h-screen flex items-center justify-center px-2 sm:px-4 md:px-0"
	style="background-color: rgba(0, 0, 0, 1); z-index: 2"
>
	<div class="max-w-2xl">
		<div
			class="w-full rounded-t-lg bg-slate-700 p-4 text-gray-300 border-b border-gray-700 flex justify-between items-center"
		>
			<div class="font-bold uppercase">{header}</div>
			<!-- svelte-ignore a11y-click-events-have-key-events -->
			<div on:click={close}>
				<CloseIcon />
			</div>
		</div>
		<div class="bg-slate-800">
			<slot name="body" />
		</div>
		<div
			class="w-full rounded-b-lg p-4 border-t flex justify-between items-center border-gray-700 text-right bg-slate-800"
		>
			<slot name="footer" />
		</div>
	</div>
</div>
