class Time {
	sleep(ms: number) {
		return new Promise((resolve) => setTimeout(resolve, ms));
	}

	add(seconds: number): number {
		// Get the current time in epoch format
		const now = Date.now();

		return now + seconds * 1000;
	}
}

export default Time;
