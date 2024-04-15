import {closeResource} from '@/api/resource.js';
import {toValue} from '@/util/vue.js';
import {computed, ref, watch} from 'vue';

/**
 * @template V - The value type
 * @template R - Stream response type
 * @param {MaybeRefOrGetter<V | string | null>} value
 * @param {(name: MaybeRefOrGetter<string>) => ToRefs<ResourceValue<V, R>>} toComp - typically a useMyFooTrait function
 * @return {ToRefs<ResourceValue<V, R>>}}
 */
export default function useValueOrQuery(value, toComp) {
  const _v = computed(() => toValue(value));
  const needsQuery = (v) => typeof v === 'string';

  const comp = ref(/** @type {ToRefs<ResourceValue<V, R>>} */ null);
  watch(_v, (newValue, oldValue) => {
    const isStr = needsQuery(newValue);
    const wasStr = needsQuery(oldValue);
    if (isStr && !wasStr) {
      // We only write this when going from null -> something.
      // Make sure it stays reactive if generated changes from something => something else
      comp.value = toComp(_v);
    } else if (!isStr && wasStr) {
      closeResource(comp.value?.stream);
      comp.value = null;
    }
  }, {immediate: true});

  return {
    value: computed(() => {
      if (needsQuery(_v.value)) {
        return comp.value?.value;
      } else {
        return _v.value;
      }
    }),
    loading: computed(() => {
      if (needsQuery(_v.value)) {
        return comp.value?.loading;
      } else {
        return false;
      }
    }),
    stream: computed(() => {
      if (needsQuery(_v.value)) {
        return comp.value?.stream;
      } else {
        return null;
      }
    }),
    streamError: computed(() => {
      if (needsQuery(_v.value)) {
        return comp.value?.streamError;
      } else {
        return null;
      }
    }),
    updateTime: computed(() => {
      if (needsQuery(_v.value)) {
        return comp.value?.updateTime;
      } else {
        return null;
      }
    })
  };
}
