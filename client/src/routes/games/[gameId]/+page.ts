import { redirect } from '@sveltejs/kit';
import AuthService from '../../../services/AuthService';

/** @type {import('./$types').PageLoad} */
export function load({ params }) {
	const auth = new AuthService();
	const user = auth.user();

	if (!user) {
		throw redirect(302, '/login');
	} else if (!user.verified) {
		throw redirect(302, '/verify-email');
	}

	return {
		gameId: params.gameId,
		token: user.token
	}
}

export const ssr = false;
