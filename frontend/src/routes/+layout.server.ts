import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ locals }) => {
	return { userData: { id: locals.userId, name: locals.name, email: locals.email } };
};
