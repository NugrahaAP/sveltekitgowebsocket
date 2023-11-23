import { responseToastEnum, type responseToast } from '$lib/types/myTypes';
import { ZodError, z } from 'zod';
import type { Actions, PageServerLoad } from './$types';
import { fail } from '@sveltejs/kit';
import { API_URL } from '$env/static/private';

export const load: PageServerLoad = async ({ cookies }) => {
	let data: responseToast;

	const session = cookies.get('session');

	let fetchUsers: Response;
	try {
		fetchUsers = await fetch(API_URL + '/user?listUser=true', {
			method: 'GET',
			headers: new Headers({ Authorization: session as string })
		});
	} catch (err) {
		console.error(err);
		data = { error: true, message: ['Something went wrong'], type: responseToastEnum.error };

		return fail(500, data);
	}

	return { users: (await fetchUsers.json()).data.users };
};

const userSchema = z.string().uuid({ message: 'Valid uuid is required' });

const createGroupSchema = z.object({
	groupName: z.string().min(1, { message: 'Group name must have atleast 1 characters' }),
	users: z.array(userSchema)
});

export const actions: Actions = {
	default: async ({ request, cookies, locals }) => {
		let data: responseToast;

		const fd = await request.formData();
		const groupName = fd.get('input-groupName');
		const users = fd.getAll('user');
		users.unshift(locals.userId);

		let session = cookies.get('session');

		console.log({ groupName, users });

		// validasi zod
		try {
			createGroupSchema.parse({ groupName, users });
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

		// fetch ke /group_chat_room

		let fetchCreateGroupRoom: Response;
		try {
			const formData = new FormData();
			formData.set('input-roomName', groupName as string);
			const usersJsonStr = `{"participantId":${JSON.stringify(users)}}`;
			console.log(usersJsonStr);
			formData.set('input-participant', usersJsonStr);

			fetchCreateGroupRoom = await fetch(API_URL + '/group_chat_room', {
				method: 'POST',
				headers: new Headers({ Authorization: session as string }),
				body: formData
			});
		} catch (err) {
			console.error(err);
			data = { error: true, message: ['Something went wrong'], type: responseToastEnum.error };

			return fail(500, data);
		}
		console.log(fetchCreateGroupRoom.status);

		if (fetchCreateGroupRoom.status == 201) {
			data = {
				error: false,
				message: ['Success'],
				type: responseToastEnum.primary
			};

			const jsonobj = await fetchCreateGroupRoom.json();
			console.log(jsonobj);

			return { ...data, groupChatRoom: jsonobj.data.groupChatRoom };
		} else {
			data = { error: true, message: ['Something went wrong'], type: responseToastEnum.error };

			return fail(500, data);
		}
	}
};
