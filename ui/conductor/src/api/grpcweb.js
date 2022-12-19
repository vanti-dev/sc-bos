import {apiToken} from '@/api/auth.js';
import {ClientReadableStream} from 'grpc-web';

/**
 * @param {import('grpc-web').GrpcWebClientBaseOptions} [options]
 * @return {import('grpc-web').GrpcWebClientBaseOptions}
 */
export function clientOptions(options = {}) {
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
          });
        }
      }],
    streamInterceptors: [
      ...(options.streamInterceptors || []),
      {
        intercept(request, invoker) {
          return new DelayedClientReadableStream(apiToken().then(token => {
            if (token) {
              request.getMetadata()['Authorization'] = `Bearer ${token}`;
            }
            return invoker(request);
          }));
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
