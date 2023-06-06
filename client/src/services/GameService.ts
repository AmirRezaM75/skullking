import type { Player, DealResponse } from './../types';
import { GameCommand, GameState } from './../constants';
import Time from './../utils/Time';

class GameService {
	players: Player[];

	// Dealt card ids for authenticated user
	cards: number[] = [];

	state: GameState = GameState.Pending;

	round = 0;

	trick = 0;

	roundNotifier = false;

	constructor() {
		this.players = [];
	}

	async handle(command: GameCommand, content: any, senderId: string) {
		if (GameCommand.Init == command) {
			this.init(content);
		}

		if (GameCommand.Joined == command) {
			this.joined(content);
		}

		if (GameCommand.Left == command) {
			this.left(content);
		}

		if (GameCommand.Deal == command) {
			await this.deal(content);
		}

		return this;
	}

	init(content: any) {
		this.state = content.state;

		content.players.forEach((player) => {
			this.addPlayer(player);
			if (player.dealtCards) {
				this.cards = player.dealtCards
			}
		});
	}

	joined(content: any) {
		this.addPlayer(content);
	}

	left(content: any) {
		this.deletePlayerById(content.playerId);
	}

	async deal(content: DealResponse) {
		this.state = content.state;
		this.round = content.round;
		this.trick = content.trick;
		this.cards = content.cards;
	}

	addPlayer(player: any) {
		// When user joins a game, it receives "INIT" and "JOINED" commands
		// at the same time, to avoid inserting auth user data twice
		// we need to check if user id is already exists or exclude
		// authenticated user id but I think using AuthService is overhead in this class
		const exists = this.existsByPlayerId(player.id);

		if (exists) return;

		const p: Player = {
			avatar: '/images/avatars/' + player.avatar,
			id: player.id,
			username: player.username,
			picking: false,
			bids: 0, // TODO: Get from server
			score: 0 // TODO: Get from server
		};

		this.players.push(p);
	}

	deletePlayerById(id: string) {
		const index = this.players.findIndex((player) => player.id === id);
		if (index !== -1) {
			this.players.splice(index, 1);
		}
	}

	existsByPlayerId(playerId: string): boolean {
		for (let i = 0; i < this.players.length; i++) {
			if (this.players[i].id === playerId) {
				return true;
			}
		}

		return false;
	}
}

export default GameService;