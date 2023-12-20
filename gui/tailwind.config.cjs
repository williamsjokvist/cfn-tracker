/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./src/**/*.{js,jsx,ts,tsx}",
  ],
  theme: {
    extend: {
      fontFamily: {
        'spartan': ['"League Spartan"', 'sans-serif']
      },
      backgroundImage: {
        'radial-purple': 'radial-gradient(circle at 52.1% -29.6%, rgb(144, 17, 105) 0%, rgb(51, 0, 131) 100.2%)'
      },
      animation: {
        'blink-pulse': 'blink 1s linear infinite',
      },
      keyframes: {
        blink: {
          '50%': { color: 'white' }
        }
      }
    },
  },
  plugins: [require("@tailwindcss/forms")],
}
