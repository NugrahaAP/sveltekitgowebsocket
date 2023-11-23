import { responseToastEnum, type responseToast, type ServerResponse } from '$lib/types/myTypes';
import { fail, type Actions, redirect } from '@sveltejs/kit';
import { ZodError, z } from 'zod';
import { AUTH_URL } from '$env/static/private';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ cookies }) => {
	const session = cookies.get('session');
	const refresh = cookies.get('refreshToken');

	if (session || refresh) {
		throw redirect(302, '/');
	}
};

const loginSchema = z.object({
	email: z
		.string({
			required_error: 'Email is required',
			invalid_type_error: 'Email must be a string'
		})
		.email({ message: 'Email is required' }),
	password: z
		.string({
			required_error: 'Username is required',
			invalid_type_error: 'Username must be a string'
		})
		.min(8, { message: 'Password must be 8 or more characters long' })
});

export const actions: Actions = {
	default: async ({ fetch, request, cookies, locals }) => {
		const formData = await request.formData();
		let data: responseToast;

		const { email, password } = Object.fromEntries(formData) as Record<string, string>;

		try {
			loginSchema.parse({ email, password });
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

			const fetchLogin = await fetch(AUTH_URL + '/login', {
				method: 'POST',
				body: fd
			});
			const jsonString = await fetchLogin.json();

			res = { status: fetchLogin.status, response: jsonString };
		} catch (err) {
			console.error(err);
			data = { error: true, message: ['Something went wrong'], type: responseToastEnum.error };
			return fail(500, data);
		}

		if (res.status != 200) {
			data = { error: true, message: res.response.messages, type: responseToastEnum.error };
			return fail(res.status, data);
		}

		const accessTokenExpirationDate = new Date();
		accessTokenExpirationDate.setTime(accessTokenExpirationDate.getTime() + 1 * 60 * 60 * 1000);

		const refreshTokenExpirationDate = new Date();
		refreshTokenExpirationDate.setTime(accessTokenExpirationDate.getTime() + 6 * 60 * 60 * 1000);

		cookies.set('session', res.response.data.accessToken, {
			expires: accessTokenExpirationDate,
			httpOnly: true,
			path: '/'
		});
		cookies.set('refreshToken', res.response.data.refreshToken, {
			expires: refreshTokenExpirationDate,
			httpOnly: true,
			path: '/'
		});

		locals.userId = res.response.data.user.id;
		locals.email = res.response.data.user.email;
		locals.name = res.response.data.user.name;

		throw redirect(302, '/');
	}
};
