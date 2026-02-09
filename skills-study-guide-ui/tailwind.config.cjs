const path = require('path')

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    path.join(__dirname, 'index.html'),
    path.join(__dirname, 'src', '**', '*.{js,jsx}'),
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        paper: {
          50: '#FAFAF8',
          100: '#F5F4F0',
          200: '#EFEEEA',
          300: '#E8E7E3',
          400: '#D4D3CF',
        },
        ink: {
          50: '#E8E8E8',
          100: '#999999',
          200: '#6B6B6B',
          300: '#3A3A3A',
          400: '#2E2E2E',
          500: '#2A2A2A',
          600: '#252525',
          700: '#1E1E1E',
          800: '#1A1A1A',
          900: '#161616',
        },
      },
      fontFamily: {
        heading: ['Newsreader', 'Georgia', 'serif'],
        body: ['Source Serif 4', 'Georgia', 'serif'],
        mono: ['JetBrains Mono', 'Monaco', 'Courier New', 'monospace'],
      },
      maxWidth: {
        'prose': '680px',
        'wide': '780px',
      },
      typography: (theme) => ({
        DEFAULT: {
          css: {
            maxWidth: '680px',
            fontFamily: theme('fontFamily.body').join(', '),
            lineHeight: '1.75',
            h1: {
              fontFamily: theme('fontFamily.heading').join(', '),
              fontWeight: '800',
              fontSize: '2.5rem',
              lineHeight: '1.15',
              letterSpacing: '-0.02em',
            },
            h2: {
              fontFamily: theme('fontFamily.heading').join(', '),
              fontWeight: '600',
              fontSize: '1.75rem',
              lineHeight: '1.25',
              letterSpacing: '-0.01em',
              marginTop: '4rem',
              paddingBottom: '0.75rem',
              borderBottomWidth: '1px',
              borderBottomColor: theme('colors.paper.300'),
            },
            h3: {
              fontFamily: theme('fontFamily.heading').join(', '),
              fontWeight: '600',
              fontSize: '1.375rem',
              lineHeight: '1.3',
              marginTop: '2.5rem',
            },
            h4: {
              fontFamily: theme('fontFamily.body').join(', '),
              fontWeight: '600',
              fontSize: '1.125rem',
            },
            code: {
              fontFamily: theme('fontFamily.mono').join(', '),
              fontSize: '0.875em',
              fontWeight: '500',
            },
            'code::before': { content: 'none' },
            'code::after': { content: 'none' },
            th: {
              fontFamily: theme('fontFamily.mono').join(', '),
              fontSize: '0.6875rem',
              textTransform: 'uppercase',
              letterSpacing: '0.06em',
              fontWeight: '600',
            },
          },
        },
        invert: {
          css: {
            h2: {
              borderBottomColor: theme('colors.ink.400'),
            },
          },
        },
      }),
    },
  },
  plugins: [
    require('@tailwindcss/typography'),
  ],
}
