import { redirect } from '@sveltejs/kit';
import AuthService from '../../../services/AuthService';

/** @type {import('./$types').PageLoad} */
export async function load({ params }) {
	const auth = new AuthService();
	const user = auth.user();

	if (!user) {
		throw redirect(302, '/login');
	} else if (!user.verified) {
		throw redirect(302, '/verify-email');
	}
}

export const ssr = false;
