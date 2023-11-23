import { API_URL } from '$env/static/private';
import type { LayoutServerLoad } from './$types';
import type { ChatRoom, GroupChatRoom, ServerResponse } from '$lib/types/myTypes';

export const load: LayoutServerLoad = async ({ cookies, params, fetch, locals }) => {
	let session = cookies.get('session');
	let resChatRoom: { status: number; response: ServerResponse } | undefined;
	let resGroupChatroom: { status: number; response: ServerResponse } | undefined;
	try {
		const fetchChatRooms = await fetch(API_URL + '/chat_room', {
			method: 'GET',
			headers: new Headers({ Authorization: session as string })
		});

		const jsonString = await fetchChatRooms.json();

		resChatRoom = { status: fetchChatRooms.status, response: jsonString };
	} catch (err) {
		console.error(err);
	}

	try {
		const fetchGroupChatRooms = await fetch(API_URL + '/group_chat_room', {
			method: 'GET',
			headers: new Headers({ Authorization: session as string })
		});

		const jsonString2 = await fetchGroupChatRooms.json();
		resGroupChatroom = { status: fetchGroupChatRooms.status, response: jsonString2 };
	} catch (err) {
		console.error(err);
	}

	console.log(resGroupChatroom);

	return {
		chatRooms: resChatRoom?.response.data.chatrooms as ChatRoom[],
		groupChatRooms: resGroupChatroom?.response.data.groupChatRooms as GroupChatRoom[],
		userData: { id: locals.userId, name: locals.name, email: locals.email }
	};
};
