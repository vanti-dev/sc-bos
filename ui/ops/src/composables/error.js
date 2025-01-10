import {StatusCode} from 'grpc-web';
import {computed, toValue} from 'vue';

export const statusCodeById = Object.entries(StatusCode).reduce((acc, [key, value]) => {
  acc[value] = key;
  return acc;
}, {});

/**
 * @param {MaybeRefOrGetter<
 *   string | import('grpc-web').RpcError | Error | {error: import('grpc-web').RpcError} | any
 * >} err
 * @return {{errStr: ComputedRef<string|string|string>, showErr: ComputedRef<boolean>}}
 */
export default function useError(err) {
  const showErr = computed(() => Boolean(toValue(err)));
  const errStr = computed(() => {
    let e = toValue(err);
    if (typeof e === 'string') return e; // simple case
    if (!e) return ''; // no error
    if (e.error) e = e.error;
    let str = '';
    if (e.code) {
      str += statusCodeById[e.code] + ': ';
    }
    str += e.message ?? '';
    return str;
  });
  return {showErr, errStr};
}
