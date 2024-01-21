import { redirect } from '@sveltejs/kit';
import AuthService from '../../../services/AuthService';
import CardService from '../../../services/CardService';
import { IntendedUrl } from '../../../constants';
import ApiService from '../../../services/ApiService';

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

	const cardService = new CardService();
	await cardService.import();

	return {
		gameId: params.gameId,
		authId: user.id,
		cardService: cardService
	};
}

export const ssr = false;
