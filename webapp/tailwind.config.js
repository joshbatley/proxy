const colors = require('tailwindcss/colors')

module.exports = {
  purge: ['./src/**/*.{js,jsx,ts,tsx}', './public/index.html'],
  darkMode: false, // or 'media' or 'class'
  theme: {
    extend: {
      colors: {
        gray: colors.blueGray,
        green: colors.emerald,
        orange: colors.orange,
      },
      fontSize: {
        'xxs': '.60rem'
      },
    },
  },
  variants: {
    extend: {},
  },
  plugins: [],
}
