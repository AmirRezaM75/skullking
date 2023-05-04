import { redirect } from '@sveltejs/kit';
import AuthService from '../../services/AuthService';

export function load() {
	const auth = new AuthService();
	const user = auth.user();
	if (!user) {
		throw redirect(302, '/login');
	} else if (user.verified) {
		throw redirect(302, '/');
	}

	return user
}

export const ssr = false;
