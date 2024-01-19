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

	let ticketId = ""

	const apiService = new ApiService

	const response = await apiService.createTicket()

	if (response.status === 201) {
		const data = await response.json()
		ticketId = data.id
	} else {
		throw redirect(302, '/');
	}

	return {
		gameId: params.gameId,
		ticketId: ticketId,
		authId: user.id,
		cardService: cardService
	};
}

export const ssr = false;
