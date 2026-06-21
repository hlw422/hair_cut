module.exports = {
  env: {
    NODE_ENV: '"development"',
    API_BASE_URL: '"http://localhost:8080/api/v1"',
  },
  defineConstants: {},
  mini: {
    // 小程序开发配置
    webpackChain(chain) {
      chain.merge({
        plugin: {
          define: {
            'process.env.API_BASE_URL': JSON.stringify('http://localhost:8080/api/v1'),
          },
        },
      });
    },
  },
  h5: {
    devServer: {
      port: 10089,
      proxy: {
        '/api': {
          target: 'http://localhost:8080',
          changeOrigin: true,
          pathRewrite: { '^/api': '' },
        },
      },
    },
  },
};
