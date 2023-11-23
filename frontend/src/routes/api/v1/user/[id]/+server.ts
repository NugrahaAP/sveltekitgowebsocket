import { json, redirect } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { responseToastEnum, type responseToast, type ServerResponse } from '$lib/types/myTypes';
import { API_URL } from '$env/static/private';

export const GET: RequestHandler = async ({ fetch, cookies, params }) => {
	let data: responseToast;

	let accessToken = cookies.get('session');
	// if (!accessToken) {
	// 	cookies.delete('session');
	// 	cookies.delete('refreshToken');

	// 	throw redirect(302, '/login');
	// }
	console.log('params', params.id);
	let res: { status: number; response: ServerResponse };
	try {
		const fetchUser = await fetch(API_URL + `/user?userId=${params.id}`, {
			method: 'GET',
			headers: new Headers({ Authorization: accessToken as string })
		});
		const jsonString = await fetchUser.json();

		res = { status: fetchUser.status, response: jsonString };
	} catch (err) {
		console.error(err);
		data = { error: true, message: ['Something went wrong'], type: responseToastEnum.error };

		return json(data, { status: 500 });
	}

	return json(res.response);
};
