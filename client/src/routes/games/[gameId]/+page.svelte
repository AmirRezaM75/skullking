<script lang="ts">
	import GameService from '../../../services/GameService.js';

	export let data;

	let game = new GameService

	const ws = new WebSocket(
		'ws://localhost:3000/games/join?gameId=' + data.gameId + '&token=' + data.token
	);

	ws.onopen = function (e) {
		console.log('OPEN');
	};

	ws.onmessage = function (e) {
		let message = JSON.parse(e.data);
		
		game = game.handle(message['command'], message['content'], message['senderId'])
	};
</script>

<div
	class="min-w-full min-h-screen bg-gray-900 flex items-center justify-center"
	style="background-color: #1B1B1B;"
>
	<div class="flex-col">
		<div class="flex items-center justify-center gap-4 flex-wrap px-2 py-4 max-w-2xl">
			{#each game.players as player}
			<div class="">
				<div class="bg-gray-700 p-6 rounded-lg text-center">
					<div class="mb-3">
						<img src="/images/avatars/1.jpg" width="100" height="100" alt="" class="rounded-full" />
					</div>
					<span class="text-gray-100 font-bold text-lg uppercase">{player.username}</span>
				</div>
			</div>
			{/each}
		</div>
		<div class="text-center mt-6">
			<button type="button" class="btn-secondary">
				<span>Start</span>
			</button>
		</div>
	</div>
</div>
