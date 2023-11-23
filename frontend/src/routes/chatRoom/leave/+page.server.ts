import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { API_URL } from '$env/static/private';
import { responseToastEnum, type responseToast } from '$lib/types/myTypes';

export const load: PageServerLoad = async ({}) => {
	throw redirect(302, '/chatRoom');
};

export const actions: Actions = {
	default: async ({ cookies, request, fetch }) => {
		const fd = await request.formData();
		const chatRoomId = fd.get('input-chatRoomId') as string;

		let data: responseToast;

		let session = cookies.get('session');

		let fetchLeave;
		try {
			fetchLeave = await fetch(API_URL + '/chat_room?crid=' + chatRoomId, {
				method: 'DELETE',
				headers: new Headers({ Authorization: session as string })
			});
		} catch (err) {
			console.error(err);

			data = { error: true, message: ['Internal server error'], type: responseToastEnum.error };

			return fail(500, data);
		}

		if (fetchLeave.status == 200) {
			data = { error: false, message: ['Success'], type: responseToastEnum.primary };

			return data;
		} else {
			data = { error: true, message: ['Internal server error'], type: responseToastEnum.error };

			return fail(500, data);
		}
	}
};
