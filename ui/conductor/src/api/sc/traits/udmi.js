import {clientOptions} from '@/api/grpcweb';
import {pullResource, setValue} from '@/api/resource';
import {UdmiServicePromiseClient} from '@sc-bos/ui-gen/proto/udmi_grpc_web_pb';
import {PullExportMessagesRequest} from '@sc-bos/ui-gen/proto/udmi_pb';

/**
 *
 * @param {string} name
 * @param {ResourceValue<MqttMessage.AsObject, MqttMessage>} resource
 */
export function pullExportMessages(name, resource) {
  if (!name) throw new Error('name must be specified');
  pullResource('UDMI.pullExportMessages', resource, endpoint => {
    const api = new UdmiServicePromiseClient(endpoint, null, clientOptions());
    const stream = api.pullExportMessages(new PullExportMessagesRequest().setName(name));
    stream.on('data', msg => {
      const obj = msg.getMessage().toObject();
      console.debug('UDMI.pullExportMessages', obj);
      setValue(resource, obj);
    });
    return stream;
  });
}
