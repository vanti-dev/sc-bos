module.exports = {
  parserOptions: {
    parser: null, // override the parser pulled in from @vanti/eslint-config-vue (babel)
    sourceType: 'module',
    ecmaVersion: 'latest'
  },
  'extends': [
    '@vanti/eslint-config-vue'
  ],
  'rules': {
    'no-unused-vars': 'warn',
    // allow custom v-models
    'vue/no-v-model-argument': 0,
    // allow errors and warnings to console
    'no-console': ['warn', {allow: ['warn', 'error']}],
    // tags should have no lines between them (the default), except @example which we allow any
    // number of lines between to allow better separation of these tags
    'jsdoc/tag-lines': ['warn', 'never', {
      'startLines': 1, // this was added in a newer version of the plugin
      'tags': {'example': {lines: 'any'}}
    }]
  },
  'ignorePatterns': ['**/dist/*']
};
