<script lang="ts">
	let isOpen = false;

	function open() {
		isOpen = true;
	}

	function close() {
		isOpen = false;
	}
</script>

<div class="w-full h-full absolute overflow-hidden">
	<svg
		xmlns="http://www.w3.org/2000/svg"
		viewBox="0 0 50 50"
		width="35px"
		height="35px"
		class="absolute bottom-10 right-10 cursor-pointer"
		on:click={open}
	>
		<path
			fill="#ddd"
			d="M47.16,21.221l-5.91-0.966c-0.346-1.186-0.819-2.326-1.411-3.405l3.45-4.917c0.279-0.397,0.231-0.938-0.112-1.282 l-3.889-3.887c-0.347-0.346-0.893-0.391-1.291-0.104l-4.843,3.481c-1.089-0.602-2.239-1.08-3.432-1.427l-1.031-5.886 C28.607,2.35,28.192,2,27.706,2h-5.5c-0.49,0-0.908,0.355-0.987,0.839l-0.956,5.854c-1.2,0.345-2.352,0.818-3.437,1.412l-4.83-3.45 c-0.399-0.285-0.942-0.239-1.289,0.106L6.82,10.648c-0.343,0.343-0.391,0.883-0.112,1.28l3.399,4.863 c-0.605,1.095-1.087,2.254-1.438,3.46l-5.831,0.971c-0.482,0.08-0.836,0.498-0.836,0.986v5.5c0,0.485,0.348,0.9,0.825,0.985 l5.831,1.034c0.349,1.203,0.831,2.362,1.438,3.46l-3.441,4.813c-0.284,0.397-0.239,0.942,0.106,1.289l3.888,3.891 c0.343,0.343,0.884,0.391,1.281,0.112l4.87-3.411c1.093,0.601,2.248,1.078,3.445,1.424l0.976,5.861C21.3,47.647,21.717,48,22.206,48 h5.5c0.485,0,0.9-0.348,0.984-0.825l1.045-5.89c1.199-0.353,2.348-0.833,3.43-1.435l4.905,3.441 c0.398,0.281,0.938,0.232,1.282-0.111l3.888-3.891c0.346-0.347,0.391-0.894,0.104-1.292l-3.498-4.857 c0.593-1.08,1.064-2.222,1.407-3.408l5.918-1.039c0.479-0.084,0.827-0.5,0.827-0.985v-5.5C47.999,21.718,47.644,21.3,47.16,21.221z M25,32c-3.866,0-7-3.134-7-7c0-3.866,3.134-7,7-7s7,3.134,7,7C32,28.866,28.866,32,25,32z"
		/>
	</svg>
	<div class="lobby-sidebar {isOpen ? 'is-open' : 'is-closed'}">
		<button class="close-button" on:click={close}>
			<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 30 30" width="20px" height="20px">
				<path
					d="M 7 4 C 6.744125 4 6.4879687 4.0974687 6.2929688 4.2929688 L 4.2929688 6.2929688 C 3.9019687 6.6839688 3.9019687 7.3170313 4.2929688 7.7070312 L 11.585938 15 L 4.2929688 22.292969 C 3.9019687 22.683969 3.9019687 23.317031 4.2929688 23.707031 L 6.2929688 25.707031 C 6.6839688 26.098031 7.3170313 26.098031 7.7070312 25.707031 L 15 18.414062 L 22.292969 25.707031 C 22.682969 26.098031 23.317031 26.098031 23.707031 25.707031 L 25.707031 23.707031 C 26.098031 23.316031 26.098031 22.682969 25.707031 22.292969 L 18.414062 15 L 25.707031 7.7070312 C 26.098031 7.3170312 26.098031 6.6829688 25.707031 6.2929688 L 23.707031 4.2929688 C 23.316031 3.9019687 22.682969 3.9019687 22.292969 4.2929688 L 15 11.585938 L 7.7070312 4.2929688 C 7.5115312 4.0974687 7.255875 4 7 4 z"
				/>
			</svg>
		</button>
		<form action="" class="px-4 mt-4">
			<div class="mb-3">
				<label for="name">Name</label>
				<input type="text" id="name" required />
			</div>
			<div class="mb-3">
				<label for="visibility">Visibility</label>
				<select name="visibility" id="visibility">
					<option value="">Private</option>
					<option value="">Public</option>
				</select>
			</div>
			<div class="mb-3">
				<label for="link">Lobby Link</label>
				<div class="relative">
					<input
						type="text"
						id="link"
						value="https://skullking.ir/lobbies/d9db1493-6d11-474e-9ede-e7536f5fc362"
						disabled
					/>
					<div class="absolute top-2 right-1 w-10 right-0 bg-slate-900">
						<div class="w-full h-full flex justify-end items-center">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								width="20"
								height="20"
								class="mr-2 mt-1 cursor-pointer"
								data-name="Group 27733"
								viewBox="0 0 230.451 235.429"
								><path
									fill="#fff"
									d="M86.832 63.594h115.407a28.247 28.247 0 0 1 28.212 28.216v115.4a28.247 28.247 0 0 1-28.212 28.216H86.832a28.247 28.247 0 0 1-28.211-28.216V91.81a28.247 28.247 0 0 1 28.211-28.216Zm115.407 157.431a13.826 13.826 0 0 0 13.808-13.812v-115.4a13.826 13.826 0 0 0-13.808-13.812H86.832a13.826 13.826 0 0 0-13.807 13.813v115.4a13.826 13.826 0 0 0 13.807 13.813Z"
									data-name="Path 25"
								/><path
									fill="#fff"
									d="M48.077 0h89.5a34.289 34.289 0 0 1 34.255 34.251v3.072a7.2 7.2 0 1 1-14.4 0v-3.072a19.875 19.875 0 0 0-19.853-19.848h-89.5a33.716 33.716 0 0 0-33.674 33.674v89.506a19.873 19.873 0 0 0 19.848 19.848 7.2 7.2 0 0 1 0 14.4A34.291 34.291 0 0 1 .001 137.58V48.078A48.131 48.131 0 0 1 48.077 0Z"
									data-name="Path 26"
								/>
							</svg>
						</div>
					</div>
				</div>
			</div>

			<button
				class="w-full mt-3 text-slate-800 bg-lime-primary py-2 rounded outline-none focus:outline-none"
				type="button"
			>
				Update
			</button>
		</form>
	</div>
</div>

<style lang="scss">
	.lobby-sidebar {
		@apply absolute w-80 h-full bg-slate-900 right-0 top-0 rounded-l-3xl;
		&.is-closed {
			-webkit-animation-name: closeSidebar;
			animation-name: closeSidebar;
		}
		&.is-open {
			-webkit-animation-name: openSidebar;
			animation-name: openSidebar;
		}
		-webkit-animation-duration: 0.8s;
		animation-duration: 0.8s;
		animation-fill-mode: forwards;
	}

	@keyframes openSidebar {
		from {
			right: -320px;
		}

		to {
			right: 0;
		}
	}

	@keyframes closeSidebar {
		from {
			right: 0;
		}

		to {
			right: -320px;
		}
	}

	.close-button {
		@apply p-2 bg-gray-600 rounded-full mt-4 ml-4;
	}
</style>
