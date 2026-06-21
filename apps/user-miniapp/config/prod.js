module.exports = {
  env: {
    NODE_ENV: '"production"',
    API_BASE_URL: '"https://api.your-haircut.com/api/v1"',
  },
  defineConstants: {},
  mini: {
    // 小程序生产配置
    optimizeSubpackages: {
      enable: true,
    },
    webpackChain(chain) {
      chain.merge({
        plugin: {
          define: {
            'process.env.API_BASE_URL': JSON.stringify('https://api.your-haircut.com/api/v1'),
          },
        },
      });
    },
  },
  h5: {
    publicPath: 'https://cdn.your-haircut.com/user-miniapp/',
    staticDirectory: 'static',
  },
};
