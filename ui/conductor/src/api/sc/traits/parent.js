import {clientOptions} from '@/api/grpcweb.js';
import {pullResource, setCollection, setError} from '@/api/resource.js';
import {ParentApiPromiseClient} from '@smart-core-os/sc-api-grpc-web/traits/parent_grpc_web_pb';
import {PullChildrenRequest} from '@smart-core-os/sc-api-grpc-web/traits/parent_pb';

/**
 * @param {string} name
 * @param {ResourceCollection<Child.AsObject, Child>} resources
 */
export function pullChildren(name, resources) {
  pullResource('Parent.Children', resources, endpoint => {
    const api = new ParentApiPromiseClient(endpoint, null, clientOptions());
    const stream = api.pullChildren(new PullChildrenRequest().setName(name));
    stream.on('data', msg => {
      const changes = msg.getChangesList();
      for (const change of changes) {
        setCollection(resources, change, child => child.name);
      }
    });
    stream.on('error', err => {
      setError(resources, err);
    });
    return stream;
  });
}
