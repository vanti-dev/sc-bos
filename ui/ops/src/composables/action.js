import {newActionTracker} from '@/api/resource.js';
import {reactive, ref, toRefs, toValue, watch} from 'vue';

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

  // What follows is a race safe way to perform requests whenever request changes.
  // We queue up tasks (calls to action) and only process the latest one.
  const activeTask = ref(null); // the task we are awaiting
  const nextTask = ref(null); // our queue, except we only care about the tip
  watch(() => toValue(request), (request) => {
    nextTask.value = {request};
  }, {immediate: true});
  watch(nextTask, (task) => {
    if (!task) return; // break the loop
    if (activeTask.value) {
      // todo: cancel active request, grpc doesn't let us do this yet
      return;
    }
    activeTask.value = task;
    nextTask.value = null;
  }, {immediate: true});
  watch(activeTask, async (task) => {
    if (!task) return; // break the loop
    try {
      const {request} = task;
      if (!request) {
        tracker.response = null;
        activeTask.value = null;
        return;
      }

      await action(request, tracker);
    } finally {
      const next = nextTask.value;
      nextTask.value = null;
      activeTask.value = next;
    }
  }, {immediate: true});

  return {
    ...toRefs(tracker),
    refresh() {
      nextTask.value = {request: toValue(request)};
    }
  };
}
