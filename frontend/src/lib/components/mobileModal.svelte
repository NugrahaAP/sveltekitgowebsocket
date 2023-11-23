<script lang="ts">
	import { fly } from 'svelte/transition';
	import X from '$lib/svg/x.svelte';
	import { LightSwitch, ProgressRadial } from '@skeletonlabs/skeleton';
	import { applyAction, enhance } from '$app/forms';
	import ArrowLeft from '$lib/svg/arrowLeft.svelte';
	import Gear from '$lib/svg/gear.svelte';
	import GridFill from '$lib/svg/grid-fill.svelte';

	export let userData:
		| {
				id: string | undefined;
				email: string | undefined;
				name: string | undefined;
		  }
		| undefined;
	export let trigger: boolean = false;
	let waitingResult = false;
	let logOutForm: HTMLFormElement;
</script>

{#if trigger}
	<div
		transition:fly={{ x: -500 }}
		class="z-10 absolute dark:bg-surface-800 dark:text-surface-300 bg-surface-200 w-3/4 h-full p-5 sm:hidden"
	>
		<div class="w-full flex flex-col">
			<button
				on:click={() => {
					trigger = false;
				}}
				class="ms-auto px-5 dark:hover:text-white hover:text-primary-500"><X /></button
			>
		</div>
		<nav
			class="dark:border-surface-600 border-surface-200 px-5 font-medium my-auto flex flex-col items-start justify-start h-full"
		>
			<ul class="flex flex-col gap-1 py-10 lg:ms-52 xl:ms-72 w-full h-full">
				<div class="bg-red-500 contents h-full w-full">
					<button
						on:click={() => {
							trigger = false;
						}}
						class="dark:hover:text-white hover:text-primary-500 dark:hover:bg-surface-500 hover:bg-surface-300 p-2.5 w-full rounded text-start"
					>
						<a href="#">Overview</a>
					</button>
					<button
						on:click={() => {
							trigger = false;
						}}
						class="dark:hover:text-white hover:text-primary-500 dark:hover:bg-surface-500 hover:bg-surface-300 p-2.5 w-full rounded text-start"
					>
						<a href="#">Customers</a>
					</button>
					<button
						on:click={() => {
							trigger = false;
						}}
						class="dark:hover:text-white hover:text-primary-500 dark:hover:bg-surface-500 hover:bg-surface-300 p-2.5 w-full rounded text-start"
					>
						<a href="#">Products</a>
					</button>
					<button
						class="dark:hover:text-white hover:text-primary-500 dark:hover:bg-surface-500 hover:bg-surface-300 p-2.5 w-full rounded text-start"
						on:click={() => {
							trigger = false;
						}}
					>
						<a href="/">Index</a>
					</button>
					<hr />
					{#if userData?.id}
						<button
							class="dark:hover:text-white hover:text-primary-500 dark:hover:bg-surface-500 hover:bg-surface-300 p-2.5 w-full rounded text-start"
							on:click={() => {
								trigger = false;
							}}
						>
							<a class="flex flex-row gap-2 items-center" href="/"><GridFill />Home</a>
						</button>
						<button
							class="dark:hover:text-white hover:text-primary-500 dark:hover:bg-surface-500 hover:bg-surface-300 p-2.5 w-full rounded text-start"
							on:click={() => {
								trigger = false;
							}}
						>
							<a href="/settings/account" class="flex flex-row gap-2 items-center"
								><Gear />Settings</a
							>
						</button>
						<hr />
						<form
							use:enhance={() => {
								waitingResult = true;
								return async ({ result }) => {
									await applyAction(result);
									waitingResult = false;
									trigger = false;
								};
							}}
							bind:this={logOutForm}
							action="/logout/?logout"
							method="POST"
							style="display: contents;"
						>
							<button
								on:click={() => {
									logOutForm.submit;
								}}
								disabled={waitingResult}
								class="dark:hover:text-white hover:text-primary-500 dark:hover:bg-surface-500 hover:bg-surface-300 p-2.5 w-full rounded text-start flex flex-row items-center justify-start gap-2 mt-auto"
							>
								{#if waitingResult == true}<ProgressRadial
										width="w-4"
										meter="stroke-surface-100 dark:stroke-surface-900"
									/>{/if}
								<ArrowLeft /> Log out
							</button>
						</form>
					{:else}
						<button
							class="dark:hover:text-white hover:text-primary-500 dark:hover:bg-surface-500 hover:bg-surface-300 p-2.5 w-full rounded text-start"
							on:click={() => {
								trigger = false;
							}}
						>
							<a href="/login">Login</a>
						</button>
					{/if}
					<div
						class="dark:hover:text-white hover:text-primary-500 p-2.5 w-full rounded self-end mt-auto"
					>
						<LightSwitch />
					</div>
				</div>
			</ul>
		</nav>
	</div>
{/if}
