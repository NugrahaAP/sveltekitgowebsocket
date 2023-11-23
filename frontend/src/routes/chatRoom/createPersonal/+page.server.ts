import { responseToastEnum, type responseToast } from '$lib/types/myTypes';
import { ZodError, z } from 'zod';
import type { Actions } from './$types';
import { fail } from '@sveltejs/kit';
import { API_URL } from '$env/static/private';

const checkEmailSchema = z.string().email({ message: 'Valid email address is required' });

export const actions: Actions = {
	createPersonal: async ({ request, cookies, fetch, locals }) => {
		let data: responseToast;
		let session = cookies.get('session');
		let fd = await request.formData();
		const email = fd.get('input-email');

		try {
			checkEmailSchema.parse(email);
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

		if (email == locals.email) {
			data = {
				error: true,
				message: ['Cannot create room with your self as participant'],
				type: responseToastEnum.warning
			};

			return fail(400, data);
		}

		let jsonRes: { isError: boolean; messages: string[]; data: any };
		let status: number;

		try {
			const fetchCheckEmail = await fetch(API_URL + '/checkEmail?email=' + email, {
				method: 'get',
				headers: new Headers({ Authorization: session as string })
			});
			jsonRes = await fetchCheckEmail.json();
			status = fetchCheckEmail.status;
		} catch (err) {
			console.error(err);
			data = { error: true, message: ['Internal server error'], type: responseToastEnum.error };
			return fail(500, data);
		}

		if (jsonRes.isError == true) {
			data = { error: true, message: jsonRes.messages, type: responseToastEnum.error };
			return fail(status, data);
		}

		// create new chatroom type personal
		let fetchCreatePersonalRoom: Response;
		try {
			const formData = new FormData();
			formData.set('input-receiverId', jsonRes.data.userId);

			fetchCreatePersonalRoom = await fetch(API_URL + '/chat_room', {
				method: 'POST',
				headers: new Headers({ Authorization: session as string }),
				body: formData
			});
		} catch (err) {
			console.error(err);
			data = { error: true, message: ['Something went wrong'], type: responseToastEnum.error };

			return fail(500, data);
		}

		console.log(fetchCreatePersonalRoom.status);

		if (fetchCreatePersonalRoom.status == 200 || fetchCreatePersonalRoom.status == 201) {
			const message = () => {
				if (fetchCreatePersonalRoom.status == 200) {
					return ['Chatroom already exist'];
				} else {
					return ['Chatroom created'];
				}
			};
			data = {
				error: false,
				message: message(),
				type: responseToastEnum.warning
			};

			return { ...data, chatRoom: (await fetchCreatePersonalRoom.json()).data.chatRoom };
		} else {
			data = { error: true, message: ['Something went wrong'], type: responseToastEnum.error };

			return fail(500, data);
		}
	}
};
