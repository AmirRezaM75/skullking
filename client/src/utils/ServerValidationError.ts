type ValidationError = {
	[key: string]: string;
};

class ServerValidationError {
	errors: ValidationError;
	constructor() {
		this.errors = {};
	}
	add(name: string, message: string) {
		this.errors[name] = message;
	}
    has(name: string): boolean {
        // eslint-disable-next-line
        return this.errors.hasOwnProperty(name)
    }

	get(name: string): string {
		return this.errors[name]
	}

	clear(name: string) {
		delete this.errors[name]
	}
}

export default ServerValidationError