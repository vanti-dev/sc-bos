module.exports = {
  'extends': [
    '@vanti/eslint-config-vue'
  ],
  'rules': {
    'no-unused-vars': 'warn',
    // allow custom v-models
    'vue/no-v-model-argument': 0,
    // allow errors and warnings to console
    'no-console': ['warn', {allow: ['warn', 'error']}]
  },
  'ignorePatterns': ['**/dist/*']
};
