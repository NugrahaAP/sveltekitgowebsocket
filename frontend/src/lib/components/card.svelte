<script lang="ts">
	import { relativeTime } from '$lib/misc/misc.js';
	export let href: string;
	export let name: string;
	export let message: string;
	export let strDateTime: string | undefined;
	export let email: string;
	export let newChat = false;
</script>

<a
	{href}
	class="h-[15vh] w-full flex flex-row gap-2 flex border-b border-surface-200 dark:border-surface-700 px-5"
>
	<div class="rounded-full bg-primary-500 w-2 h-2 my-2.5" class:opacity-0={!newChat} />
	<div class="flex flex-col w-full h-full overflow-hidden">
		<h4 class="h4 font-bold pt-2.5">{name}</h4>
		<p class="text-surface-400 pb-1">
			<small>{email}</small>
		</p>
		<article
			class="text-surface-400 w-[90%]"
			style="white-space: wrap;
        overflow: hidden;
        text-overflow: ellipsis;"
		>
			{message}
		</article>

		<p class="text-surface-400 ms-auto mt-auto pb-2.5">
			<small>{strDateTime ? relativeTime(strDateTime) : ''}</small>
		</p>
	</div>
</a>
