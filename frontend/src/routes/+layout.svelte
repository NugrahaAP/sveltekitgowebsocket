<script lang="ts">
	import '../app.postcss';
	import { Toast, getToastStore } from '@skeletonlabs/skeleton';
	import { ProgressRadial, type ToastSettings } from '@skeletonlabs/skeleton';
	import { initializeStores } from '@skeletonlabs/skeleton';
	import { toastData } from '$lib/store';
	import { responseToastEnum } from '$lib/types/myTypes';
	import Gear from '$lib/svg/gear.svelte';
	import ArrowLeft from '$lib/svg/arrowLeft.svelte';
	import { scale } from 'svelte/transition';
	import { applyAction, enhance } from '$app/forms';
	import { page } from '$app/stores';
	import { LightSwitch } from '@skeletonlabs/skeleton';
	import List from '$lib/svg/list.svelte';
	import MobileModal from '$lib/components/mobileModal.svelte';
	import Menu from '$lib/components/menu.svelte';

	initializeStores();

	const toastStore = getToastStore();
	const t: ToastSettings = {
		message: 'Custom Toast',
		hoverable: true
	};
	let toastType: responseToastEnum | undefined = responseToastEnum.error;
	let showUserModal = false;
	let logoutForm: HTMLFormElement;
	let waitingResult = false;
	let showMobileModal = false;

	$: if ($toastData.message) {
		$toastData.message.forEach((data) => {
			toastStore.trigger({
				message: data,
				hoverable: true,
				background: `variant-filled-${$toastData.type ? $toastData.type : 'error'}`
			});
		});

		toastType = $toastData.type;
	}
</script>

<Toast position="br" background="variant-filled-warning" />

{#if showUserModal}
	<div
		transition:scale={{
			duration: 100
		}}
		class="absolute dark:bg-surface-900 2xl:w-1/12 lg:w-1/6 w-1/5 flex flex-col gap-4 p-5 border dark:border-surface-600 rounded shadow-md justify-center top-[3.5em] right-10 2xl:right-[33.5em] dark:text-surface-300 font-medium"
	>
		<a
			on:click={() => {
				showUserModal = false;
			}}
			class="flex flex-row items-center justify-start gap-2 mx-auto dark:hover:text-white hover:text-primary-500"
			href="/settings/account"
		>
			<Gear /> Settings</a
		>
		<hr />
	</div>
{/if}

<MobileModal
	userData={$page?.data?.userData ? $page?.data?.userData : undefined}
	bind:trigger={showMobileModal}
/>

<div class="flex flex-col min-w-screen min-h-screen dark:text-surface-300">
	{#if $page.data?.userData?.id}
		<nav class="dark:border-surface-600 border-surface-200 border-b px-5 font-medium w-full">
			<ul class="flex flex-row gap-4 py-4 w-3/4 mx-auto max-lg:w-full">
				<div class="bg-red-500 contents max-sm:hidden">
					<li class="dark:hover:text-white hover:text-primary-500"><a href="/">Index</a></li>
					<button
						class="flex flex-row items-center justify-start gap-2 dark:hover:text-white hover:text-primary-500"
						on:click={() => {
							toastStore.trigger(t);
						}}>trigger</button
					>

					<div class="flex flex-row items-center gap-2 ms-auto">
						<Menu>
							<p slot="menu">{$page.data.userData.name}</p>
							<svelte:fragment slot="options">
								<form
									use:enhance={() => {
										waitingResult = true;
										return async ({ result }) => {
											await applyAction(result);
											waitingResult = false;
											showUserModal = false;
										};
									}}
									bind:this={logoutForm}
									action="/logout/?logout"
									method="POST"
									style="display: contents;"
								>
									<button
										on:click={() => {
											logoutForm.submit;
										}}
										disabled={waitingResult}
										class="flex btn text-primary-500 flex-row items-center justify-start gap-2 mx-auto"
									>
										{#if waitingResult == true}<ProgressRadial
												width="w-4"
												meter="stroke-surface-100 dark:stroke-surface-900"
											/>{/if}
										<ArrowLeft /> Log out
									</button>
								</form>
							</svelte:fragment>
						</Menu>
						<LightSwitch />
					</div>
				</div>
				<div class="max-sm:flex hidden flex-row items-center gap-2 w-full">
					<button
						on:click={() => {
							showMobileModal = true;
						}}
						class="dark:hover:text-white hover:text-primary-500 my-auto"
					>
						<List />
					</button>
					<a class="my-auto mx-auto dark:hover:text-white hover:text-primary-500" href="/"
						>App name</a
					>
				</div>
			</ul>
		</nav>
	{:else}
		<nav class="dark:border-surface-600 border-surface-200 border-b px-5 font-medium">
			<ul class="flex flex-row gap-4 py-4 px-5 font-medium lg:ms-52 xl:ms-72">
				<div class="bg-red-500 contents max-sm:hidden">
					<button
						class="flex flex-row items-center justify-start gap-2 dark:hover:text-white hover:text-primary-500"
						on:click={() => {
							toastStore.trigger(t);
						}}>trigger</button
					>
					<li
						class="flex flex-row items-center justify-start gap-2 ms-auto dark:hover:text-white hover:text-primary-500"
					>
						<a href="/login">login</a>
					</li>
					<LightSwitch />
				</div>
				<div class="max-sm:flex hidden flex-row items-center gap-2 w-full">
					<button
						on:click={() => {
							showMobileModal = !showMobileModal;
						}}
						class="dark:hover:text-white hover:text-primary-500 my-auto"
					>
						<List />
					</button>
					<a class="my-auto mx-auto dark:hover:text-white hover:text-primary-500" href="/"
						>App name</a
					>
				</div>
			</ul>
		</nav>
	{/if}
	<slot />
</div>
