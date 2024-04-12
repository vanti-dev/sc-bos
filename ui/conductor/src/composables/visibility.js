import {ref} from 'vue';

/**
 * Manage the visibility of something.
 *
 * @return {{
 *   hide: (function(): boolean),
 *   show: (function(): boolean),
 *   active: Ref<boolean>,
 *   toggle: (function(): boolean)
 * }}
 */
export default function useVisibility() {
  const active = ref(false);
  const show = () => active.value = true;
  const hide = () => active.value = false;
  const toggle = () => active.value = !active.value;
  return {active, show, hide, toggle};
}
