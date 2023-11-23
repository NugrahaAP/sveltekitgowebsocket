import type { ServerResponse } from '$lib/types/myTypes';
import { redirect, type Handle } from '@sveltejs/kit';
import { API_URL, AUTH_URL } from '$env/static/private';

export const handle: Handle = async ({ event, resolve }) => {
	if (event.route.id == '/login' || event.route.id == '/register' || event.route.id == '/logout') {
		return resolve(event);
	}

	let session = event.cookies.get('session');
	let refreshToken = event.cookies.get('refreshToken');

	// throw ke /login jika tidak ada session dan refreshToken
	if (!session || !refreshToken) {
		event.cookies.delete('session');
		event.cookies.delete('refreshToken');
		throw redirect(302, '/login');
	}

	// verify access token
	try {
		// fetch to /rahasia, if not 200 then refresh jwt
		const fetchRahasia = await fetch(API_URL + '/rahasia', {
			method: 'GET',
			headers: new Headers({ Authorization: session as string })
		});

		if (fetchRahasia.status != 200) {
			const fd = new FormData();
			fd.set('refresh-token', refreshToken as string);

			// try refreshing jwt
			const fetchRefreshJWT = await fetch(AUTH_URL + '/refresh', {
				method: 'POST',
				body: fd
			});

			const jsonString: ServerResponse = await fetchRefreshJWT.json();

			// if tidak 201 throw new error
			if (fetchRefreshJWT.status != 201) {
				throw new Error('Invalid refresh token');
			}

			const newJWT = {
				accessToken: jsonString.data.accessToken,
				refreshToken: jsonString.data.refreshToken
			};

			const accessTokenExpirationDate = new Date();
			accessTokenExpirationDate.setTime(accessTokenExpirationDate.getTime() + 1 * 60 * 60 * 1000);

			const refreshTokenExpirationDate = new Date();
			refreshTokenExpirationDate.setTime(accessTokenExpirationDate.getTime() + 6 * 60 * 60 * 1000);

			event.cookies.set('session', newJWT.accessToken, {
				expires: accessTokenExpirationDate,
				httpOnly: true,
				path: '/'
			});
			event.cookies.set('refreshToken', newJWT.refreshToken, {
				expires: refreshTokenExpirationDate,
				httpOnly: true,
				path: '/'
			});

			event.locals.userId = jsonString.data.userId;
		}
	} catch (err) {
		console.error(err);
		// delete cookie and refresh token
		event.cookies.delete('session');
		event.cookies.delete('refreshToken');
		event.locals.userId = '';
		event.locals.email = '';
		event.locals.name = '';

		throw redirect(302, '/login');
	}

	let res: { status: number; response: ServerResponse };
	// get user
	try {
		const fetchUser = await fetch(API_URL + '/user', {
			method: 'GET',
			headers: new Headers({ Authorization: session as string })
		});
		const jsonString = await fetchUser.json();
		res = { status: fetchUser.status, response: jsonString };
	} catch (err) {
		console.error(err);
		// delete cookie and refreshToken
		event.cookies.delete('session');
		event.cookies.delete('refreshToken');
		event.locals.userId = '';
		event.locals.email = '';
		event.locals.name = '';
		throw redirect(302, '/login');
	}

	// console.log('hook');
	// console.log(res.response.data);

	event.locals.userId = res.response.data.id;
	event.locals.email = res.response.data.email;
	event.locals.name = res.response.data.name;

	return await resolve(event);
};
