<script lang="ts">
	export let message: string;
	export let strDateTime: string;
	import { relativeTime } from '$lib/misc/misc.js';
</script>

<div
	class="w-fit text-white bg-secondary-500 dark:bg-secondary-700 flex flex-col max-w-[40%] p-5 m-2.5
	my-5 ms-auto rounded-tl-xl rounded-bl-xl rounded-tr-xl"
>
	<p style="text-wrap: wrap;">{message}</p>

	<p class="py-1 ms-auto text-surface-100"><small>{relativeTime(strDateTime) || ''}</small></p>
</div>
