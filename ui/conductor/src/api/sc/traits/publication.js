import {setProperties} from '@/api/convpb.js';
import {clientOptions} from '@/api/grpcweb.js';
import {pullResource, setCollection, setError, trackAction} from '@/api/resource.js';
import {PublicationApiPromiseClient} from '@smart-core-os/sc-api-grpc-web/traits/publication_grpc_web_pb';
import {
  AcknowledgePublicationRequest,
  CreatePublicationRequest,
  Publication,
  PullPublicationsRequest,
  UpdatePublicationRequest
} from '@smart-core-os/sc-api-grpc-web/traits/publication_pb';

/**
 * @param {string} name
 * @param {ResourceCollection<Publication.AsObject, OpenClosePositions>} resource
 */
export function pullPublications(name, resource) {
  pullResource('Publication', resource, endpoint => {
    const api = new PublicationApiPromiseClient(endpoint, null, clientOptions());
    const stream = api.pullPublications(new PullPublicationsRequest().setName(name));
    stream.on('data', msg => {
      const changes = msg.getChangesList();
      for (const change of changes) {
        setCollection(resource, change, v => v.id);
      }
    });
    stream.on('error', err => {
      setError(resource, err);
    });
    return stream;
  });
}

/**
 * @param {string} name
 * @param {Publication.AsObject} publication
 * @param {ActionTracker<Publication.AsObject>} [tracker]
 * @return {Promise<Publication.AsObject>}
 */
export function createPublication(name, publication, tracker) {
  return trackAction('Publication.createPublication', tracker ?? {}, endpoint => {
    const api = new PublicationApiPromiseClient(endpoint, null, clientOptions());
    return api.createPublication(new CreatePublicationRequest()
        .setName(name)
        .setPublication(fromObject(publication)));
  });
}

/**
 * @param {string} name
 * @param {Publication.AsObject} publication
 * @param {ActionTracker} [tracker]
 * @return {Promise<Publication.AsObject>}
 */
export async function updatePublication(name, publication, tracker) {
  return trackAction('Publication.updatePublication', tracker ?? {}, endpoint => {
    const api = new PublicationApiPromiseClient(endpoint, null, clientOptions());
    return api.updatePublication(new UpdatePublicationRequest()
        .setName(name)
        .setVersion(publication.version)
        .setPublication(fromObject(publication)));
  });
}

/**
 * @param {AcknowledgePublicationRequest.AsObject} request
 * @param {ActionTracker} [tracker]
 * @return {Promise<Publication.AsObject>}
 */
export async function acknowledgePublication(request, tracker) {
  return trackAction('Publication.acknowledgePublication', tracker ?? {}, endpoint => {
    const api = new PublicationApiPromiseClient(endpoint, null, clientOptions());
    return api.acknowledgePublication(new AcknowledgePublicationRequest()
        .setName(request.name)
        .setId(request.id)
        .setVersion(request.version)
        .setReceipt(request.receipt ?? Publication.Audience.Receipt.ACCEPTED)
        .setReceiptRejectedReason(request.receiptRejectedReason ?? '')
        .setAllowAcknowledged(request.allowAcknowledged ?? false));
  });
}

/**
 * @param {Publication.AsObject} obj
 * @return {Publication}
 */
export function fromObject(obj) {
  if (!obj) return undefined;

  const publication = new Publication();
  setProperties(publication, obj, 'id', 'body', 'mediaType', 'version');
  publication.setAudience(audienceFromObject(obj.audience));
  return publication;
}

/**
 * @param {Publication.Audience.AsObject} obj
 * @return {undefined|Publication.Audience}
 */
export function audienceFromObject(obj) {
  if (!obj) return undefined;
  const dst = new Publication.Audience();
  setProperties(dst, obj, 'name', 'receipt', 'receiptRejectedReason');
  return dst;
}
