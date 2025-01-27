import * as components from './src/components';

export * from './src/components';

/** @type {import('Vue').FunctionPlugin} */
export function install(app) {
  for (const component in components) {
    app.component(component, components[component]);
  }
}
