import { redirect } from '@sveltejs/kit';
import AuthService from '../../../services/AuthService';
import { IntendedUrl } from '../../../constants';

/** @type {import('./$types').PageLoad} */
export async function load({ params }) {
	const auth = new AuthService();
	const user = auth.user();

	if (!user) {
		sessionStorage.setItem(IntendedUrl, window.location.href);
		throw redirect(302, '/login');
	} else if (!user.verified) {
		throw redirect(302, '/verify-email');
	}

	['start'].forEach((filename) => {
		const audio = new Audio(`/sounds/${filename}.mp3`);
		audio.preload = 'auto';
	});

	return {
		lobbyId: params.lobbyId,
		auth: user
	};
}

export const ssr = false;
