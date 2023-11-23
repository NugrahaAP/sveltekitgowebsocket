<script lang="ts">
	import { applyAction, enhance } from '$app/forms';
	import { goto } from '$app/navigation';
	import { toastData } from '$lib/store';
	import { responseToastEnum } from '$lib/types/myTypes';
	import { ProgressRadial } from '@skeletonlabs/skeleton';
	import type { ActionData } from './$types';

	export let form: ActionData;
	let waitingResult = false;

	$: toast = {
		message: form?.message,
		type: form?.type
	};

	$: $toastData = toast;
</script>

<div class="m-auto sm:w-1/2 md:w-1/3 lg:w-1/4 xl:w-1/5 my-auto">
	<h4 class="h3 font-medium py-5 text-center">Login to your account</h4>

	<form
		use:enhance={() => {
			waitingResult = true;
			return async ({ result }) => {
				if (result.type == 'redirect') {
					toast = { message: ['Login success'], type: responseToastEnum.primary };
					$toastData = toast;
					goto(result.location);
				} else {
					await applyAction(result);
				}
				waitingResult = false;
			};
		}}
		action="?login"
		method="POST"
		class="flex flex-col w-full gap-2"
	>
		<input
			class="bg-transparent dark:border-surface-600 border-surface-300 border rounded p-1.5 ps-2"
			type="text"
			name="email"
			id="email"
			placeholder="email"
			disabled={waitingResult}
			required
		/>
		<input
			class="bg-transparent dark:border-surface-600 border-surface-300 border rounded p-1.5 ps-2"
			type="password"
			name="password"
			id="password"
			placeholder="Password"
			disabled={waitingResult}
			required
		/>
		<button
			disabled={waitingResult}
			type="submit"
			class="btn variant-filled-primary rounded font-medium gap-2"
			>Login {#if waitingResult}
				<ProgressRadial
					width="w-4"
					meter="stroke-surface-100 dark:stroke-surface-900"
				/>{/if}</button
		>
	</form>
	<div class="flex flex-row justify-between py-5 w-full">
		<p class="text-surface-400">
			<small>
				<a class="text-primary-500" href="/register">Register</a>
			</small>
		</p>
	</div>
</div>
