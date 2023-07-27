import AuthService from '../services/AuthService';

export function load() {
	const auth = new AuthService();
	const user = auth.user();

	const data = {
		action: 'login'
	};

	if (user && user.verified) {
		data.action = 'play';
	}

	return data;
}

export const ssr = false;
