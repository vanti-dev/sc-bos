import parseTemplate from 'json-templates';
import {merge as _merge} from 'lodash';

/**
 * Load a template from a path and process it to apply any templating to the elements within.
 *
 * @param {string} path
 * @return {Promise<Object|null>}
 */
export async function loadConfig(path) {
  if (!path) return null; // no path, no config
  const json = await fetch(path).then(r => r.json());
  const templates = {};
  if (json.templates) {
    for (const [name, template] of Object.entries(json.templates)) {
      templates[name] = parseTemplate(template);
    }
  }

  for (let i = 0; i < json.elements.length; i++) {
    const element = json.elements[i];
    if (element.template) {
      const ref = element.template.ref;
      const context = element.template;
      const defaults = element;
      const template = templates[ref];
      // todo: validate that each required property of the template is present in context
      const templateOutput = template(context);
      json.elements[i] = _merge(defaults, templateOutput);
    }
  }

  return json;
}
