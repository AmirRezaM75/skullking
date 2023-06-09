import type { Message } from '../types';

class QueueService {
	isProcessing = false;

	private messages: Message[] = [];

	push(message: Message) {
		this.messages.push(message);
	}

	pop(): Message | undefined {
		if (this.isProcessing || this.messages.length === 0) {
			return undefined;
		}

		const message = this.messages.shift();

		if (message) {
			this.isProcessing = true;
		}

		return message;
	}
}

export default QueueService;
