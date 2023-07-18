import {pullResource, setValue, trackAction} from '@/api/resource.js';
import {MeterApiPromiseClient, MeterInfoPromiseClient} from '@sc-bos/ui-gen/proto/meter_grpc_web_pb';
import {clientOptions} from '@/api/grpcweb';
import {PullMeterReadingsRequest, DescribeMeterReadingRequest} from '@sc-bos/ui-gen/proto/meter_pb';

/**
 * @param {string} name
 * @param {ResourceValue<MeterReading.AsObject, MeterReading>} resource
 */
export function pullMeterReading(name, resource) {
  pullResource('MeterApi.MeterReading', resource, (endpoint) => {
    const api = new MeterApiPromiseClient(endpoint, null, clientOptions());
    const stream = api.pullMeterReadings(new PullMeterReadingsRequest().setName(name));
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
 * @param {string} name
 * @param {ActionTracker<MeterReadingSupport.AsObject>} tracker
 * @return {Promise<MeterReadingSupport.AsObject>}
 */
export function describeMeterReading(name, tracker) {
  return trackAction('MeterApi.DescribeMeterReading', tracker ?? {}, (endpoint) => {
    const api = new MeterInfoPromiseClient(endpoint, null, clientOptions());
    return api.describeMeterReading(new DescribeMeterReadingRequest().setName(name));
  });
}
