/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#f5f7ff',
          100: '#ebedfe',
          200: '#d6dbfd',
          300: '#b3bcfb',
          400: '#8a93f7',
          500: '#667eea',
          600: '#5568d3',
          700: '#4553b8',
          800: '#374295',
          900: '#2f3775',
        }
      }
    },
  },
  plugins: [],
}
