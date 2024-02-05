import colors from 'vuetify/lib/util/colors';

/**
 * Get the hex color from the class name or Vuetify theme
 *
 * @example hexColor('amber lighten-2') returns #FFD54F
 * @param {string} name
 * @param {import('vuetify').Vuetify} vuetifyInstance
 * @return {string}
 */
export const hexColor = (name, vuetifyInstance) => {
  const currentTheme = vuetifyInstance.theme.dark ? 'dark' : 'light';

  const [nameFamily, nameModifier] = name.split(' ');
  const shades = ['black', 'white', 'transparent'];
  const util = {family: null, modifier: null};

  if (shades.includes(nameFamily)) {
    util.family = 'shades';
    util.modifier = nameFamily;
  } else {
    const [firstWord, secondWord] = nameFamily.split('-');
    util.family = `${firstWord}${secondWord ?
        secondWord.charAt(0).toUpperCase() + secondWord.slice(1) :
        ''}`;
    util.modifier = nameModifier ?
        nameModifier.replace('-', '') :
        'base';
  }

  // Attempt to get the color from the Vuetify theme first
  const themeColor = vuetifyInstance?.theme?.themes[currentTheme][util.family]?.[util.modifier];

  // Fallback to Vuetify's default colors if not defined in the theme
  return themeColor || colors[util.family][util.modifier];
};


/**
 * Get the rgb color from the name
 *
 * @example rgbColor('amber lighten-2') returns rgb(255, 213, 79)
 * @param {string} name
 * @return {string}
 */
export const rgbColor = (name) => {
  return hexColor(name).replace('#', '').match(/.{2}/g).map(val => parseInt(val, 16)).join(', ');
};

/**
 * Get the rgba color from the name
 *
 * @example rgbaColor('amber lighten-2', 0.5) returns rgba(255, 213, 79, 0.5)
 * @param {string} name
 * @param {number} alpha
 * @return {string}
 */
export const rgbaColor = (name, alpha) => {
  return `rgba(${rgbColor(name)}, ${alpha})`;
};

/**
 * Set hex color opacity
 *
 * @example hexOpacity('#FFD54F', 10) returns #1AFFD54F
 * @param {string} hex
 * @param {number} opacityLevel
 * @return {string}
 */
export const hexOpacity = (hex, opacityLevel) => {
  // Ensure opacity is between 0 and 100
  const clampedOpacity = Math.max(0, Math.min(100, opacityLevel));

  // Convert the opacity level to a hexadecimal value (0 - 255)
  const opacityHex = Math.floor((clampedOpacity / 100) * 255).toString(16).padStart(2, '0');

  // Return the hex color with opacity
  return hex.replace('#', `#${opacityHex}`);
};
