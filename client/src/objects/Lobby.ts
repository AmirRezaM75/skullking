import type { Player } from "../lobby-types";

class Lobby {
    constructor(
        public id: string,
        public name: string,
        public players: Player[],
        public creatorId: string,
        public managerId: string,
        public createdAt: number,
    ) {

    }

    getCreatedAt(): string {
        const date = new Date(this.createdAt * 1000)

        const options = {
            day: 'numeric',
            month: 'short',
            hour: 'numeric',
            minute: 'numeric',
        };

        return new Intl.DateTimeFormat('en-US', options).format(date);
    }

    getManager(): Player|null {
        for (let i = 0; i < this.players.length; i++) {
            if (this.players[i].id == this.managerId) {
                return this.players[i]
            }
        }

        return null
    }
}

export default Lobby;