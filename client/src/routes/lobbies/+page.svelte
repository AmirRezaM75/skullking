<script lang="ts">
	import AnimatedBackground from '../../components/AnimatedBackground.svelte';
import Navigation from '../../components/Navigation.svelte';
</script>

<svelte:head>
	<title>Lobbies - Kenopsia</title>
</svelte:head>

<section class="min-h-screen flex justify-center bg-slate-900">
	<Navigation />
	<AnimatedBackground/>
	<div class="w-full max-w-2xl px-4 2xl:px-0 mt-24" style="z-index: 2;">
		<div class="bg-slate-800 shadow-lg rounded-lg">
			<div class="px-4 py-3">
				<div class="flex flex-wrap items-center">
					<div class="flex-1">
						<h3 class="font-semibold text-base text-gray-100">Lobbies</h3>
					</div>
					<div class="flex-1 text-right">
						<button
							class="bg-indigo-500 text-white active:bg-indigo-600 font-bold uppercase px-3 py-1 rounded outline-none focus:outline-none ease-linear transition-all duration-150"
							type="button"
						>
							Create
						</button>
					</div>
				</div>
			</div>

			<div class="block w-full overflow-x-auto">
				<table class="items-center bg-slate-800 w-full border-collapse rounded-b-lg">
					<thead>
						<tr class="align-middle text-xs uppercase whitespace-nowrap bg-slate-900 border border-solid border-gray-600 border-l-0 border-r-0 text-gray-400 font-semibold">
							<th class="px-4 py-3 text-left">Lobby Name</th>
							<th class="px-4 py-3 text-center">Players</th>
							<th class="px-4 py-3 text-center">Created At</th>
							<th class="px-4 py-3"></th>
						</tr>
					</thead>

					<tbody>
						{#each {length: 5} as _, i}
						<tr class="text-gray-200">
							<td
								class="align-middle whitespace-nowrap p-4 text-left"
							>
								Intelligent Skull King
							</td>
							<td
								class="align-middle whitespace-nowrap p-4 text-center"
							>
								7
							</td>
							<td
								class="align-center whitespace-nowrap p-4 text-center"
							>
								2023-12-01 12:10:30
							</td>
							<td
								class="align-middle whitespace-nowrap p-4 text-right"
							>
								<a
									href="/"
									class="border border-lime-primary text-lime-primary hover:text-slate-800 hover:bg-lime-primary active:bg-lime-primary font-bold uppercase px-3 py-1 rounded outline-none focus:outline-none ease-linear transition-all duration-150"
									type="button"
								>
									Join
								</a>
							</td>
						</tr>
						{/each}
						
					</tbody>
				</table>
			</div>
		</div>
	</div>
</section>
