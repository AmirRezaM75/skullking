import type { User } from '../types';
import JwtService from './JwtService';

class AuthService {
	user(): User | null {
		const item = localStorage.getItem('user');

		if (item === null) {
			return null;
		}

		// TODO: Check token expiration is greater than 1 min
		const user = JSON.parse(item) as User;

		const jwtService = new JwtService();
		const payload = jwtService.decode(user.token);
		const now = Math.round(Date.now() / 1000);
		const exp = payload.exp;

		if (now > exp) {
			return null;
		}

		return user;
	}

	save(user: User) {
		const item = JSON.stringify(user);
		localStorage.setItem('user', item);
	}

	markEmailAsVerified() {
		const user = this.user();
		if (!user) return;
		user.verified = true;
		this.save(user);
	}
}

export default AuthService;
