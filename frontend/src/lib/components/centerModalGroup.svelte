<script lang="ts">
	import { browser } from '$app/environment';
	import X from '$lib/svg/x.svelte';
	import type { User } from '$lib/types/myTypes';
	import { ProgressRadial } from '@skeletonlabs/skeleton';
	import { onDestroy, onMount } from 'svelte';
	import { fade, scale } from 'svelte/transition';

	type ComponentResponse = {
		isError: boolean;
		messages: string[];
		data: User;
	};
	export let show = false;
	export let userId: string;
	let isWaiting = true;
	let userData: ComponentResponse | undefined;
	async function fetchUserDetail(userId: string) {
		isWaiting = true;

		//fetch user detail
		const userDetail = await fetch(`http://localhost:5173/api/v1/user/${userId}`);
		const jsonString = await userDetail.json();
		isWaiting = false;
		userData = jsonString;
	}

	$: if (browser && show) {
		console.log('show');
		fetchUserDetail(userId);
	}
</script>

{#if show}
	<!-- svelte-ignore a11y-click-events-have-key-events -->
	<!-- svelte-ignore a11y-no-static-element-interactions -->

	<div
		on:click={() => {
			show = !show;
		}}
		transition:fade={{ duration: 250 }}
		class="absolute bg-black/50 w-full h-full top-0 left-0"
	>
		{#if isWaiting}
			<div class="flex h-full">
				<div class="m-auto">
					<ProgressRadial />
				</div>
			</div>
		{:else}
			<div transition:scale class="flex h-full">
				<div
					class="relative w-1/2 bg-surface-50 dark:bg-surface-900 m-auto rounded p-5 border border-surface-200 dark:border-surface-700"
				>
					<h4 class="text-2xl font-bold">User detail</h4>
					<div class="p-2.5 my-2.5 rounded border border-surface-200 dark:border-surface-700">
						<p class="text-surface-400">
							<small>Name</small>
						</p>
						<h4 class="h4 font-bold">{userData?.data.name}</h4>
					</div>
					<div class="p-2.5 my-2.5 rounded border border-surface-200 dark:border-surface-700">
						<p class="text-surface-400">
							<small>Email</small>
						</p>
						<h4 class="h4 font-bold">{userData?.data.email}</h4>
					</div>

					<p class="text-surface-400">
						<small>Joined At: {new Date(String(userData?.data.createdAt)).toDateString()} </small>
					</p>
					<button
						class="w-5 h-5 flex items-center justify-center absolute -top-0 -right-0 variant-filled-primary rounded-tr rounded-bl p-1"
					>
						<X />
					</button>
				</div>
			</div>
		{/if}
	</div>
{/if}
