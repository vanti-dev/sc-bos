import {fieldMaskFromObject, setProperties} from '@/api/convpb.js';
import {clientOptions} from '@/api/grpcweb.js';
import {pullResource, setValue, trackAction} from '@/api/resource.js';
import {
  TransportApiPromiseClient,
  TransportInfoPromiseClient
} from '@vanti-dev/sc-bos-ui-gen/proto/transport_grpc_web_pb.js';
import {DescribeTransportRequest, PullTransportRequest} from '@vanti-dev/sc-bos-ui-gen/proto/transport_pb';

/**
 * @param {Partial<PullTransportRequest.AsObject>} request
 * @param {ResourceValue<Transport.AsObject, PullTransportResponse>} resource
 */
export function pullTransport(request, resource) {
  pullResource('Transport.pullTransport', resource, endpoint => {
    const api = apiClient(endpoint);
    const stream = api.pullTransport(pullTransportRequestFromObject(request));
    stream.on('data', msg => {
      const changes = msg.getChangesList();
      for (const change of changes) {
        setValue(resource, change.getTransport().toObject());
      }
    });
    return stream;
  });
}

/**
 *
 * @param {Partial<DescribeTransportRequest.AsObject>} request
 * @param {ActionTracker<TransportSupport.AsObject>} [tracker]
 * @return {Promise<TransportSupport.AsObject>}
 */
export function describeTransport(request, tracker) {
  return trackAction('Transport.describeTransport', tracker ?? {}, endpoint => {
    const api = infoClient(endpoint);
    return api.describeTransport(describeTransportRequestFromObject(request));
  });
}

/**
 * @param {string} endpoint
 * @return {TransportApiPromiseClient}
 */
function apiClient(endpoint) {
  return new TransportApiPromiseClient(endpoint, null, clientOptions());
}

/**
 * @param {string} endpoint
 * @return {TransportInfoPromiseClient}
 */
function infoClient(endpoint) {
  return new TransportInfoPromiseClient(endpoint, null, clientOptions());
}

/**
 * @param {Partial<DescribeTransportRequest.AsObject>} obj
 * @return {undefined|DescribeTransportRequest}
 */
function describeTransportRequestFromObject(obj) {
  if (!obj) return undefined;
  const dst = new DescribeTransportRequest();
  setProperties(dst, obj, 'name');
  return dst;
}


/**
 * @param {Partial<PullTransportRequest.AsObject>} obj
 * @return {PullTransportRequest|undefined}
 */
function pullTransportRequestFromObject(obj) {
  if (!obj) return undefined;

  const dst = new PullTransportRequest();
  setProperties(dst, obj, 'name', 'updatesOnly');
  dst.setReadMask(fieldMaskFromObject(obj.readMask));
  return dst;
}

