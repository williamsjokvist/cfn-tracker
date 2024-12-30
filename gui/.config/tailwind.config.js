/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./src/**/*.{js,jsx,ts,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        'highlight': 'var(--color-highlight)',
        'divider': 'var(--color-divider)'
      },
      opacity: {
        'enth': 'var(--enth-only)'
      },
      fontFamily: {
        'app': 'var(--app-font)',
      },
      backgroundImage: {
        'main': 'var(--bg-main)'
      },
      animation: {
        'blink': 'blink 1s linear infinite',
      },
      keyframes: {
        blink: {
          '50%': { color: 'white' }
        }
      }
    },
  },
  plugins: [require('@tailwindcss/forms')],
}
