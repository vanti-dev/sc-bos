import {fieldMaskFromObject, setProperties} from '@/api/convpb';
import {clientOptions} from '@/api/grpcweb';
import {pullResource, setValue, trackAction} from '@/api/resource.js';
import {MeterApiPromiseClient, MeterInfoPromiseClient} from '@sc-bos/ui-gen/proto/meter_grpc_web_pb';
import {DescribeMeterReadingRequest, PullMeterReadingsRequest} from '@sc-bos/ui-gen/proto/meter_pb';

/**
 * @param {PullMeterReadingsRequest.AsObject} request
 * @param {ResourceValue<MeterReading.AsObject, PullMeterReadingsResponse>} resource
 */
export function pullMeterReading(request, resource) {
  pullResource('MeterApi.PullMeterReadings', resource, (endpoint) => {
    const api = apiClient(endpoint);
    const stream = api.pullMeterReadings(pullMeterReadingsRequestFromObject(request));
    stream.on('data', (msg) => {
      const changes = msg.getChangesList();
      for (const change of changes) {
        setValue(resource, change.getMeterReading().toObject());
      }
    });
    return stream;
  });
}

/**
 *
 * @param {DescribeMeterReadingRequest.AsObject} request
 * @param {ActionTracker<MeterReadingSupport.AsObject>} tracker
 * @return {Promise<MeterReadingSupport.AsObject>}
 */
export function describeMeterReading(request, tracker) {
  return trackAction('MeterInfo.DescribeMeterReading', tracker ?? {}, (endpoint) => {
    const api = infoClient(endpoint);
    return api.describeMeterReading(describeMeterReadingRequestFromObject(request));
  });
}

/**
 * @param {string} endpoint
 * @return {MeterApiPromiseClient}
 */
function apiClient(endpoint) {
  return new MeterApiPromiseClient(endpoint, null, clientOptions());
}

/**
 * @param {string} endpoint
 * @return {MeterInfoPromiseClient}
 */
function infoClient(endpoint) {
  return new MeterInfoPromiseClient(endpoint, null, clientOptions());
}

/**
 * @param {PullMeterReadingsRequest.AsObject} obj
 * @return {PullMeterReadingsRequest|undefined}
 */
function pullMeterReadingsRequestFromObject(obj) {
  if (!obj) return undefined;
  const dst = new PullMeterReadingsRequest();
  setProperties(dst, obj, 'name', 'updatesOnly');
  dst.setReadMask(fieldMaskFromObject(obj.readMask));
  return dst;
}

function describeMeterReadingRequestFromObject(obj) {
  if (!obj) return undefined;
  const dst = new DescribeMeterReadingRequest();
  setProperties(dst, obj, 'name');
  return dst;
}
