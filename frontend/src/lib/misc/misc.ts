import type { User } from '$lib/types/myTypes';

export function relativeTime(strDateTime: string) {
	let result;
	try {
		const createdAtDate = new Date(strDateTime);
		const now = new Date();

		const elapsedMilliseconds = now.getTime() - createdAtDate.getTime();
		const seconds = Math.floor(elapsedMilliseconds / 1000);

		const rtf = new Intl.RelativeTimeFormat('en', { numeric: 'auto' });

		if (seconds < 60) {
			result = rtf.format(-seconds, 'second');
		} else if (seconds < 3600) {
			const minutes = Math.floor(seconds / 60);
			result = rtf.format(-minutes, 'minute');
		} else if (seconds < 86400) {
			const hours = Math.floor(seconds / 3600);
			result = rtf.format(-hours, 'hour');
		} else {
			const days = Math.floor(seconds / 86400);
			result = rtf.format(-days, 'day');
		}
	} catch (err) {
		result = '';
	}

	return result;
}

export function getOtherUserId(currentUserId: string, p: User[] | undefined) {
	let userId: string | undefined;

	if (p) {
		p.forEach((u) => {
			if (u.id != currentUserId) {
				userId = u.id;
			}
		});
	}

	return userId;
}
