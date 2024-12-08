import { mergeConfig } from 'vite'
import configArcoResolverPlugin from './plugin/arcoResolver'
import configCompressPlugin from './plugin/compress'
import configImageminPlugin from './plugin/imagemin'
import configVisualizerPlugin from './plugin/visualizer'
import baseConfig from './vite.config.base'

export default mergeConfig(
  {
    mode: 'production',
    plugins: [
      configCompressPlugin('gzip'),
      configVisualizerPlugin(),
      configArcoResolverPlugin(),
      configImageminPlugin(),
    ],
    build: {
      rollupOptions: {
        output: {
          manualChunks: {
            arco: ['@arco-design/web-vue'],
            chart: ['echarts', 'vue-echarts'],
            vue: ['vue', 'vue-router', 'pinia', '@vueuse/core', 'vue-i18n'],
          },
        },
      },
      chunkSizeWarningLimit: 2000,
    },
  },
  baseConfig,
)
