import {computed, ref, toValue, watch} from 'vue';

/**
 *
 * @param {FilterCtx} filterCtx
 * @param {MaybeRefOrGetter<FilterKey|Filter|null>} page
 * @return {{
 *   filter: ComputedRef<Filter|null>,
 *   key: ComputedRef<FilterKey|null>,
 *   type: ComputedRef<FilterType|null>,
 *   title: ComputedRef<string>,
 *   isSelected: ComputedRef<boolean>,
 *   clear: () => void,
 *   isNonDefaultChoice: ComputedRef<boolean>,
 *   choose: (ChoiceValue) => boolean,
 *   choice: ComputedRef<Choice>,
 *   defaultChoice: ComputedRef<Choice|null>,
 *   value: ComputedRef<ChoiceValue>,
 *   text: ComputedRef<string>,
 *   search: Ref<string>,
 *   items: ComputedRef<FilterItem[]>
 * }}
 */
export default function usePageCtx(filterCtx, page) {
  /** @type {ComputedRef<Filter|null>} */
  const filter = computed(() => {
    const p = toValue(page);
    if (p === null || p === undefined) return null;
    if (typeof p === 'string') {
      if (!Object.hasOwn(filterCtx.filtersByKey.value, p)) throw new Error(`No filter with key ${p}`);
      return filterCtx.filtersByKey.value[p];
    }
    if (typeof p === 'object') {
      if (!Object.hasOwn(p, 'key')) throw new Error('Filter object must have a key');
      return p;
    }

    throw new Error(`Invalid filter ${p}`);
  });
  /** @type {ComputedRef<FilterKey|null>} */
  const key = computed(() => filter.value?.key ?? null);
  /** @type {ComputedRef<FilterType|null>} */
  const type = computed(() => filter.value?.type ?? null);
  /** @type {ComputedRef<string>} */
  const title = computed(() => filter.value?.title ?? '');
  const isSelected = computed(() => filter.value !== null);
  /** @type {ComputedRef<Choice>} */
  const choice = computed(() => filterCtx.choices[filter.value?.key] ?? {filter: filter.value.key});
  /** @type {ComputedRef<Choice|null>} */
  const defaultChoice = computed(() => filterCtx.defaultsByKey.value[key.value]);
  /** @type {ComputedRef<ChoiceValue>} */
  const value = computed(() => choice.value?.value);
  /** @type {ComputedRef<string>} */
  const text = computed(() => {
    if (choice.value === null || choice.value === undefined) return 'unset';
    const v = choice.value.text ?? choice.value.value;
    if (v === undefined) return `${title.value}: All`;
    return `${v}`;
  });
  /**
   * @param {ChoiceValue} value
   * @return {boolean} - true if a change has been applied.
   */
  const choose = (value) => {
    if (!isSelected.value) {
      throw new Error('No filter selected');
    }
    return filterCtx.choose(key.value, value);
  };

  const clear = () => {
    if (isSelected.value) {
      filterCtx.clear(key.value);
    }
  };
  /** @type {ComputedRef<boolean>} */
  const isNonDefaultChoice = computed(() => {
    if (!isSelected.value) {
      return false;
    }
    return !filterCtx.isDefaultChoice(choice.value);
  });

  const search = ref('');
  // if the filter changes, clear the search text.
  watch(filter, () => search.value = '');
  const searchNorm = computed(() => search.value?.trim().toLowerCase() ?? '');
  const items = computed(() => {
    if (!filter.value) return [];
    if (searchNorm.value === '') return filter.value.items;
    return filter.value.items.filter(i => {
      const v = `${i?.title ?? i?.value ?? i}`;
      return v.toLowerCase().includes(searchNorm.value);
    });
  });
  return {
    filter,
    key,
    type,
    title,
    isSelected,
    choice,
    defaultChoice,
    value,
    text,
    choose,
    clear,
    isNonDefaultChoice,

    search,
    items
  };
}
