import AuthService from './AuthService';
class ApiService {
	userServiceBaseUrl: string = import.meta.env.VITE_USER_SERVICE_URL;
	lobbyServiceBaseUrl: string = import.meta.env.VITE_LOBBY_SERVICE_URL;
	skullkingBaseURL: string = import.meta.env.VITE_SKULLKING_URL;

	private authService;

	constructor() {
		this.authService = new AuthService();
	}

	createGame(lobbyId: string): Promise<Response> {
		const user = this.authService.user();

		if (!user) throw new Error('Unauthenticated');

		return fetch(this.skullkingBaseURL + '/games', {
			method: 'POST',
			headers: {
				Authorization: `Bearer ${user.token}`,
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({
				lobbyId
			}),
		});
	}

	joinGame(gameId: string, ticketId: string): WebSocket {
		const baseURL = this.skullkingBaseURL.replace('http', 'ws');
		return new WebSocket(`${baseURL}/games/join?gameId=${gameId}&ticketId=${ticketId}`);
	}

	forgotPassword(email: string): Promise<Response> {
		return fetch(this.userServiceBaseUrl + '/forgot-password', {
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
		return fetch(this.userServiceBaseUrl + '/login', {
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
		return fetch(this.userServiceBaseUrl + '/register', {
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
		return fetch(this.userServiceBaseUrl + '/reset-password', {
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
		return fetch(this.userServiceBaseUrl + path, {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json'
			}
		});
	}

	sendEmailVerificationNotification(): Promise<Response> {
		const user = this.authService.user();

		if (!user) throw new Error('Unauthenticated');

		return fetch(this.userServiceBaseUrl + '/email/verification-notification', {
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

		return fetch(this.userServiceBaseUrl + '/users/avatar', {
			method: 'PATCH',
			headers: {
				Authorization: `Bearer ${user.token}`,
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ avatarId })
		});
	}

	getCards(): Promise<Response> {
		return fetch(this.skullkingBaseURL + '/games/cards', {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json'
			}
		});
	}

	getLobbies(ticketId: string): EventSource {
		return new EventSource(this.lobbyServiceBaseUrl + '/lobbies?ticketId=' + ticketId)
	}

	createLobby(): Promise<Response> {
		const user = this.authService.user();

		if (!user) throw new Error('Unauthenticated');

		return fetch(this.lobbyServiceBaseUrl + '/lobbies', {
			method: 'POST',
			headers: {
				Authorization: `Bearer ${user.token}`,
				'Content-Type': 'application/json'
			},
		});
	}

	joinLobby(lobbyId: string, ticketId: string): EventSource {
		return new EventSource(this.lobbyServiceBaseUrl + `/lobbies/${lobbyId}/join?ticketId=` + ticketId)
	}
	

	createTicket(): Promise<Response> {
		const user = this.authService.user();

		if (!user) throw new Error('Unauthenticated');

		return fetch(this.userServiceBaseUrl + '/tickets', {
			method: 'POST',
			headers: {
				Authorization: `Bearer ${user.token}`,
				'Content-Type': 'application/json'
			},
		});
	}
}

export default ApiService;
