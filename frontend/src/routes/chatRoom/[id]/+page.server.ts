import { API_URL } from '$env/static/private';
import type { Actions, PageServerLoad } from './$types';
import {
	responseToastEnum,
	type ChatRoom,
	type ServerResponse,
	type responseToast
} from '$lib/types/myTypes';
import { fail } from '@sveltejs/kit';
import { ZodError, z } from 'zod';

export const load: PageServerLoad = async ({ cookies, params, fetch }) => {
	let session = cookies.get('session');
	let res: { status: number; response: ServerResponse } | undefined;
	try {
		const fetchChatRooms = await fetch(API_URL + '/chat_room?crid=' + params.id, {
			method: 'GET',
			headers: new Headers({ Authorization: session as string })
		});

		const jsonString = await fetchChatRooms.json();

		res = { status: fetchChatRooms.status, response: jsonString };
	} catch (err) {
		console.error(err);
	}

	return { chatRoom: res?.response.data.chatRoom as ChatRoom };
};

const messageSchema = z.object({
	messageBody: z.string(),
	messageType: z.string(),
	chatRoomId: z.string(),
	messageLink: z.string()
});

export const actions: Actions = {
	default: async ({ fetch, request, cookies }) => {
		console.log('hit');
		const session = cookies.get('session');
		const formData = await request.formData();
		const { messageBody, messageType, chatRoomId, messageLink } = Object.fromEntries(
			formData
		) as Record<string, string>;

		let data: responseToast;
		let res: { status: number; response: ServerResponse };

		try {
			messageSchema.parse({ messageBody, messageType, chatRoomId, messageLink });
		} catch (err) {
			console.error(err);
			if (err instanceof ZodError) {
				const error = err.errors.map((e) => {
					return e.message;
				});

				data = { error: true, message: error, type: responseToastEnum.warning };
				return fail(400, data);
			} else {
				data = { error: true, message: ['Internal server error'], type: responseToastEnum.error };
				return fail(500, data);
			}
		}

		try {
			const fd = new FormData();
			fd.append('input-messageBody', messageBody);
			fd.append('input-messageType', messageType);
			fd.append('input-messageLink', messageLink);
			fd.append('input-chatRoomId', chatRoomId);

			const fetchMessage = await fetch(API_URL + '/message', {
				method: 'POST',
				headers: new Headers({ Authorization: session as string }),
				body: fd
			});
			const jsonString = await fetchMessage.json();
			console.log(jsonString);
			res = { status: fetchMessage.status, response: jsonString };
		} catch (err) {
			console.error(err);
			data = { error: true, message: ['Something went wrong'], type: responseToastEnum.error };
			return fail(500, data);
		}

		if (res.status != 201) {
			data = { error: true, message: res.response.messages, type: responseToastEnum.error };
			return fail(res.status, data);
		}

		console.log(res);
		return res;
	}
};
