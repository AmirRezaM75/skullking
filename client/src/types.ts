export type User = {
	email: string;
	username: string;
	verified: boolean;
	token: string;
};

export type CreateGameResponse = {
	id: string
	statusCode: number
}