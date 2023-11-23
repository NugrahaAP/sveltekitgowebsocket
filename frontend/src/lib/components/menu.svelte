<script lang="ts">
	import { slide } from 'svelte/transition';

	export let show: boolean = false;
	export let position: string = 'right-5 2xl:right-[20.5em]';
</script>

<div class="ms-auto">
	<button
		class="rounded dark:hover:bg-surface-700 hover:bg-surface-300 flex justify-center items-center gap-2"
		class:bg-surface-700={show}
		on:click={() => {
			show = true;
		}}
		on:keyup={(e) => {
			if (e.key == 'Escape') {
				show = false;
			}
		}}
	>
		<slot name="menu" />
	</button>

	{#if show == true}
		<div
			transition:slide
			class="absolute bg-surface-50 dark:bg-surface-900 border border-surface-300 dark:border-surface-700 flex flex-col w-1/6 2xl:w-[10%] max-h-[40%] rounded overflow-x-hidden mt-2.5 z-10 {position}"
		>
			<slot name="options" />
		</div>
	{/if}
</div>

{#if show == true}
	<!-- svelte-ignore a11y-click-events-have-key-events -->
	<!-- svelte-ignore a11y-no-static-element-interactions -->
	<div
		class="bg-transparent w-screen h-screen absolute top-0 left-0"
		on:click={() => {
			show = false;
		}}
	>
		<br />
	</div>
{/if}
