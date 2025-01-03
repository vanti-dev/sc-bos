import {onScopeDispose, ref, toValue, watch} from 'vue';

/**
 * Returns the natural size of an image element, once it's loaded.
 *
 * @param {MaybeRefOrGetter<HTMLImageElement>} img
 * @return {{
 *   naturalWidth: Ref<number | undefined>,
 *   naturalHeight: Ref<number | undefined>
 * }}
 */
export function useNaturalSize(img) {
  const naturalWidth = ref();
  const naturalHeight = ref();

  const h = (e) => {
    const t = e.target;
    naturalWidth.value = t.naturalWidth;
    naturalHeight.value = t.naturalHeight;
  };
  watch(() => toValue(img), (n, o) => {
    if (o) o.removeEventListener('load', h);
    if (n) {
      n.addEventListener('load', h);
      naturalWidth.value = n.naturalWidth;
      naturalHeight.value = n.naturalHeight;
    }
  });
  onScopeDispose(() => {
    const i = toValue(img);
    if (i) i.removeEventListener('load', h);
  });

  return {
    naturalWidth,
    naturalHeight
  };
}
