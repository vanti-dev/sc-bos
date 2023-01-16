module.exports = {
  'extends': [
    '@vanti/eslint-config-vue'
  ],
  'rules': {
    'no-unused-vars': 'warn',
    // todo: remove this once we publish the new Vanti eslint config
    'jsdoc/tag-lines': ['warn', 'never', {
      'tags': {'example': {lines: 'any'}}
    }]
  }
};
