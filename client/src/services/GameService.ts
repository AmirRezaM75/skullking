import type { Player } from './../types';
import { GameCommand } from './../constants';

class GameService {
	players: Player[];

	constructor() {
		this.players = [];
	}

	handle(command: GameCommand, content: any, senderId: string) {
		if (GameCommand.Init == command) {
			this.init(content);
		}

		if (GameCommand.Joined == command) {
			this.joined(content);
		}

        if (GameCommand.Left == command) {
			this.left(content);
		}

		return this;
	}

	init(content: any) {
		content.players.forEach((player) => {
			this.addPlayer(player);
		});
	}

	joined(content: any) {
		this.addPlayer(content);
	}

    left(content: any) {
        this.deletePlayerById(content.playerId)
    }

	addPlayer(player: any) {
        console.log(player)
		// When user joins a game, it receives "INIT" and "JOINED" commands
		// at the same time, to avoid inserting auth user data twice
		// we need to check if user id is already exists or exclude
		// authenticated user id but I think using AuthService is overhead in this class
		const exists = this.checkPlayerExistsById(player.id);

		if (exists) return;

		const p: Player = {
			avatar: player.avatar,
			id: player.id,
			username: player.username
		};

		this.players.push(p);
	}

    deletePlayerById(id: string) {
        const index = this.players.findIndex(player => player.id === id);
        if (index !== -1) {
            this.players.splice(index, 1);
        }
    }
    

	checkPlayerExistsById(playerId: string): boolean {
        for (let i = 0; i < this.players.length; i++) {
            if (this.players[i].id === playerId) {
                return true;
            }
        }

		return false;
	}
}

export default GameService;
