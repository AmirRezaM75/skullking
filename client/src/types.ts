import type { GameState, CardType, GameCommand } from './constants';

export type User = {
	id: string;
	email: string;
	username: string;
	verified: boolean;
	token: string;
};

export type CreateGameResponse = {
	id: string;
	statusCode: number;
};

export type StartBiddingResponse = {
	endsAt: number;
	state: GameState;
};

export type EndBiddingResponse = {
	bids: {
		playerId: string;
		number: number;
	}[];
};

export type BadeResponse = {
	number: number;
};

export type StartPickingResponse = {
	playerId: string;
	endsAt: number;
	cardIds: number[];
	state: GameState;
};

export type Player = {
	avatar: string;
	id: string;
	username: string;
	score: number;
	picking: boolean;
	bid: number;
	wonTricksCount: number;
};

export type DealResponse = {
	round: number;
	trick: number;
	cards: number[];
	state: GameState;
};

export type PickedResponse = {
	playerId: string;
	cardId: number;
};

export type AnnounceTrickWinnerResponse = {
	playerId: string;
	cardId: number;
};

export type NextTrickResponse = {
	round: number;
	trick: number;
};

export type AnnounceScoresResponse = {
	scores: {
		playerId: string;
		score: number;
	}[];
};

export type Card = {
	id: number;
	type: CardType;
	number: number;
	borderColor: string;
	backgroundColor: string;
	textColor: string;
	imageURL: string;
	isWinner: boolean;
	disabled: boolean;
	ownerUsername: string;
};

export type Table = {
	cards: Card[];
	hasWinner: boolean;
};

export type ReportErrorResponse = {
	message: string;
};

export type Content =
	| StartBiddingResponse
	| AnnounceTrickWinnerResponse
	| NextTrickResponse
	| PickedResponse
	| DealResponse
	| StartPickingResponse
	| BadeResponse
	| EndBiddingResponse
	| ReportErrorResponse;

export type Message = {
	command: GameCommand;
	content: Content;
};

export type Countdown = {
	id: number;
	audio: HTMLAudioElement;
};

export type PlayerResponse = {
	id: string;
	username: string;
	avatar: string;
	score: number;
	bid: number;
	handCardIds: number[];
	pickableCardIds: number[];
	wonTricksCount: number;
};

export type InitResponse = {
	round: number;
	trick: number;
	state: GameState;
	expirationTime: number;
	pickingUserId: string;
	players: PlayerResponse[];
	creatorId: string;
	tableCardIds: number[];
};

export type LeftResponse = {
	playerId: string;
};
