import antfu from '@antfu/eslint-config'

export default antfu({
  ignores: ['*.json', '*.json/**', '*.js', '*.js/**', 'dist', '**/dist/**'],
  unocss: true,
  vue: true,
  rules: {
    'no-console': 'off',
    'eslint-comments/no-unlimited-disable': 'off',
  },
})
