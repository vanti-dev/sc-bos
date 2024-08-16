module.exports = {
  parserOptions: {
    parser: null,
    sourceType: 'module',
    ecmaVersion: 'latest'
  },
  ignorePatterns: ['dist/**/*'],
  env: {
    browser: true
  },
  // https://github.com/feross/standard/blob/master/RULES.md#javascript-standard-style
  extends: [
    'google',
    'plugin:vue/recommended',
    'plugin:vuetify/recommended',
    'plugin:jsdoc/recommended'
  ],
  rules: {
    // allow paren-less arrow functions
    'arrow-parens': 0,
    // we don't like trailing commas
    'comma-dangle': ['error', 'never'],
    // allow async-await
    'generator-star-spacing': 0,
    'linebreak-style': 0,
    // we have modern monitors these days
    'max-len': ['error', 120],
    // allow console debugger during development
    'no-debugger': process.env.NODE_ENV === 'production' ? 'error' : 'off',
    'vue/attributes-order': 0,
    // let us call things single words
    'vue/multi-word-component-names': 0,
    'vue/html-closing-bracket-newline': ['error', {
      'multiline': 'never'
    }],
    'vue/html-closing-bracket-spacing': ['error', {
      'selfClosingTag': 'never'
    }],
    'vue/html-indent': ['error', 2, {
      'attribute': 2
    }],
    'vue/max-attributes-per-line': [2, {
      'singleline': 10,
      'multiline': 1
    }],
    // This rule is in conflict with the vue/no-v-for-template-key-on-child rule, which should be used instead.
    // This rule seems to be enabled by default, even though we're using vue 3, not sure why.
    'vue/no-v-for-template-key': 0,
    'vue/first-attribute-linebreak': ['error', {
      'singleline': 'beside'
    }],
    'vue/singleline-html-element-content-newline': ['off'],
    // 'turn on' script-setup-uses-vars to allow declaration of vars in <script setup> tags
    'vue/script-setup-uses-vars': 'error',
    'vuetify/no-deprecated-colors': ['error', {
      // from src/plugins/vuetify.js
      themeColors: ['primary', 'secondary', 'primaryTeal', 'accent', 'neutral', 'error', 'success', 'info', 'warning']
    }],
    // we use the jsdoc plugin instead
    'valid-jsdoc': 'off',
    'require-jsdoc': 'off',
    // because we haven't implemented modules correctly across most of our code
    'jsdoc/no-undefined-types': 0,
    // because it's the types that are more important, it's still preferred but we don't want
    // all the explicit warnings
    'jsdoc/require-param-description': 0,
    'jsdoc/require-property-description': 0,
    'jsdoc/require-returns-description': 0,
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
  settings: {
    jsdoc: {
      // Enable import('foo').Foo syntax
      'mode': 'typescript',
      'tagNamePreference': {
        'returns': 'return'
      },
      'preferredTypes': {
        'object': 'Object',
        'Function': 'function'
      }
    }
  }
};
