import type LobbySidebar from "../components/LobbySidebar.svelte";
import type { ListLobbiesResponse, Player, ShowLobbyResponse, SomeoneJoinedLobbbyResponse, SomeoneLeftLobbyResponse, UserUpdatedResponse } from "../lobby-types";
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
        }
        return this
    }

    joined(content: SomeoneJoinedLobbbyResponse) {
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
}

export default LobbyService;