import type { Config } from 'tailwindcss';

const config: Config = {
  content: [
    './src/pages/**/*.{js,ts,jsx,tsx,mdx}',
    './src/components/**/*.{js,ts,jsx,tsx,mdx}',
    './src/app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  darkMode: 'class',
  theme: {
    extend: {
      // 品牌色彩系统（与后台统一）
      colors: {
        primary: {
          DEFAULT: '#C8A882',
          50: '#FBF7F0',
          100: '#F5EDE0',
          200: '#E8D5B8',
          300: '#DBC390',
          400: '#CDB168',
          500: '#C8A882',
          600: '#B8956A',
          700: '#94784D',
          800: '#6F5A3A',
          900: '#4A3C28',
        },
      },
      
      fontFamily: {
        sans: ['PingFang SC', '-apple-system', 'BlinkMacSystemFont', 'Segoe UI', 'Roboto', 'sans-serif'],
        serif: ['Georgia', 'serif'],
      },
      
      animation: {
        'fade-in': 'fadeIn 600ms ease-out forwards',
        'slide-up': 'slideUp 500ms ease-out',
        'count-up': 'countUp 1.5s ease-out',
      },
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        slideUp: {
          '0%': { transform: 'translateY(20px)', opacity: '0' },
          '100%': { transform: 'translateY(0)', opacity: '1' },
        },
        countUp: {
          '0%': { opacity: '0', transform: 'translateY(10px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
      },
    },
  },
  plugins: [],
};

export default config;
