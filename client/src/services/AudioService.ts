class AudioService {
	private backgroundAudio: HTMLAudioElement;

	public isBackgroundAudioPlaying = false;

	constructor() {
		this.backgroundAudio = new Audio('/sounds/background.mp3');
		this.backgroundAudio.loop = true;
	}

	public async toggleBackgroundAudio(): Promise<this> {
		if (this.isBackgroundAudioPlaying) {
			await this.pauseBackgroundAudio();
		} else {
			await this.playBackgroundAudio();
		}
		return this;
	}

	private async playBackgroundAudio() {
		try {
			await this.backgroundAudio.play();
			this.isBackgroundAudioPlaying = true;
		} catch (error) {
			this.isBackgroundAudioPlaying = false;
		}
	}

	private async pauseBackgroundAudio() {
		this.backgroundAudio.pause();
		this.isBackgroundAudioPlaying = false;
	}
}

export default AudioService;
