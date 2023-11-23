<script lang="ts">
	import type { Message, User } from '$lib/types/myTypes';
	import type { LayoutData } from './$types';
	import { relativeTime } from '$lib/misc/misc.js';
	import EnvelopePlus from '$lib/svg/envelopePlus.svelte';
	import Menu from '$lib/components/menu.svelte';
	import Person from '$lib/svg/person.svelte';
	import People from '$lib/svg/people.svelte';
	import { page } from '$app/stores';
	import { groupName } from '$lib/store';

	export let data: LayoutData;

	function getOtherUser(u: User[]) {
		let name = '';
		u.forEach((user) => {
			if (user.id != data.userData.id) {
				name = user.name;
			}
		});

		return name;
	}

	function getLastMessage(m: Message) {
		try {
			if (m.sender.id == data.userData.id) {
				return 'You: ' + m.messageBody;
			} else {
				return m.messageBody;
			}
		} catch (error) {
			return '';
		}
	}

	$: console.log(data);

	let layoutChatRoomState = { showCreateChatRoom: false };
</script>

{#if !String($page.route.id).startsWith('/chatRoom/create')}
	<div class="h-[94vh] w-screen 2xl:w-3/4 2xl:mx-auto flex flex-row">
		<div class="overflow-y-scroll w-1/2 border-x border-surface-200 dark:border-surface-700">
			<div
				class="top-0 bg-surface-50 dark:bg-surface-900 flex p-5 border-b border-surface-200 dark:border-surface-700"
			>
				<h4 class="h4 font-bold">Message</h4>
				<Menu position="" bind:show={layoutChatRoomState.showCreateChatRoom}>
					<button slot="menu" class="ms-auto btn btn-sm variant-filled-primary rounded">
						<EnvelopePlus />
					</button>
					<div slot="options" class="flex flex-col">
						<a
							on:click={() => {
								layoutChatRoomState.showCreateChatRoom = false;
							}}
							class="flex flex-row items-center gap-2 justify-between p-5 font-bold hover:bg-surface-100 dark:hover:bg-surface-800"
							href="/chatRoom/createPersonal">Personal <Person /></a
						>
						<a
							on:click={() => {
								layoutChatRoomState.showCreateChatRoom = false;
							}}
							class="flex flex-row text-primary-500 items-center gap-2 justify-between p-5 font-bold hover:bg-surface-100 dark:hover:bg-surface-800"
							href="/chatRoom/createGroup">Group<People /></a
						>
					</div>
				</Menu>
			</div>

			{#each data.groupChatRooms as group}
				<a
					on:click={() => {
						$groupName = group.roomName;
					}}
					href="/chatRoom/{group.chatRoom.id}?groupName={group.roomName}"
					class="border-b border-surface-200 dark:border-surface-700 p-5 flex flex-col"
				>
					<h4 class="h4 font-bold">
						{group.roomName}
					</h4>
					<article class="text-surface-400">
						{getLastMessage(group.chatRoom.messages[0]) || ''}
					</article>
					<p class="ms-auto text-surface-400">
						<small>{relativeTime(group.chatRoom.messages[0]?.createdAt)}</small>
					</p>
				</a>
			{/each}

			{#each data.chatRooms as room}
				<a
					href="/chatRoom/{room.id}"
					class="border-b border-surface-200 dark:border-surface-700 p-5 flex flex-col"
				>
					<h4 class="h4 font-bold">
						{getOtherUser(room.participant) || ''}
					</h4>
					<article class="text-surface-400">
						{getLastMessage(room.messages[0]) || ''}
					</article>
					<p class="ms-auto text-surface-400">
						<small>{relativeTime(room.messages[0]?.createdAt)}</small>
					</p>
				</a>
			{/each}
		</div>

		<div class="w-full border-e dark:border-surface-700 border-surface-200">
			<slot />
		</div>
	</div>
{:else}
	<div class="w-full">
		<slot />
	</div>
{/if}
