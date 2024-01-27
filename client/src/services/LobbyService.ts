import { goto } from "$app/navigation";
import type { GameCreatedResponse, SomeoneJoinedLobbyResponse as SomeoneJoinedLobbyResponse, SomeoneLeftLobbyResponse, UserUpdatedResponse } from "../lobby-types";
import { EventType } from "../lobby-types";
import LobbyModel from "../objects/Lobby";

class LobbyService {
    lobby: LobbyModel | null = null

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
        }
        return this
    }

    joined(content: SomeoneJoinedLobbyResponse) {
        const lobby = content.lobby

        const model = new LobbyModel(
            lobby.id,
            lobby.name,
            lobby.players,
            lobby.creatorId,
            lobby.createdAt
        )

        this.lobby = model
    }

    left(content: SomeoneLeftLobbyResponse) {
        if (! this.lobby) return

        const playerIndex = this.lobby.players.findIndex(p => p.id === content.playerId)

        if (playerIndex !== -1) {
			this.lobby.players.splice(playerIndex, 1);
		}
    }

    userUpdated(content: UserUpdatedResponse) {
        if (! this.lobby) return

        const player = this.lobby.players.find(p => p.id === content.userId)

        if (player) {
            player.avatarId = content.avatarId
        }
    }

    gameCreated(content: GameCreatedResponse) {
        goto(`/games/${content.gameId}`)
    }

    reportError() {
        // TODO:
    }
}

export default LobbyService;