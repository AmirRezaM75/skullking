import type { CreateGameResponse } from './../types';
import AuthService from './AuthService';
class ApiService {
	baseURL = 'http://localhost:3000';

	createGame(): Promise<Response> {
		return fetch(this.baseURL + '/games', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			}
		});
	}

	forgotPassword(email: string): Promise<Response> {
		return fetch(this.baseURL + '/forgot-password', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({
				email
			})
		});
	}

	login(username: string, password: string): Promise<Response> {
		return fetch(this.baseURL + '/login', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({
				username,
				password
			})
		});
	}

	register(username: string, email: string, password: string): Promise<Response> {
		return fetch(this.baseURL + '/register', {
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
		return fetch(this.baseURL + '/reset-password', {
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
		return fetch(this.baseURL + path, {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json'
			}
		});
	}

	sendEmailVerificationNotification(): Promise<Response> {
		const authService = new AuthService();
		const user = authService.user();
		
		if (!user) throw new Error("Unauthenticated");

		return fetch(this.baseURL + '/email/verification-notification', {
			method: 'POST',
			headers: {
				Authorization: `Bearer ${user.token}`,
				'Content-Type': 'application/json'
			}
		});
	}
}

export default ApiService;
