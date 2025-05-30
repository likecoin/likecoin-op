const defaultTheme = require('tailwindcss/defaultTheme');

/** @type {import('tailwindcss').Config} */
export default {
  content: [],
  theme: {
    extend: {
      colors: {
        'like-green': '#28646c',
        'like-cyan-light': '#aaf1e7',
        likecoin: {
          black: '#222222',
          buttonbg: '#AAF1E7',
          buttontext: '#27646E',
          darkgreen: '#1A4951',
          darkgrey: '#4A4A4A',
          grey: '#D8D8D8',
          lightergrey: '#F8F8F8',
          votecolor: {
            abstain: '#B7B7B7',
            yes: '#6DCDBC',
            no: '#C72F2F',
          },
          white: '#ffffff',
        },
      },
      fontSize: {
        xs: ['0.75rem', { lineHeight: '1.25rem' }],
        sm: [
          '0.875rem',
          {
            lineHeight: '1.09375rem',
          },
        ],
        base: [
          '1rem',
          {
            lineHeight: '1.25rem',
          },
        ],
        '3xl': [
          '1.875rem',
          {
            lineHeight: '2.34375rem',
          },
        ],
      },
      fontFamily: {
        inter: ['Inter', ...defaultTheme.fontFamily.sans],
      },
    },
  },
  plugins: [],
};
