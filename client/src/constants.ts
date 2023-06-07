export enum GameCommand {
	Init = 'INIT',
	Joined = 'JOINED',
	Left = 'LEFT',
	Start = 'START',
	Deal = 'DEAL'
}

export enum GameState {
	Pending = 'PENDING',
	Dealing = 'DEALING'
}

export enum CardType {
	King = 'king',
	Whale = 'whale',
	Kraken = 'kraken',
	Mermaid = 'mermaid',
	Parrot = 'parrot',
	Map = 'map',
	Chest = 'chest',
	Roger = 'roger',
	Pirate = 'pirate',
	Escape = 'escape'
}
