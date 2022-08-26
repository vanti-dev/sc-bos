import {apiToken} from '@/api/auth.js';
import {parseISO} from 'date-fns';
import {Timestamp} from 'google-protobuf/google/protobuf/timestamp_pb';
import {ClientReadableStream} from 'grpc-web';

/**
 * @param {import('grpc-web').GrpcWebClientBaseOptions} [options]
 * @returns {import('grpc-web').GrpcWebClientBaseOptions}
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
 * @extends {ClientReadableStream}
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

/**
 * @param proto
 * @param obj
 * @param props
 */
export function simpleFromObject(proto, obj, ...props) {
  for (const prop of props) {
    if (obj[prop]) {
      proto[`set${prop[0].toUpperCase()}${prop.substring(1)}`](obj[prop]);
    }
  }
}

/**
 * @param proto
 * @param obj
 * @param props
 */
export function timestampsFromObject(proto, obj, ...props) {
  for (const prop of props) {
    if (obj[prop]) {
      proto[`set${prop[0].toUpperCase()}${prop.substring(1)}`](timestampsFromObject(obj[prop]));
    }
  }
}

/**
 * @param {Timestamp.AsObject|string|Date} obj
 * @return {Timestamp}
 */
export function timestampFromObject(obj) {
  if (typeof obj === 'string') return timestampFromObject(parseISO(obj));
  if (obj instanceof Date) return Timestamp.fromDate(obj);

  return new Timestamp()
      .setSeconds(obj.seconds)
      .setNanos(obj.nanos);
}


/**
 * @param {google_protobuf_timestamp_pb.Timestamp|google_protobuf_timestamp_pb.Timestamp.AsObject} ts
 * @return {Date}
 */
export function timestampToDate(ts) {
  if (ts instanceof Timestamp) return ts.toDate();
  if (ts.hasOwnProperty('nanos') && ts.hasOwnProperty('seconds')) return timestampToDate(new Timestamp().setSeconds(ts.seconds).setNanos(ts.seconds));

  // be kind
  if (ts instanceof Date) return ts;
  throw new Error('cannot convert ' + ts + ' to Date, unknown format');
}


