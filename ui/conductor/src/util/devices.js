/**
 * Return whether the given item has the given trait.
 * Trait can either be fully qualified (`smartcore.traits.OnOff`) or just the local name (`OnOff`).
 *
 * @param {Device.AsObject} item
 * @param {string} trait
 * @return {boolean}
 */
export function hasTrait(item, trait) {
  for (const traitObj of (item?.metadata?.traitsList || [])) {
    const traitName = traitObj.name;
    if (traitName === trait) return true;
    const localName = traitName.split('.').at(-1);
    if (localName === trait) return true;
  }
  return false;
}
