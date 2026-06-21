/** @type {import('next').NextConfig} */
const nextConfig = {
  // 图片优化配置
  images: {
    domains: [
      'localhost',
      'images.unsplash.com',
      'your-minio-endpoint.com',
    ],
    formats: ['image/avif', 'image/webp'],
  },

  // 环境变量（客户端可访问）
  env: {
    NEXT_PUBLIC_API_URL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1',
    NEXT_PUBLIC_SITE_URL: process.env.NEXT_PUBLIC_SITE_URL || 'http://localhost:3000',
  },

  // 国际化支持
  i18n: {
    locales: ['zh-CN'],
    defaultLocale: 'zh-CN',
  },

  // React 严格模式
  reactStrictMode: true,

  // 实验性功能
  experimental: {
    // 启用部分预渲染（提高性能）
    optimizeCss: true,
  },

  // 输出配置（SSR模式）
  output: 'standalone',

  // 安全头配置
  async headers() {
    return [
      {
        source: '/(.*)',
        headers: [
          {
            key: 'X-Frame-Options',
            value: 'DENY',
          },
          {
            key: 'X-Content-Type-Options',
            value: 'nosniff',
          },
          {
            key: 'Referrer-Policy',
            value: 'origin-when-cross-origin',
          },
        ],
      },
    ];
  },
};

module.exports = nextConfig;
