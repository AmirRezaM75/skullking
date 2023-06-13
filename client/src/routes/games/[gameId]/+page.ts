import { redirect } from '@sveltejs/kit';
import AuthService from '../../../services/AuthService';
import CardService from '../../../services/CardService';
import { IntendedGameId } from '../../../constants';

/** @type {import('./$types').PageLoad} */
export async function load({ params }) {
	const auth = new AuthService();
	const user = auth.user();

	if (!user) {
		sessionStorage.setItem(IntendedGameId, params.gameId)
		throw redirect(302, '/login');
	} else if (!user.verified) {
		throw redirect(302, '/verify-email');
	}

	const cardService = new CardService();
	await cardService.import();

	return {
		gameId: params.gameId,
		token: user.token,
		authId: user.id,
		cardService: cardService
	}
}

export const ssr = false;
