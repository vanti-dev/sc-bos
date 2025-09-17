import {closeResource, newActionTracker, newResourceValue} from '@/api/resource.js';
import {updateAirTemperature} from '@/api/sc/traits/air-temperature.js';
import {pullOnOff, updateOnOff} from '@/api/sc/traits/on-off.js';
import {setRequestName, toQueryObject, watchResource} from '@/util/traits.js';
import {computed, onScopeDispose, reactive, toRefs, toValue} from 'vue';


/**
 * @param {MaybeRefOrGetter<string|PullOnOffRequest.AsObject>} query
 * @param {MaybeRefOrGetter<boolean>=} paused
 * @return {ToRefs<UnwrapNestedRefs<ResourceValue<OnOff.AsObject, PullOnOffResponse>>>}
 */
export function usePullOnOff(query, paused = false) {
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