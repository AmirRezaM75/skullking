export type Lobby = {
    id: string
    name: string
    players: Player[]
    creatorId: string
    createdAt: number
}

export type Player = {
    id: string
    username: string
    verified: boolean
    avatarId: number
}

export enum EventType {
    List = 'list',
    Joined = 'joined',
    Left = 'left',
    UserUpdated = 'user_updated',
    ReportError = 'report_error',
    GameCreated = 'game_created',
}

export type ListLobbiesResponse = Lobby[]

export type SomeoneJoinedLobbbyResponse = {
    lobby: Lobby
    userId: string
}

export type SomeoneLeftLobbyResponse = {
    lobbyId: string
    playerId: string
}

export type UserUpdatedResponse = {
    userId: string
    avatarId: number
}

export type GameCreatedResponse = {
    gameId: string
}