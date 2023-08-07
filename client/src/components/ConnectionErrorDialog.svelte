<script lang="ts">
	function reload() {
		window.location.reload();
	}
</script>

<div class="fluid-container">
	<div class="dialog-box">
		<p class="text-xl text-red-400 font-semibold uppercase mb-2">Connection Error</p>
		<p class="text-gray-100">
			Connection to the server timed out. Please reload the page by clicking on the designated
			button or using the appropriate keyboard shortcut.
		</p>
		<button type="button" class="btn light-gradient mt-4" on:click={reload}>
			<img
				class="mr-2"
				width="20"
				alt="Reload page icon"
				src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABgAAAAYCAYAAADgdz34AAAACXBIWXMAAAsTAAALEwEAmpwYAAABXklEQVR4nN2VTytEURjGf1yZQtn5EGRWltZKQonYzIbFWBmKnWRhQ7P1CZDY+gY2VlZIM9NoLGwVaWak0NGjTjPn/pl7uwueOnXrfd7nOe97z3kP/HdkgBxwDjwADeAVKAFHwDzQE1d8CXgCvkKWMZu28naASpCwBxxaAtdAHhgBhoABYAzYAKrifAJ7wK6V54tf8TqwAnSFbGYL+HBU5sSiJT4etIsWHEcxyFg9X+5AfNvn37QhZ/U8qC2tqDjEL9tYwJmCq6SEkgyG0zJ4lsFgWgZ1GfQl0PCAK+DCFbyRQTaBQVYaZVfwVMFCAoM1aZy4grMKmkq6Y7bnThpzLoKZio8irMcwKCi3FjRhJzW4GvqOigngXQZTYeSiiE1VYkr3gydOUzkHUXe0b137W2BTJ8TckX5gVML31rgudjhiflpUjfDgmFk0Q0z0KtkcO1PJC/BmPZkLSZ7Mv4FvtId4idfl6NcAAAAASUVORK5CYII="
			/>
			<span>Reload</span>
		</button>
	</div>
</div>

<style lang="scss">
	.dialog-box {
		@apply p-5 max-w-md rounded-md bg-slate-700 text-lg;
	}
</style>
