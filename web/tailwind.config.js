/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./layouts/**/*.html",
    "./content/**/*.md",
    "./assets/js/**/*.js",
    "./static/**/*.html",
  ],
  theme: {
    extend: {
      colors: {
        // Kanvas & Permukaan
        canvas: {
          DEFAULT: '#FBFBFA',
          pure: '#FFFFFF',
        },
        surface: {
          DEFAULT: '#FFFFFF',
          subtle: '#F7F6F3',
          card: '#FFFFFF',
        },
        // Garis batas struktural (Hairline borders)
        border: {
          DEFAULT: '#EAEAEA',
          subtle: 'rgba(0, 0, 0, 0.06)',
          hover: '#CCCCCC',
        },
        // Teks & Tinta
        charcoal: {
          DEFAULT: '#111111',
          muted: '#2F3437',
        },
        secondary: '#787774',
        ghost: '#999999',

        // Muted Pastels (Aksen semantik terkontrol)
        pastel: {
          green: {
            bg: '#EDF3EC',
            text: '#346538',
            border: '#D3E5D1',
          },
          blue: {
            bg: '#E1F3FE',
            text: '#1F6C9F',
            border: '#C3E4FC',
          },
          yellow: {
            bg: '#FBF3DB',
            text: '#956400',
            border: '#F4E6B8',
          },
          red: {
            bg: '#FDEBEC',
            text: '#9F2F2D',
            border: '#F8D2D4',
          },
        },
      },
      fontFamily: {
        serif: ['Newsreader', 'Instrument Serif', 'Playfair Display', 'Georgia', 'serif'],
        sans: ['Geist Sans', 'SF Pro Display', 'Helvetica Neue', 'sans-serif'],
        mono: ['Geist Mono', 'SF Mono', 'JetBrains Mono', 'Menlo', 'monospace'],
      },
      letterSpacing: {
        'tight-title': '-0.03em',
        'subtle-tight': '-0.015em',
        'badge': '0.05em',
      },
      lineHeight: {
        'title': '1.15',
        'body': '1.6',
      },
      boxShadow: {
        'diffuse': '0 2px 8px rgba(0, 0, 0, 0.04)',
        'dropdown': '0 4px 16px rgba(0, 0, 0, 0.06)',
      },
      borderRadius: {
        'crisp': '4px',
        'card': '8px',
        'card-lg': '12px',
      },
    },
  },
  plugins: [],
};
