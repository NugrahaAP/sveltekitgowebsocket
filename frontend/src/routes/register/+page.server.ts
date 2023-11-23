import { ZodError, z } from 'zod';
import type { Actions } from './$types';
import { responseToastEnum, type responseToast, type ServerResponse } from '$lib/types/myTypes';
import { fail, redirect } from '@sveltejs/kit';
import { AUTH_URL } from '$env/static/private';

const registerSchema = z.object({
	email: z
		.string({
			required_error: 'Email is required',
			invalid_type_error: 'Email must be a string'
		})
		.email({ message: 'Email is required' }),
	name: z
		.string({
			required_error: 'Name is required',
			invalid_type_error: 'Name must be a string'
		})
		.min(1, { message: 'Name must be 1 or more characters long' }),
	password: z
		.string({
			required_error: 'Username is required',
			invalid_type_error: 'Username must be a string'
		})
		.min(8, { message: 'Password must be 8 or more characters long' })
});

export const actions: Actions = {
	default: async ({ request, fetch }) => {
		let data: responseToast;

		const formData = await request.formData();

		const { email, name, password, confirm_password } = Object.fromEntries(formData) as Record<
			string,
			string
		>;

		if (password != confirm_password) {
			data = { error: true, message: ["Password didn't match"], type: responseToastEnum.warning };
			return fail(400, data);
		}

		try {
			registerSchema.parse({ email, name, password });
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

		let res: { status: number; response: ServerResponse };

		try {
			const fd = new FormData();
			fd.append('input-email', email);
			fd.append('input-password', password);
			fd.append('input-name', name);

			const fetchRegis = await fetch(AUTH_URL + '/register', {
				method: 'POST',
				body: fd
			});
			const jsonString = await fetchRegis.json();

			res = { status: fetchRegis.status, response: jsonString };
		} catch (err) {
			console.error(err);
			data = { error: true, message: ['Something went wrong'], type: responseToastEnum.error };
			return fail(500, data);
		}

		if (res.status != 201) {
			data = { error: true, message: res.response.messages, type: responseToastEnum.error };
			return fail(res.status, data);
		}

		throw redirect(302, '/login');
	}
};
