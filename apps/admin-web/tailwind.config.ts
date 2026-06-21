import type { Config } from 'tailwindcss';

const config: Config = {
  darkMode: ['class'],
  content: [
    './src/pages/**/*.{ts,tsx}',
    './src/components/**/*.{ts,tsx}',
    './src/layouts/**/*.{ts,tsx}',
    './index.html',
  ],
  theme: {
    extend: {
      // 品牌色彩系统
      colors: {
        primary: {
          50: '#FBF7F0',   // 极浅玫瑰金（背景）
          100: '#F5EDE0',  // 浅玫瑰金
          200: '#E8D5B8',  // 浅色
          300: '#DBC390',  // 中浅
          400: '#CDB168',  // 中等
          500: '#C8A882',  // 主品牌色 - 深玫瑰金
          600: '#B8956A',  // 深色调
          700: '#94784D',  // 更深
          800: '#6F5A3A',  // 深色
          900: '#4A3C28',  // 最深
          DEFAULT: '#C8A882',
        },
        background: {
          DEFAULT: '#FFFFFF',
          secondary: '#F7F8FA',
          dark: '#1A1A1A',
          'dark-secondary': '#252525',
        },
        border: {
          DEFAULT: '#E5E7EB',
          light: '#F3F4F6',
        },
      },
      
      // 字体系统
      fontFamily: {
        sans: ['PingFang SC', '-apple-system', 'BlinkMacSystemFont', 'Segoe UI', 'Roboto', 'sans-serif'],
      },

      // 圆角系统
      borderRadius: {
        'xl': '12px',
        '2xl': '16px',
        '3xl': '20px',
        '4xl': '24px',
      },

      // 阴影系统
      boxShadow: {
        card: '0 2px 12px rgba(0, 0, 0, 0.06)',
        'card-hover': '0 4px 20px rgba(200, 168, 130, 0.15)',
        dropdown: '0 4px 16px rgba(0, 0, 0, 0.08)',
        modal: '0 8px 32px rgba(0, 0, 0, 0.12)',
      },

      // 动画
      keyframes: {
        'fade-in': {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        'slide-up': {
          '0%': { transform: 'translateY(10px)', opacity: '0' },
          '100%': { transform: 'translateY(0)', opacity: '1' },
        },
        shimmer: {
          '100%': { transform: 'translateX(100%)' },
        },
      },
      animation: {
        'fade-in': 'fade-in 200ms ease-out',
        'slide-up': 'slide-up 300ms ease-out',
        'shimmer': 'shimmer 2s infinite',
      },
    },
  },
  plugins: [require('tailwindcss-animate')],
};

export default config;
