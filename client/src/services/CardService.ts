import type { Card } from './../types';
import { CardType } from './../constants';
import ApiService from './ApiService';

interface Response {
	items: {
		id: number;
		number: number;
		type: string;
	}[];
}

class CardService {
	cards: Card[] = [];

	async import() {
		const apiService = new ApiService();
		await apiService
			.getCards()
			.then((response) => response.json())
			.then((data: Response) => {
				data.items.forEach((c) => {
					this.cards.push({
						id: c.id,
						type: c.type as CardType,
						number: c.number,
						borderColor: this.getBorderColor(c.type),
						backgroundColor: this.getBackgroundColor(c.type),
						textColor: this.getTextColor(c.type),
						imageURL: this.getImageURL(c.type),
						isWinner: false,
						disabled: false,
						ownerUsername: ''
					});
				});
			});
	}

	// I could create Card model and put these logic as getters in the class
	// But we didn't have models directory so far and I didn't want to add
	// another concept to the project. second only calculate these stuff once not for each round.
	getBorderColor(type: string): string {
		switch (type) {
			case CardType.Chest:
				return 'border-yellow-500';
			case CardType.Parrot:
				return 'border-green-500';
			case CardType.Map:
				return 'border-purple-500';
			case CardType.Roger:
				return 'border-black';
			default:
				return 'border-white';
		}
	}

	// I have to be explicit about class name otherwise tailwincss purges the classes
	// @link https://tailwindcss.com/docs/optimizing-for-production#writing-purgeable-html
	// @link https://stackoverflow.com/questions/61312762/purgecss-cant-recognize-conditional-classes?answertab=scoredesc
	getBackgroundColor(type: string): string {
		switch (type) {
			case CardType.Chest:
				return 'bg-yellow-500';
			case CardType.Parrot:
				return 'bg-green-500';
			case CardType.Map:
				return 'bg-purple-500';
			case CardType.Roger:
				return 'bg-black';
			default:
				return 'bg-white';
		}
	}

	getTextColor(type: string): string {
		switch (type) {
			case CardType.Roger:
				return 'text-white';
			default:
				return 'text-black';
		}
	}

	getImageURL(type: string): string {
		return `/images/cards/${type}.jpg`;
	}

	findById(id: number): Card | null {
		for (let i = 0; i < this.cards.length; i++) {
			if (this.cards[i].id === id) {
				return this.cards[i];
			}
		}
		return null;
	}
}

export default CardService;
