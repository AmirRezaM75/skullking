import type {
	Player,
	Card,
	DealResponse,
	StartBiddingResponse,
	StartPickingResponse,
	BadeResponse,
	EndBiddingResponse,
	PickedResponse,
	AnnounceTrickWinnerResponse,
	NextTrickResponse,
	ReportErrorResponse
} from './../types';
import { GameCommand, GameState } from './../constants';
import type CardService from './CardService';
import Time from '../utils/Time';

class GameService {
	// Authenticated user id
	authId: string;

	cardService: CardService;

	players: Player[];

	// Dealt cards for authenticated user
	cards: Card[] = [];

	tableCards: Card[] = [];

	state: GameState = GameState.Pending;

	round = 0;

	trick = 1;

	timer = 0;

	notifierMessage = '';

	bids: number[] = [];

	bid = 0;

	// Epoch time when we are not going to process to next commands from server
	// because we are busy to complete CSS animations.
	waiter = 0;

	roundNotifier = false;

	showCountdown = false;

	countdownColor: 'blue' | 'red' = 'blue';

	constructor(cardService: CardService, authId: string) {
		this.players = [];
		this.cardService = cardService;
		this.authId = authId;
	}

	async handle(command: GameCommand, content: any) {
		const now = Date.now();

		if (this.waiter != 0 && this.waiter > now) {
			const time = new Time();
			await time.sleep(this.waiter - now);
			this.waiter = 0;
		}

		if (GameCommand.Init == command) {
			this.init(content);
		}

		if (GameCommand.Joined == command) {
			this.joined(content);
		}

		if (GameCommand.Left == command) {
			this.left(content);
		}

		if (GameCommand.Deal == command) {
			this.deal(content);
		}

		if (GameCommand.StartBidding === command) {
			this.startBidding(content);
		}

		if (GameCommand.EndBidding === command) {
			this.endBidding(content);
		}

		if (GameCommand.Bade === command) {
			this.bade(content);
		}

		if (GameCommand.StartPicking === command) {
			this.startPicking(content);
		}

		if (GameCommand.Picked === command) {
			this.picked(content);
		}

		if (GameCommand.AnnounceTrickWinner === command) {
			this.announceTrickWinner(content);
		}

		if (GameCommand.NextTrick === command) {
			this.nextTrick(content);
		}

		if (GameCommand.ReportError == command) {
			this.reportError(content);
		}

		return this;
	}

	init(content: any) {
		this.state = content.state;

		content.players.forEach((player) => {
			this.addPlayer(player);
			if (player.dealtCards) {
				this.cards = player.dealtCards;
			}
		});
	}

	joined(content: any) {
		this.addPlayer(content);
	}

	left(content: any) {
		this.deletePlayerById(content.playerId);
	}

	// To determine when next round is started, use deal()
	nextTrick(content: NextTrickResponse) {
		this.tableCards = [];
		this.round = content.round;
		this.trick = content.trick;
	}

	deal(content: DealResponse) {
		this.tableCards = [];
		this.state = content.state;
		this.round = content.round;
		this.trick = content.trick;
		// Using map() raises typescript warning because of nullable findById
		content.cards.forEach((cardId) => {
			const card = this.cardService.findById(cardId);
			if (card) {
				this.cards.push(card);
			}
		});
	}

	startBidding(content: StartBiddingResponse) {
		// As the server will bid number 0 by default,
		// It is reasonable to indicate that to player.
		this.bid = 0;
		const now = new Date().getTime() / 1000;
		this.timer = Math.floor(content.endsAt - now);
		this.state = content.state;
		this.bids = [...Array(this.round + 1).keys()];
		this.notifierMessage = 'Bidding Time';
		this.countdownColor = 'blue';
		this.showCountdown = true;
	}

	endBidding(content: EndBiddingResponse) {
		this.showCountdown = false;
		this.bids = [];
		this.state = GameState.EndBidding;
		// I wanna show bids order by players listed on left side of the screen
		this.players.forEach((player) => {
			content.bids.forEach((bid) => {
				if (bid.playerId == player.id) {
					player.bid = bid.number;
					this.bids.push(bid.number);
				}
			});
		});
		const time = new Time();
		this.waiter = time.add(2);
	}

	bade(content: BadeResponse) {
		this.bid = content.number;
	}

	startPicking(content: StartPickingResponse) {
		this.bids = [];

		const now = new Date().getTime() / 1000;
		this.timer = Math.floor(content.endsAt - now);
		this.state = content.state;
		this.players.forEach((player) => {
			player.picking = false;
			if (player.id === content.playerId) {
				player.picking = true;
				this.notifierMessage =
					player.id === this.authId ? 'Pick your card' : `${player.username} is picking`;
				this.countdownColor = player.id === this.authId ? 'blue' : 'red';
			}
		});
		this.showCountdown = true;
	}

	picked(content: PickedResponse) {
		this.showCountdown = false;
		const index = this.cards.findIndex((card) => card.id === content.cardId);

		if (index !== -1) {
			this.cards.splice(index, 1);
		}

		const card = this.cardService.findById(content.cardId);
		if (card) {
			this.tableCards.push(card);
		}

		// When this is the last person who picked the card
		// before announcing the trick winner, we need time
		// to make sure picked-card-animation is done
		// before winner animation
		if (this.tableCards.length === this.players.length) {
			const time = new Time();
			this.waiter = time.add(2);
		}
	}

	announceTrickWinner(content: AnnounceTrickWinnerResponse) {
		this.players.forEach((player) => {
			player.picking = false;
		});

		if (content.playerId === '') {
			this.notifierMessage = `No one won the trick.`;
		} else {
			this.tableCards.forEach((card) => {
				if (card.id === content.cardId) {
					card.isWinner = true;
				}
			});

			this.players.forEach((player) => {
				if (player.id === content.playerId) {
					this.notifierMessage = `${player.username} Won the trick.`;
				}
			});
		}

		const time = new Time();
		this.waiter = time.add(2);
	}

	reportError(content: ReportErrorResponse) {
		alert(content.message);
	}

	addPlayer(player: any) {
		// When user joins a game, it receives "INIT" and "JOINED" commands
		// at the same time, to avoid inserting auth user data twice
		// we need to check if user id is already exists or exclude
		// authenticated user id but I think using AuthService is overhead in this class
		const exists = this.existsByPlayerId(player.id);

		if (exists) return;

		const p: Player = {
			avatar: '/images/avatars/' + player.avatar,
			id: player.id,
			username: player.username,
			picking: false,
			bid: 0, // TODO: Get from server
			score: 0 // TODO: Get from server
		};

		this.players.push(p);
	}

	deletePlayerById(id: string) {
		const index = this.players.findIndex((player) => player.id === id);
		if (index !== -1) {
			this.players.splice(index, 1);
		}
	}

	existsByPlayerId(playerId: string): boolean {
		for (let i = 0; i < this.players.length; i++) {
			if (this.players[i].id === playerId) {
				return true;
			}
		}

		return false;
	}
}

export default GameService;
