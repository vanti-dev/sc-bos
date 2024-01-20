import {useAccountStore} from '@/stores/account';
import {StatusCode} from 'grpc-web';

/**
 * @param {import('grpc-web').GrpcWebClientBaseOptions} [options]
 * @return {import('grpc-web').GrpcWebClientBaseOptions}
 */
export function clientOptions(options = {}) {
  const account = useAccountStore();

  /**
   * Handle a logout event
   */
  const handleLogout = () => {
    // Initiates the logout process without waiting for it to complete
    account.logout().catch(e => console.error('Logout failed', e));
  };

  /**
   * Handle an error from a gRPC call
   *
   * @param {Error} e
   */
  const handleError = (e) => {
    // Log the user out if we get a permission denied error (this going to clear the pinia state too)
    if (e.code === StatusCode.PERMISSION_DENIED || e.code === StatusCode.UNAUTHENTICATED) {
      handleLogout();
    }
  };

  /**
   * Triggers a token refresh if necessary and returns the new token
   *
   * @return {Promise<string>}
   */
  const refreshToken = () => {
    // Initiates the refresh token process without waiting for it to complete
    return account.refreshToken().then(() => {
      return account.authenticationDetails.token;
    });
  };


  /**
   * @template T
   * @param {*} request
   * @param {function(*): T} invoker
   * @return {Promise<T>}
   */
  const addRequestHeader = (request, invoker) => {
    return refreshToken().then(token => {
      if (token) {
        request.getMetadata()['Authorization'] = `Bearer ${token}`;
      }
      return request;
    }).then(request => {
      return invoker(request);
    });
  };

  return {
    ...options,
    unaryInterceptors: [
      ...(options.unaryInterceptors || []),
      {
        intercept(request, invoker) {
          return addRequestHeader(request, invoker).catch(e => {
            handleError(e);
            throw e;
          });
        }
      }],
    streamInterceptors: [
      ...(options.streamInterceptors || []),
      {
        intercept(request, invoker) {
          const s = new DelayedClientReadableStream(addRequestHeader(request, invoker));

          s.on('error', (err) => {
            handleError(err);
          });

          return s;
        }
      }
    ]

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
  }

  cancel() {
    this.other.then(o => o.cancel());
  }
}
