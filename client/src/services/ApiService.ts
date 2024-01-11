import AuthService from './AuthService';
class ApiService {
	userServiceBaseURL: string = import.meta.env.VITE_USER_SERVICE_URL;
	skullkingServiceBaseURL: string = import.meta.env.VITE_USER_SERVICE_URL;

	private authService;

	constructor() {
		this.authService = new AuthService();
	}

	createGame(): Promise<Response> {
		const user = this.authService.user();

		if (!user) throw new Error('Unauthenticated');

		return fetch(this.skullkingServiceBaseURL + '/games', {
			method: 'POST',
			headers: {
				Authorization: `Bearer ${user.token}`,
				'Content-Type': 'application/json'
			}
		});
	}

	joinGame(gameId: string, token: string): WebSocket {
		const baseURL = this.skullkingServiceBaseURL.replace('http', 'ws');
		return new WebSocket(`${baseURL}/games/join?gameId=${gameId}&token=${token}`);
	}

	forgotPassword(email: string): Promise<Response> {
		return fetch(this.userServiceBaseURL + '/forgot-password', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({
				email
			})
		});
	}

	login(identifier: string, password: string): Promise<Response> {
		return fetch(this.userServiceBaseURL + '/login', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({
				identifier,
				password
			})
		});
	}

	register(username: string, email: string, password: string): Promise<Response> {
		return fetch(this.userServiceBaseURL + '/register', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({
				username,
				email,
				password
			})
		});
	}

	resetPassword(email: string, token: string, password: string): Promise<Response> {
		return fetch(this.userServiceBaseURL + '/reset-password', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({
				email,
				token,
				password
			})
		});
	}

	verifyEmail(path: string): Promise<Response> {
		return fetch(this.userServiceBaseURL + path, {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json'
			}
		});
	}

	sendEmailVerificationNotification(): Promise<Response> {
		const user = this.authService.user();

		if (!user) throw new Error('Unauthenticated');

		return fetch(this.userServiceBaseURL + '/email/verification-notification', {
			method: 'POST',
			headers: {
				Authorization: `Bearer ${user.token}`,
				'Content-Type': 'application/json'
			}
		});
	}

	updateAvatarId(avatarId: number): Promise<Response> {
		const user = this.authService.user();

		if (!user) throw new Error('Unauthenticated');

		return fetch(this.userServiceBaseURL + '/users/avatar', {
			method: 'PATCH',
			headers: {
				Authorization: `Bearer ${user.token}`,
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({avatarId})
		});
	}

	getCards(): Promise<Response> {
		return fetch(this.skullkingServiceBaseURL + '/games/cards', {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json'
			}
		});
	}
}

export default ApiService;
