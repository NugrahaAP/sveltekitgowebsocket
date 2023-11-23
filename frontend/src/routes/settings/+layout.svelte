<script>
	import ArrowLeftBack from '$lib/svg/arrowLeftBack.svelte';
	import Person from '$lib/svg/person.svelte';
	import CreditCard from '$lib/svg/creditCard.svelte';
	import { page } from '$app/stores';
</script>

<div class="min-h-[93vh] flex flex-col">
	<a
		href="/"
		class=" p-2.5 border-b font-medium dark:border-surface-600 border-surface-200 dark:hover:text-white hover:text-primary-500"
	>
		<div class="flex flex-row items-center gap-2 lg:ms-52 xl:ms-72">
			<ArrowLeftBack />
			Back
		</div>
	</a>
	<div class="flex flex-row h-full grow px-2.5 max-sm:flex-col-reverse max-sm:p-0">
		<div
			class="border-e max-sm:border-t max-sm:border-e-0 dark:border-surface-600 border-surface-200 pe-2.5 py-2.5 font-medium lg:ms-52 xl:ms-72"
		>
			<div
				class="h-full flex flex-col max-sm:justify-center max-sm:flex-row max-sm:p-2.5 gap-2 pt-5"
			>
				<a
					href="/settings/account"
					class="flex flex-row items-center gap-2 p-5 max-sm:p-4 rounded dark:hover:text-white hover:text-primary-500 w-full hover:bg-surface-300/50 dark:hover:bg-surface-500/50 max-sm:w-fit {$page
						.route.id == '/settings/account'
						? 'bg-surface-300/50 dark:bg-surface-500/50'
						: ''}"
					style="text-wrap: nowrap;"
				>
					<Person />
					<p class="flex items-center gap-2 lg:block hidden">Account Setting</p>
				</a>
				<a
					href="/settings/billing"
					class="flex flex-row items-center gap-2 p-5 max-sm:p-4 rounded dark:hover:text-white hover:text-primary-500 w-full hover:bg-surface-300/50 dark:hover:bg-surface-500/50 max-sm:w-fit {$page
						.route.id == '/settings/billing'
						? 'bg-surface-300/50 dark:bg-surface-500/50'
						: ''}"
					style="text-wrap: nowrap;"
				>
					<CreditCard />
					<p class="flex items-center gap-2 lg:block hidden">Billing Setting</p>
				</a>
			</div>
		</div>
		<slot />
	</div>
</div>
