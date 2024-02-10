import type {
	ListLobbiesResponse,
	LobbyNameUpdatedResponse,
	SomeoneJoinedLobbyResponse,
	SomeoneLeftLobbyResponse
} from '../lobby-types';
import { EventType } from '../lobby-types';
import LobbyModel from '../objects/Lobby';

class LobbiesService {
	lobbies: LobbyModel[] = [];

	handle(type: EventType, content: any): this {
		switch (type) {
			case EventType.Joined:
				this.joined(content);
				break;
			case EventType.Left:
				this.left(content);
				break;
			case EventType.List:
				this.list(content);
				break;
			case EventType.LobbyNameUpdated:
				this.updateName(content);
				break;
		}
		return this;
	}

	joined(content: SomeoneJoinedLobbyResponse) {
		const lobby = content.lobby;
		const index = this.lobbies.findIndex((l) => l.id === lobby.id);

		const model = new LobbyModel(
			lobby.id,
			lobby.name,
			lobby.players,
			lobby.creatorId,
			lobby.managerId,
			lobby.createdAt
		);

		if (index == -1) {
			this.lobbies.push(model);
		} else {
			this.lobbies[index] = model;
		}
	}

	left(content: SomeoneLeftLobbyResponse) {
		const lobbyIndex = this.lobbies.findIndex((l) => l.id === content.lobbyId);

		if (lobbyIndex === -1) {
			return;
		}

		const playerIndex = this.lobbies[lobbyIndex].players.findIndex(
			(p) => p.id === content.playerId
		);

		if (playerIndex !== -1) {
			this.lobbies[lobbyIndex].players.splice(playerIndex, 1);
		}

		if (this.lobbies[lobbyIndex].players.length === 0) {
			this.lobbies.splice(lobbyIndex, 1);
		}
	}

	list(lobbies: ListLobbiesResponse) {
		lobbies.forEach((lobby) => {
			this.lobbies.push(
				new LobbyModel(
					lobby.id,
					lobby.name,
					lobby.players,
					lobby.creatorId,
					lobby.managerId,
					lobby.createdAt
				)
			);
		});
	}

	updateName(content: LobbyNameUpdatedResponse) {
		const lobby = this.lobbies.find((l) => l.id === content.id);

		if (lobby) {
			lobby.name = content.name;
		}
	}
}

export default LobbiesService;
