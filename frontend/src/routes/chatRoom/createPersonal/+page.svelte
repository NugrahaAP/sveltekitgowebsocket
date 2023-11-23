<script lang="ts">
	import { applyAction, enhance } from '$app/forms';
	import { goto, invalidateAll } from '$app/navigation';
	import { toastData } from '$lib/store';
	import PersonAdd from '$lib/svg/personAdd.svelte';
	import type { ActionData } from './$types';

	export let form: ActionData;

	let waitingResult = false;

	$: toast = {
		message: form?.message,
		type: form?.type
	};

	$: $toastData = toast;
</script>

<main class="w-full min-h-[94vh] bg-surface-50 dark:bg-surface-900 flex">
	<div class="border border-surface-200 dark:border-surface-700 p-5 rounded w-1/2 m-auto">
		<h1 class="h3 font-bold">Start new conversation</h1>
		<div class="py-2.5 flex flex-col">
			<p class="text-surface-400 py-1">Email</p>
			<form
				use:enhance={() => {
					waitingResult = true;
					return async ({ result }) => {
						await applyAction(result);
						await invalidateAll();
						console.log(result);
						if (result.status == 200 || result.status == 201) {
							goto('/chatRoom/' + form?.chatRoom.id);
						}
						waitingResult = false;
					};
				}}
				method="POST"
				action="?/createPersonal"
				class="w-full border border-surface-200 dark:border-surface-700 rounded flex"
			>
				<input
					type="text"
					name="input-email"
					class="p-2 bg-transparent my-1 outline-none"
					placeholder="Ketik email teman mu.."
					disabled={waitingResult}
					required
				/>
				<button
					disabled={waitingResult}
					class="btn variant-filled-primary rounded-tr rounded-br rounded-tl-none rounded-bl-none ms-auto border-y border-e border-surface-200 dark:border-surface-700"
				>
					<PersonAdd />
				</button>
			</form>
		</div>
	</div>
</main>
