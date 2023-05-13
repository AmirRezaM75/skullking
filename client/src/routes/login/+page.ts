import { redirect } from '@sveltejs/kit';
import AuthService from '../../services/AuthService';

export function load() {
	const auth = new AuthService();
	const user = auth.user();
	if (user) {
        if (user.verified) {
            throw redirect(302, '/')
        } else {
            throw redirect(302, 'verify-email')
        }
	}
}

export const ssr = false;
