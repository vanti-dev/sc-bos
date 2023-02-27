module.exports = {
  'extends': [
    '@vanti/eslint-config-vue'
  ],
  'rules': {
    'no-unused-vars': 'warn',
    // allow custom v-models
    'vue/no-v-model-argument': 0
  }
};
