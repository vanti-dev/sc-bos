import {pullResource, setValue} from '@/api/resource.js';
import {MeterApiPromiseClient} from '@sc-bos/ui-gen/proto/meter_grpc_web_pb';
import {clientOptions} from '@/api/grpcweb';
import {PullMeterReadingsRequest} from '@sc-bos/ui-gen/proto/meter_pb';


/**
 * @param {string} name
 * @param {ResourceValue<MeterReading.AsObject, MeterReading>} resource
 */
export function pullMeterReading(name, resource) {
  pullResource('MeterApi.MeterReading', resource, endpoint => {
    const api = new MeterApiPromiseClient(endpoint, null, clientOptions());
    const stream = api.pullMeterReadings(new PullMeterReadingsRequest().setName(name));
    stream.on('data', msg => {
      const changes = msg.getChangesList();
      for (const change of changes) {
        setValue(resource, change.getMeterReading().toObject());
      }
    });
    return stream;
  });
}
