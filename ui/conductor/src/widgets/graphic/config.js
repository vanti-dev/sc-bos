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
      const render = parseTemplate(template);
      templates[name] = (ctx) => {
        const output = render(ctx);
        return replaceTags(['{[', ']}'], ['{{', '}}'], output);
      };
    }
  }

  for (let i = 0; i < json.elements.length; i++) {
    const element = json.elements[i];
    if (element.template) {
      const ref = element.template.ref;
      const context = element.template;
      const defaults = element;
      const template = templates[ref];
      if (!template) {
        console.warn(`element refers to unknown template: ${ref}`, {element});
        continue;
      }
      // todo: validate that each required property of the template is present in context
      const templateOutput = template(context);
      json.elements[i] = _merge(defaults, templateOutput);
    }
  }

  return json;
}

/**
 * @param {[string,string]} needle
 * @param {[string,string]} replacement
 * @param {T} haystack
 * @return {T}
 * @template T
 */
function replaceTags(needle, replacement, haystack) {
  switch (typeOf(haystack)) {
    case 'string':
      return replaceStringTags(needle, replacement, haystack);
    case 'object':
      return replaceObjectTags(needle, replacement, haystack);
    case 'array':
      return replaceArrayTags(needle, replacement, haystack);
    default:
      return haystack;
  }
}

/**
 * @param {[string,string]} needle
 * @param {[string,string]} replacement
 * @param {string} haystack
 * @return {string}
 */
function replaceStringTags(needle, replacement, haystack) {
  for (let i = 0; i < needle.length; i++) {
    haystack = haystack.replace(needle[i], replacement[i]);
  }
  return haystack;
}

/**
 * @param {[string,string]} needle
 * @param {[string,string]} replacement
 * @param {Object} haystack
 * @return {Object}
 */
function replaceObjectTags(needle, replacement, haystack) {
  for (const [key, value] of Object.entries(haystack)) {
    // todo: should we replace in the key too?
    haystack[key] = replaceTags(needle, replacement, value);
  }
  return haystack;
}

/**
 * @param {[string,string]} needle
 * @param {[string,string]} replacement
 * @param {Array} haystack
 * @return {Array}
 */
function replaceArrayTags(needle, replacement, haystack) {
  for (let i = 0; i < haystack.length; i++) {
    haystack[i] = replaceTags(needle, replacement, haystack[i]);
  }
  return haystack;
}

/**
 * @param {any} value
 * @return {string}
 */
function typeOf(value) {
  if (Array.isArray(value)) return 'array';
  if (value === null) return 'null';
  return typeof value;
}
