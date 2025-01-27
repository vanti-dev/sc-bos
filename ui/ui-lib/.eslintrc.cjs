module.exports = {
  parserOptions: {
    parser: null,
    sourceType: 'module',
    ecmaVersion: 'latest'
  },
  extends: [
    'plugin:vue/recommended',
    'plugin:jsdoc/recommended'
  ],
  ignorePatterns: ['dist/**/*'],
  env: {
    browser: true
  }
};
