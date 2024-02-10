import { redirect } from '@sveltejs/kit';
import AuthService from '../../services/AuthService';
import { IntendedUrl } from '../../constants';

export function load() {
	const auth = new AuthService();
	const user = auth.user();
	if (!user) {
		sessionStorage.setItem(IntendedUrl, window.location.href);
		throw redirect(302, '/login');
	} else if (user.verified) {
		throw redirect(302, '/');
	}

	return user;
}

export const ssr = false;
