<script lang="ts">
	import Menu from '$lib/components/menu.svelte';
	import ThreeDot from '$lib/svg/threeDot.svelte';
	import type { PageData } from './$types';
	import PaperPlane from '$lib/svg/paperPlane.svelte';
	import { messageTypeEnum, type Message, type User, type WsResponse } from '$lib/types/myTypes';
	import CenterModal from '$lib/components/centerModal.svelte';
	import { getOtherUserId } from '$lib/misc/misc';
	import UserChatBubble from '$lib/components/userChatBubble.svelte';
	import OtherUserChatBubble from '$lib/components/otherUserChatBubble.svelte';
	import { onDestroy, onMount } from 'svelte';
	import { navigating, page } from '$app/stores';
	import { applyAction, enhance } from '$app/forms';
	import { newWebsocketConn } from '$lib/websocket/websocket';
	import { currentRoomId, groupName } from '$lib/store';
	import { browser } from '$app/environment';
	import CenterModalGroup from '$lib/components/centerModalGroup.svelte';
	import { goto, invalidateAll } from '$app/navigation';

	export let data: PageData;
	let inputMessage: string;
	let chatContainer: HTMLDivElement;
	let waitingResult = false;
	let wsConn: WebSocket | undefined;
	let wsPayload: WsResponse;
	let isWaiting = false;

	$: if ($groupName == '') {
		$groupName = $page.url.searchParams.get('groupName') as string;
	}

	function setChatTitle(idata: any) {
		if (data.chatRoom.chatRoomType == 'group') {
			return $groupName;
		}
		if (data.chatRoom.chatRoomType == 'personal') {
			let title;
			data.chatRoom.participant.forEach((u: User) => {
				if (u.id != data.userData.id) {
					title = u.name;
				}
			});
			return title;
		}
	}

	onMount(() => {
		setTimeout(() => {
			chatContainer.scrollTo({ top: chatContainer.scrollHeight, behavior: 'smooth' });
		}, 100);
	});

	onDestroy(() => {
		if (wsConn) {
			wsConn.close();
			wsConn = undefined;
		}
	});

	$: if (($navigating || browser) && $page.params.id) {
		// check ws
		// if exists and currentroomid != "" && currentroomid != data.chatRoomid
		if (wsConn && $currentRoomId.roomId != '' && $currentRoomId.roomId != data.chatRoom.id) {
			// close ws conn
			console.log(
				'websocket connection already opened to: ',
				$currentRoomId.roomId,
				'\n',
				'closing...'
			);
			wsConn.close();
			wsConn = undefined;
		}

		if (!wsConn) {
			// connect to data.chatrom.id
			console.log('Opening new connection to: ', data.chatRoom.id);
			wsConn = newWebsocketConn(data.chatRoom.id);
			$currentRoomId.roomId = data.chatRoom.id;
		}

		// onMessage = ()=>{}
		wsConn.onmessage = (e) => {
			console.log(e.data);
			const wsRes: WsResponse = JSON.parse(e.data);
			const newWsMsg: Message = {
				id: 'ws chat',
				chatRoomId: data.chatRoom.id,
				createdAt: wsRes.createdAt,
				deletedAt: null,
				messageBody: wsRes.message,
				messageLink: '',
				messageType: messageTypeEnum.message,
				sender: {
					id: wsRes.senderId,
					createdAt: new Date().toString(),
					deletedAt: null,
					email: 'ws email',
					name: wsRes.name,
					updatedAt: new Date().toString()
				},

				updatedAt: new Date().toString(),
				userId: wsRes.senderId
			};
			data.chatRoom.messages.push(newWsMsg);
			data.chatRoom.messages = data.chatRoom.messages;
			setTimeout(() => {
				chatContainer.scrollTo({ top: chatContainer.scrollHeight, behavior: 'smooth' });
			}, 100);
		};
		console.log('done creating msg');
	}

	let otherUserId = getOtherUserId(data.userData.id, data.chatRoom.participant);
	let chatRoomState = { showDotMenu: false, showDetailPesonal: false };
</script>

{#if data.chatRoom.chatRoomType == 'personal'}
	<CenterModal userId={String(otherUserId)} bind:show={chatRoomState.showDetailPesonal} />
{/if}

<div class="w-full h-full flex flex-col">
	<div class="border-b border-surface-200 dark:border-surface-700 p-5 flex flex-row justify-center">
		<h4 class="h4 font-bold">{setChatTitle(data)}</h4>
		<Menu bind:show={chatRoomState.showDotMenu}>
			<button
				on:click={() => {
					chatRoomState.showDotMenu = !chatRoomState.showDotMenu;
				}}
				class="p-1"
				slot="menu"><ThreeDot /></button
			>
			<div class="flex flex-col" slot="options">
				{#if data.chatRoom.chatRoomType == 'personal'}
					<button
						on:click={() => {
							chatRoomState.showDotMenu = !chatRoomState.showDotMenu;

							chatRoomState.showDetailPesonal = true;
						}}
						class="btn font-bold">Detail</button
					>
				{/if}
				{#if data.chatRoom.chatRoomType == 'group'}
					<form
						use:enhance={() => {
							isWaiting = true;
							return async ({ result }) => {
								await invalidateAll();
								if (result.status == 200) {
									goto('/chatRoom');
								}

								isWaiting = false;
								await applyAction(result);
							};
						}}
						action="/chatRoom/leave/?leave"
						method="POST"
						style="display: contents;"
					>
						<button
							on:click={() => {
								chatRoomState.showDotMenu = !chatRoomState.showDotMenu;
							}}
							class="btn font-bold text-primary-500">Leave</button
						>
						<input type="hidden" name="input-chatRoomId" value={data.chatRoom.id} />
					</form>
				{/if}
			</div>
		</Menu>
	</div>

	<!-- chat container -->
	<div bind:this={chatContainer} class="overflow-x-hidden px-10 overflow-y-scroll">
		{#each data.chatRoom.messages as msg}
			{#if msg.userId == data.userData.id}
				<UserChatBubble message={msg.messageBody} strDateTime={msg.createdAt} />
			{:else}
				<OtherUserChatBubble
					name={msg.sender.name}
					message={msg.messageBody}
					strDateTime={msg.createdAt}
				/>
			{/if}
		{/each}
	</div>

	<form
		use:enhance={() => {
			waitingResult = true;
			return async ({ result }) => {
				wsPayload = {
					action: 'show',
					createdAt: new Date().toString(),
					message: inputMessage,
					name: data.userData.name,
					senderId: data.userData.id
				};
				wsConn?.send(JSON.stringify(wsPayload));

				inputMessage = '';
				waitingResult = false;
				await applyAction(result);
			};
		}}
		action="?chatroom/{data.chatRoom.id}"
		method="POST"
		class="mt-auto flex flex-row border-t border-surface-200 dark:border-surface-700 dark:bg-surface-800"
	>
		<input
			class="w-full p-2 bg-transparent m-2 border rounded border-surface-200 dark:border-surface-700 dark:bg-surface-700"
			style="outline: none;"
			bind:value={inputMessage}
			type="text"
			name="messageBody"
			minlength="1"
			required
			placeholder="Ketik sesuatu..."
		/>
		<input type="hidden" name="messageType" value="message" />
		<input type="hidden" name="chatRoomId" value={data.chatRoom.id} />
		<input type="hidden" name="messageLink" value="" />
		<button class="p-2 m-2 w-16 btn variant-filled-primary rounded">
			<PaperPlane />
		</button>
	</form>
</div>
