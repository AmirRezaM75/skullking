export type User = {
	email: string;
	username: string;
	verified: boolean;
	token: string;
};

export type CreateGameResponse = {
	id: string;
	statusCode: number;
};

export type Player = {
	avatar: string;
	id: string;
	username: string;
};
