export enum GameCommand {
	// Client side commands
	Start = 'START',
	Bid = 'BID',
	Pick = 'PICK',
	// Server side commands
	Init = 'INIT',
	Joined = 'JOINED',
	Left = 'LEFT',
	Deal = 'DEAL',
	AnnounceScores = 'ANNOUNCE_SCORES',
	StartBidding = 'START_BIDDING',
	EndBidding = 'END_BIDDING',
	Bade = 'BADE',
	StartPicking = 'START_PICKING',
	Picked = 'PICKED',
	AnnounceTrickWinner = 'ANNOUNCE_TRICK_WINNER',
	NextTrick = 'NEXT_TRICK',
	ReportError = 'REPORT_ERROR',
	EndGame = 'END_GAME'
}

export enum GameState {
	Pending = 'PENDING',
	Dealing = 'DEALING',
	Picking = 'PICKING',
	Bidding = 'BIDDING',
	EndBidding = 'END_BIDDING', // TODO: Not set in server
	AnnounceTrickWinner = 'ANNOUNCE_TRICK_WINNER' // TODO: Not set in server
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

export const IntendedGameId = 'intended_game_id'