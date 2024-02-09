import { goto } from '$app/navigation';
import type {
	GameCreatedResponse,
	LobbyNameUpdatedResponse,
	ManagerChangedResponse,
	SomeoneJoinedLobbyResponse as SomeoneJoinedLobbyResponse,
	SomeoneLeftLobbyResponse,
	UserUpdatedResponse
} from '../lobby-types';
import { EventType } from '../lobby-types';
import LobbyModel from '../objects/Lobby';

class LobbyService {
	lobby: LobbyModel | null = null;

	authId: string;

	constructor(authId: string) {
		this.authId = authId;
	}

	handle(type: EventType, content: any): this {
		switch (type) {
			case EventType.Joined:
				this.joined(content);
				break;
			case EventType.Left:
				this.left(content);
				break;
			case EventType.UserUpdated:
				this.userUpdated(content);
				break;
			case EventType.GameCreated:
				this.gameCreated(content);
				break;
			case EventType.ReportError:
				this.reportError();
				break;
			case EventType.LobbyNameUpdated:
				this.updateName(content);
				break;
		}
		return this;
	}

	joined(content: SomeoneJoinedLobbyResponse) {
		const lobby = content.lobby;

		const model = new LobbyModel(
			lobby.id,
			lobby.name,
			lobby.players,
			lobby.creatorId,
			lobby.managerId,
			lobby.createdAt
		);

		this.lobby = model;
	}

	left(content: SomeoneLeftLobbyResponse) {
		if (!this.lobby) return;

		const playerIndex = this.lobby.players.findIndex((p) => p.id === content.playerId);

		if (playerIndex !== -1) {
			this.lobby.players.splice(playerIndex, 1);
		}

		if (this.authId === content.playerId) {
			goto('/lobbies');
		}
	}

	userUpdated(content: UserUpdatedResponse) {
		if (!this.lobby) return;

		const player = this.lobby.players.find((p) => p.id === content.userId);

		if (player) {
			player.avatarId = content.avatarId;
		}
	}

	gameCreated(content: GameCreatedResponse) {
		// In order to close SSE connection we can't use goto() method.
		window.location.href = `/games/${content.gameId}`;
	}

	reportError() {
		// TODO:
	}

	updateManager(content: ManagerChangedResponse) {
		if (this.lobby) {
			this.lobby.managerId = content.userId;
		}
	}

	updateName(content: LobbyNameUpdatedResponse) {
		if (this.lobby) {
			this.lobby.name = content.name;
		}
	}

	isManager(): boolean {
		if (this.lobby) {
			return this.lobby.managerId === this.authId;
		}

		return false;
	}
}

export default LobbyService;
