import AuthService from '../services/AuthService';

export function load() {
	const auth = new AuthService();
	const user = auth.user();

	const data = {
		isUserVerified: user && user.verified
	};

	return data;
}

export const ssr = false;
