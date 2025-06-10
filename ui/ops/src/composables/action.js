import {newActionTracker} from '@/api/resource.js';
import {asyncWatch} from '@/util/vue.js';
import {reactive, toRefs, toValue} from 'vue';

/**
 * @typedef {import('@/api/resource.js').ActionTracker<Res>} UseActionResponse
 * @property {function(): void} refresh - re-executes the action in the background
 * @template Res
 */

/**
 * Calls action whenever request changes, returning a reactive tracker to monitor the action.
 * If request is null or undefined, the action is not called and the trackers response is reset.
 * If request changes frequently, only the latest request is processed.
 *
 * @param {import('vue').MaybeRefOrGetter<Partial<Req> | null | undefined>} request
 * @param {function(Req, ActionTracker<Res>): Promise<Res>} action
 * @return {ToRefs<UnwrapNestedRefs<UseActionResponse<Res>>>}
 * @template Req
 * @template Res
 */
export function useAction(request, action) {
  const tracker = reactive(/** @type {ActionTracker<Res>} */ newActionTracker());
  const {refresh} = asyncWatch(() => toValue(request), async (request) => {
    if (!request) {
      tracker.response = null;
      return;
    }
    await action(request, tracker);
    // todo: cancel active request, grpc doesn't let us do this yet
    // onCleanup(() => {});
  }, {immediate: true});

  return {
    ...toRefs(tracker),
    refresh,
  };
}
