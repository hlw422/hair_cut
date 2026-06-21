const config = {
  projectName: 'haircut-user-miniapp',
  date: '2024-1-1',
  designWidth: 375, // 设计稿宽度（iPhone标准）
  deviceRatio: {
    640: 2.34 / 2,
    750: 1,
    828: 1.81 / 2,
    375: 2 / 1,
  },
  sourceRoot: 'src',
  outputRoot: 'dist',
  plugins: ['@tarojs/plugin-framework-react', '@tarojs/plugin-platform-weapp'],
  defineConstants: {},
  copy: {
    patterns: [],
    options: {},
  },
  framework: 'react',
  compiler: {
    type: 'webpack5',
    prebundle: { enable: false },
  },
  cache: {
    enable: false, // Webpack 持久化缓存
  },
  mini: {
    webpackChain(chain) {
      chain.merge({
        module: {
          rule: [
            {
              test: /\.s[ac]ss$/i,
              use: [{ loader: 'sass-loader' }],
            },
          ],
        },
      });
    },
    postcss: {
      pxtransform: {
        enable: true,
        configPath: '',
      },
      cssModules: {
        enable: true, // 默认为 true，如不需要 CSS Modules 可设置为 false
        config: {
          namingPattern: 'module', // 命名模式
          generateScopedName: '[name]__[local]___[hash:base64:5]',
        },
      },
    },
  },
  h5: {
    publicPath: '/',
    staticDirectory: 'static',
    es Modules: true,
    postcss: {
      autoprefixer: {
        enable: true,
        config: {},
      },
      cssModules: {
        enable: false, // H5 端默认不启用 CSS Modules
      },
    },
    devServer: {
      port: 10089,
      host: 'localhost',
      proxy: {
        '/api': {
          target: 'http://localhost:8080',
          changeOrigin: true,
          pathRewrite: { '^/api': '/api/v1' },
        },
      },
    },
  },
  rn: {
    appName: 'HairCut用户端',
    postcss: {
      cssModules: {
        enable: false, // RN 端默认禁用 CSS Modules
      },
    },
  },
};

module.exports = function (merge) {
  if (process.env.NODE_ENV === 'production') {
    return merge({}, config, require('./prod'));
  }
  return merge({}, config, require('./dev'));
};
