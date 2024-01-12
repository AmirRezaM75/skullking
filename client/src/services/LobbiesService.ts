import type { ListLobbiesResponse, SomeoneJoinedLobbbyResponse, SomeoneLeftLobbyResponse } from "../lobby-types";
import { EventType } from "../lobby-types";
import LobbyModel from "../objects/Lobby";

class LobbiesService {
    lobbies: LobbyModel[] = []

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
        }
        return this
    }

    joined(content: SomeoneJoinedLobbbyResponse) {
        const lobby = content.lobby
        const index = this.lobbies.findIndex(l => l.id === lobby.id)

        const model = new LobbyModel(
            lobby.id,
            lobby.name,
            lobby.players,
            lobby.creatorId,
            lobby.createdAt
        )

        if (index == -1) {
            this.lobbies.push(model)
        } else {
            this.lobbies[index] = model
        }
    }

    left(content: SomeoneLeftLobbyResponse) {
        const lobbyIndex = this.lobbies.findIndex((l) => l.id === content.lobbyId)

        if (lobbyIndex === -1) {
            return
        }

        const playerIndex = this.lobbies[lobbyIndex].players.findIndex(p => p.id === content.playerId)

        if (playerIndex !== -1) {
			this.lobbies[lobbyIndex].players.splice(playerIndex, 1);
		}

        if (this.lobbies[lobbyIndex].players.length === 0) {
            this.lobbies.splice(lobbyIndex, 1)
        }
    }

    list(lobbies: ListLobbiesResponse) {
        lobbies.forEach(lobby => {
            this.lobbies.push(new LobbyModel(
                lobby.id,
                lobby.name,
                lobby.players,
                lobby.creatorId,
                lobby.createdAt
            ))
        })
    }
}

export default LobbiesService;