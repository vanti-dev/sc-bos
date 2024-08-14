import choiceRangeStr from '@/components/filter/choiceRange.js';
import useVisibility from '@/composables/visibility.js';
import {toValue} from '@/util/vue.js';
import {computed, reactive, ref, watch} from 'vue';
import deepEqual from 'fast-deep-equal';

/**
 * @typedef {Object} Options
 * @property {Filter[]} filters
 * @property {Choice[]=} defaults
 */

/** @typedef {string} FilterKey */
/** @typedef {'boolean'|'list'|'range'} FilterType */
/**
 * @typedef {Object} Filter
 * @property {FilterKey} key - Identifier for the filter, useful for processing choices
 * @property {string} title - For example "Severity"
 * @property {string} icon - For example "mdi-alert"
 * @property {FilterType} type - Type of UI to show for this filter
 * @property {FilterItem[]} items - For type 'list' or 'range', the items to choose from.
 * @property {(value: ChoiceValue) => string} [valueToString] - Convert a chosen value into a string.
 *   For boolean defaults to `${title}: Yes|No`,
 *   for list defaults to `${title}: ${value.title}`,
 *   for range defaults to choiceRangeStr on the value.
 */
/**
 * @typedef {FilterItemObject|string} FilterItem
 */
/**
 * @typedef {Object} FilterItemObject
 * @property {string=} title - For example "Critical", the display text. Defaults to value.
 * @property {*} value - For example "critical", the value passed back to the caller when selected.
 */
/**
 * @typedef {Object} ChoiceRange
 * @property {FilterItem=} from
 * @property {FilterItem=} to
 */
/**
 * @typedef ChoiceValue
 * @type {
 *   boolean |
 *   FilterItem |
 *   ChoiceRange
 * }
 */
/**
 * @typedef {Object} Choice
 * @property {FilterKey} filter - the filter this choice is for
 * @property {ChoiceValue=} value - the chosen value, or undefined if no value is chosen
 * @property {string=} text - Text representation of the choice.
 */

/**
 * @typedef {Object} FilterCtx
 * @property {Ref<boolean>} active - is the menu visible or not
 * @property {() => void} show - show the menu
 * @property {() => void} hide - hide the menu
 * @property {() => void} toggle - toggle the menu
 * @property {ComputedRef<boolean>} badgeShown - is the badge shown
 * @property {ComputedRef<string>} badgeColor - the color of the badge
 * @property {{[key: FilterKey]: Choice}} choices - all choices
 * @property {(key: FilterKey, value?: ChoiceValue) => boolean} choose - choose a value for a filter
 * @property {ComputedRef<Choice[]>} sortedChoices - all choices in the order of the filters
 * @property {(c: Choice, defaults?: {[key: FilterKey]: Choice}) => boolean} isDefaultChoice - is the choice the default
 * @property {ComputedRef<Choice[]>} nonDefaultChoices - all choices that have a value, where that value is not the
 *   default
 * @property {ComputedRef<boolean>} hasNonDefaultChoices - are there any choices that have a value, where that value is
 * @property {{[key: FilterKey]: Choice}} defaultsByKey - the default choices indexed by filter key
 * @property {ComputedRef<Filter[]>} filters - all filters
 * @property {ComputedRef<{[key: FilterKey]: Filter}>} filtersByKey - all filters indexed by key
 * @property {Ref<string>} filterSearch - the search term for filtering filters
 * @property {Filter[]} displayFilters - the filters that match the search term
 * @property {(key?: FilterKey|null) => void} clear - clear all choices or a specific choice
 */

/**
 * @param {MaybeRefOrGetter<Options>} opts
 * @return {FilterCtx}
 */
export default function useFilterCtx(opts) {
  const {active, show, hide, toggle} = useVisibility();

  /**
   * @param {Filter[]} filters
   * @return {{[key: FilterKey]: Filter}}
   */
  const indexFilters = (filters) => filters.reduce((acc, f) => {
    acc[f.key] = f;
    return acc;
  }, /** @type {{[key: FilterKey]: Filter}} */ {});
  const filters = computed(() => toValue(opts).filters
      .map(((f, i) => {
        if (f.key) return f;
        return {...f, key: `i${i}`};
      })));
  const filtersByKey = computed(() => {
    return indexFilters(filters.value);
  });

  const filterSearch = ref('');
  const filterSearchNorm = computed(() => filterSearch.value?.trim().toLowerCase() ?? '');
  const displayFilters = computed(() => {
    if (filterSearchNorm.value === '') {
      return filters.value;
    }
    return filters.value.filter(f => f.title.toLowerCase().includes(filterSearchNorm.value));
  });

  const defaultsByKey = computed(() => toValue(opts).defaults?.reduce((acc, c) => {
    acc[c.filter] = c;
    return acc;
  }, /** @type {{[key: FilterKey]: Choice}} */ {}) ?? {});

  // Entries are managed by watching opts to keep them in sync without forgetting choices.
  // The various forms will update choices['someFilter'].value = 'someValue' to keep track of the selected choice.
  const choices = reactive(
      /** @type {{[key: FilterKey]: Choice}} */
      {}
  );
  watch(filters, (newFilters, oldFilters) => {
    // Don't use filtersByKey because it's computed from the same object this watcher is triggered by.
    // We don't know that when we look in it to check diffs that it hasn't already changed to match the new opts.
    const oldFiltersByKey = indexFilters(oldFilters ?? []);

    const toAdd = /** @type {Filter[]} */ [];
    const toRemove = new Set(Object.keys(choices));
    for (const f of newFilters) {
      toRemove.delete(f.key);

      if (!choices.hasOwnProperty(f.key)) {
        toAdd.push(f);
      } else {
        // check if the option and the existing choice are still compatible
        const choice = choices[f.key];
        if (f.type !== oldFiltersByKey[choice.filter]?.type) {
          // if the filter type isn't the same then we have to replace it and clear the choice.
          toRemove.add(f.key);
          toAdd.push(f);
        } else if (choice.value === undefined) {
          // if no item is selected then we don't have to do anything
        } else if (f.type === 'list') {
          // for lists, we maintain the choice if the new filter includes the item
          const item = /** @type {FilterItem} */ choice.value;
          if (!f.items.some(i => i === item || i.value === item)) {
            // replace
            toRemove.add(f.key);
            toAdd.push(f);
          }
        } else if (f.type === 'range') {
          // for ranges, we maintain the choice if the new filter includes the ranges from and to
          const range = /** @type {ChoiceRange} */ choice.value;
          const hasFrom = f.items.some(i => range.from === undefined || range.from === i);
          const hasTo = f.items.some(i => range.to === undefined || range.to === i);
          if (!(hasFrom && hasTo)) {
            // part of the range is now invalid, replace
            toRemove.add(f.key);
            toAdd.push(f);
          }
        }
      }
    }

    toRemove.forEach(key => delete(choices[key]));
    const defaults = defaultsByKey.value;
    toAdd.forEach(f => choices[f.key] = {filter: f.key, value: defaults[f.key]?.value});
  }, {deep: true, immediate: true});

  /**
   * @param {FilterKey} key
   * @param {ChoiceValue=} value - the value must match the filter type, null or undefined clears the choice.
   * @return {boolean} - if the value changed
   */
  const choose = (key, value) => {
    if (!choices.hasOwnProperty(key)) {
      throw new Error(`No filter with key ${key}`);
    }
    if (value === undefined || value === null) {
      const wasSet = choices[key].value !== undefined;
      choices[key].value = undefined;
      delete(choices[key].text);
      return wasSet;
    }

    // type check the value matches the type of filter
    const filter = filtersByKey.value[key];
    // sanity check
    if (!filter) {
      throw new Error(`No filter with key ${key}, this shouldn't happen as we've already checked for the choice`);
    }

    if (filter.type === 'boolean') {
      if (typeof value !== 'boolean') {
        throw new Error(`Expected boolean for filter ${key}, got ${typeof value}`);
      }
    } else if (filter.type === 'list') {
      if (!filter.items.some(i => i === value || deepEqual(i.value, value))) {
        throw new Error(`Invalid value ${value} for filter ${key}`);
      }
    } else if (filter.type === 'range') {
      if (typeof value !== 'object') {
        throw new Error(`Expected object for filter ${key}, got ${typeof value}`);
      }
      if (value.from !== undefined && !filter.items.some(i => i === value.from || deepEqual(i.value, value.from))) {
        throw new Error(`Invalid from value ${value.from} for filter ${key}`);
      }
      if (value.to !== undefined && !filter.items.some(i => i === value.to || deepEqual(i.value, value.to))) {
        throw new Error(`Invalid to value ${value.to} for filter ${key}`);
      }
    } else {
      throw new Error(`Unknown filter type ${filter.type}, a developer needs to update this code`);
    }

    const changed = !deepEqual(choices[key].value, value);
    choices[key].value = value;

    // set text if we can.
    if (filter.valueToString) {
      choices[key].text = filter.valueToString(value);
    } else {
      switch (filter.type) {
        case 'boolean':
          const boolS = value ? 'Yes' : 'No';
          choices[key].text = `${filter.title}: ${boolS}`;
          break;
        case 'list':
          const listS = value.title ?? value.value ?? value;
          choices[key].text = `${listS}`;
          break;
        case 'range':
          choices[key].text = choiceRangeStr(value);
          break;
        default:
          delete(choices[key].text);
      }
    }
    return changed;
  };

  /**
   * @param {FilterKey|null=} key - the choice to clear, null to clear all.
   */
  const clear = (key = null) => {
    if (key === null) {
      for (const k of Object.keys(choices)) {
        clear(k);
      }
    } else {
      const defaultVal = defaultsByKey.value[key]?.value;
      choose(key, defaultVal);
    }
  };

  /** @type {ComputedRef<Choice[]>} */
  const sortedChoices = computed(() => filters.value
      .map(f => choices[f.key]));

  /**
   * @param {Choice} c
   * @param {{[key: FilterKey]: Choice}=} defaults
   * @return {boolean}
   */
  const isDefaultChoice = (c, defaults = defaultsByKey.value) => {
    const defaultChoice = defaults[c.filter];
    return deepEqual(c.value, defaultChoice?.value);
  };
  /**
   * All choices that have a value, where that value is not the default.
   *
   * @type {ComputedRef<Choice>}
   */
  const nonDefaultChoices = computed(() => {
    // call defaultsByKey.value once instead of in each iteration
    const defaults = defaultsByKey.value;
    return sortedChoices.value
        .filter(c => !isDefaultChoice(c, defaults));
  });

  const hasNonDefaultChoices = computed(() => nonDefaultChoices.value.length > 0);
  const badgeShown = hasNonDefaultChoices;
  const badgeColor = computed(() => 'red');

  return {
    active,
    show, hide, toggle,

    badgeShown, badgeColor,

    choose, clear,
    choices, sortedChoices, nonDefaultChoices,
    hasNonDefaultChoices, isDefaultChoice,
    defaultsByKey,

    filters,
    filtersByKey,

    filterSearch,
    displayFilters
  };
}

export const filterCtxSymbol =
    /** @type {import('vue').InjectionKey<FilterCtx>} */
    Symbol('filterCtx');
