import {closeResource, newActionTracker, newResourceValue} from '@/api/resource.js';
import {pullOnOff, updateOnOff} from '@/api/sc/traits/on-off.js';
import {setRequestName, toQueryObject, watchResource} from '@/util/traits.js';
import {computed, onScopeDispose, reactive, toRefs, toValue} from 'vue';


/**
 * @param {MaybeRefOrGetter<string|PullOnOffRequest.AsObject>} query
 * @param {MaybeRefOrGetter<boolean>=} paused
 * @return {ToRefs<ResourceValue<OnOff.AsObject, PullOnOffResponse>>}
 */
export function usePullOnOff(query, paused = false) {
 console.debug('usePullOnOff');
 const resource = reactive(
   /** @type {ResourceValue<OnOff.AsObject, PullOnOffResponse>} */
   newResourceValue()
 );
 onScopeDispose(() => closeResource(resource));
 
 const queryObject = computed(() => toQueryObject(query));
 
 watchResource(
   () => toValue(queryObject),
   () => toValue(paused),
   (req) => {
    pullOnOff(req, resource);
    return () => closeResource(resource);
   }
 );
 
 return toRefs(resource);
}

/**
 * @typedef UpdateOnOffRequest
 * @type {number|Partial<OnOff.AsObject>|Partial<UpdateOnOffRequest.AsObject>}
 */
export function useUpdateOnOff(name) {
 const tracker = reactive(
   /** @type {ActionTracker<AirTemperature.AsObject>} */
   newActionTracker()
 );
 
 /**
  * @param {MaybeRefOrGetter<UpdateOnOffRequest>} req
  * @return {UpdateOnOffRequest.AsObject}
  */
 const toRequestObject = (req) => {
  req = toValue(req);
  if (!Object.hasOwn(req, 'state')) {
   req = {state: /** @type {AirTemperature.AsObject} */ req};
  }
  return setRequestName(req, name);
 };
 
 return {
  ...toRefs(tracker),
   updateOnOff: (req) => {
   return updateOnOff(toRequestObject(req), tracker);
  }
 };
}

/**
 * @param {MaybeRefOrGetter<OnOff.AsObject|null>} value
 * @returns {{state: ComputedRef<string|*>, table: ComputedRef<[{label: string, value: (string|*)}]>}}
 */
export function useOnOff(value) {
  const _v = computed(() => toValue(value));
  
  const state = computed(() => {
   const v = _v.value;
   if (!v) return '';
   return v.onOff?.state ?? '';
  })
 
 const table = computed(() => {
  return [{
   label: 'State',
   value: state.value
  }];
 });
 
  return {
    state,
    table
  }
}