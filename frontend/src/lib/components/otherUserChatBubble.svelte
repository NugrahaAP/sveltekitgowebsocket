<script lang="ts">
	import { relativeTime } from '$lib/misc/misc';

	export let name: string;
	export let message: string;
	export let strDateTime: string;
</script>

<div
	class="w-fit bg-surface-200 dark:bg-surface-800 flex flex-col w-[40%] p-5 m-2.5 me-auto my-5 rounded-tl-xl rounded-br-xl rounded-tr-xl"
>
	<h4 class="h4 font-bold text-primary-500">{name}</h4>
	<p>{message}</p>

	<p class="py-1 me-auto"><small>{relativeTime(strDateTime) || ''}</small></p>
</div>
