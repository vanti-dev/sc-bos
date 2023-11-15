import {apiToken} from '@/api/auth.js';
import {useAccountStore} from '@/stores/account';
import {StatusCode} from 'grpc-web';

/**
 * @param {import('grpc-web').GrpcWebClientBaseOptions} [options]
 * @return {import('grpc-web').GrpcWebClientBaseOptions}
 */
export function clientOptions(options = {}) {
  const account = useAccountStore();

  return {
    ...options,
    unaryInterceptors: [
      ...(options.unaryInterceptors || []),
      {
        intercept(request, invoker) {
          return apiToken().then(token => {
            if (token) {
              request.getMetadata()['Authorization'] = `Bearer ${token}`;
            }
            return invoker(request);
          }).catch(e => {
            // Log the user out if we get a permission denied error
            // and clear the local storage
            if (e.code === StatusCode.PERMISSION_DENIED || e.code === StatusCode.UNAUTHENTICATED) {
              account.logout();
              localStorage.clear();
            }

            throw e;
          });
        }
      }],
    streamInterceptors: [
      ...(options.streamInterceptors || []),
      {
        intercept(request, invoker) {
          const s = new DelayedClientReadableStream(apiToken().then(token => {
            if (token) {
              request.getMetadata()['Authorization'] = `Bearer ${token}`;
            }
            return invoker(request);
          }));
          s.on('error', (err) => {
            if (err.code === StatusCode.PERMISSION_DENIED || err.code === StatusCode.UNAUTHENTICATED) {
              account.logout();
              localStorage.clear();
            }
          });
          return s;
        }
      }]
  };
}

/**
 * A ClientReadableStream that wraps a promise of a ClientReadableStream. This is generally ok on the surface but some
 * subtle issue might arise if, for example, the close method throws an exception which won't be propagated to the
 * caller. Similarly we assume that p.then(a); p.then(b); invokes a then b in that order when p settles, otherwise some
 * calling code might not subscribe and unsubscribe is the right order.
 *
 * @augments {ClientReadableStream}
 */
class DelayedClientReadableStream {
  /**
   * @param {Promise<ClientReadableStream>} other
   */
  constructor(other) {
    this.other = other;
  }

  on(eventType, callback) {
    this.other.then(o => o.on(eventType, callback));
    return this;
  }

  removeListener(eventType, callback) {
    this.other.then(o => o.removeListener(eventType, callback));
    ;
  }

  cancel() {
    this.other.then(o => o.cancel());
  }
}
