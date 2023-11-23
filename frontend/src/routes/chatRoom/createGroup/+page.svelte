<script lang="ts">
	import { applyAction, enhance } from '$app/forms';
	import { goto, invalidateAll } from '$app/navigation';
	import { toastData } from '$lib/store';
	import People from '$lib/svg/people.svelte';
	import type { ActionData, PageData } from './$types';

	export let form: ActionData;
	export let data: PageData;
	let waitingResult = false;

	$: toast = {
		message: form?.message,
		type: form?.type
	};

	$: $toastData = toast;

	$: console.log(form);
</script>

<main class="w-full min-h-[94vh] bg-surface-50 dark:bg-surface-900 flex">
	<div class="border border-surface-200 dark:border-surface-700 p-5 rounded w-1/2 m-auto">
		<h1 class="h3 font-bold">Start new group chat</h1>
		<form
			use:enhance={() => {
				waitingResult = true;
				return async ({ result }) => {
					console.log(result.status);
					await applyAction(result);
					await invalidateAll();

					if (result.status == 200) {
						console.log('redirect');
						goto('/chatRoom/' + form?.groupChatRoom.chatRoom.id);
					}
					waitingResult = false;
				};
			}}
			method="POST"
			action="?createGroup"
			class="py-2.5 flex flex-col"
		>
			<p class="text-surface-400 py-1 font-bold">Group name</p>
			<input
				type="text"
				name="input-groupName"
				class="p-2 bg-transparent my-1 border border-surface-200 dark:border-surface-700 rounded outline-none"
				placeholder="Circle kita.."
				disabled={waitingResult}
				required
			/>
			<div class="py-5">
				<p class="text-surface-400 py-2.5 font-bold">Participant</p>
				<div
					class="border max-h-36 overflow-y-scroll border-surface-200 dark:border-surface-700 rounded"
				>
					{#each data.users as user, index (user.id)}
						<div class="flex w-full border-b border-surface-200 dark:border-surface-700">
							<input class="ms-5 my-auto checkbox" type="checkbox" name={'user'} value={user.id} />
							<div class="p-5 w-full">
								<h4 class="h4 font-bold">{user.name}</h4>
								<p class="text-surface-400">
									<small>{user.email}</small>
								</p>
							</div>
						</div>
					{/each}
				</div>
			</div>
			<button
				type="submit"
				disabled={waitingResult}
				class="mt-12 btn variant-filled-primary rounded ms-auto"
			>
				<People />
			</button>
		</form>
	</div>
</main>
