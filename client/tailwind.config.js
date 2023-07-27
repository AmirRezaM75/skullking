/** @type {import('tailwindcss').Config} */
export default {
	content: ['./src/**/*.{html,js,svelte,ts}'],
	theme: {
		extend: {
			colors: {
				lime: {
					primary: 'var(--lime-primary)',
					secondary: 'var(--lime-secondary)'
				}
			}
		}
	},
	plugins: []
};
