import {PublicationApiPromiseClient} from '@smart-core-os/sc-api-grpc-web/traits/publication_grpc_web_pb.js';
import {
  AcknowledgePublicationRequest,
  CreatePublicationRequest,
  Publication,
  PullPublicationsRequest,
  UpdatePublicationRequest
} from '@smart-core-os/sc-api-grpc-web/traits/publication_pb.js';
import {pullResource, setCollection, trackAction} from './resource.js';

/**
 * @param {string} name
 * @param {ResourceCollection<Publication.AsObject, OpenClosePositions>} resource
 */
export function pullPublications(name, resource) {
  pullResource('Publication', resource, endpoint => {
    const api = new PublicationApiPromiseClient(endpoint);
    const stream = api.pullPublications(new PullPublicationsRequest().setName(name));
    stream.on('data', msg => {
      const changes = msg.getChangesList();
      for (const change of changes) {
        setCollection(resource, change, v => v.id);
      }
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
    const api = new PublicationApiPromiseClient(endpoint);
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
    const api = new PublicationApiPromiseClient(endpoint);
    return api.updatePublication(new UpdatePublicationRequest()
        .setName(name)
        .setVersion(publication.version)
        .setPublication(fromObject(publication)))
  });
}

/**
 * @param {AcknowledgePublicationRequest.AsObject} request
 * @param {ActionTracker} [tracker]
 * @return {Promise<Publication.AsObject>}
 */
export async function acknowledgePublication(request, tracker) {
  return trackAction('Publication.acknowledgePublication', tracker ?? {}, endpoint => {
    const api = new PublicationApiPromiseClient(endpoint);
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
 * @returns {Publication}
 */
export function fromObject(obj) {
  const publication = new Publication();
  for (const prop of ['id', 'body', 'mediaType', 'version']) {
    if (obj.hasOwnProperty(prop)) {
      publication['set' + prop[0].toUpperCase() + prop.substring(1)](obj[prop]);
    }
  }

  if (obj.hasOwnProperty('audience')) {
    const src = obj.audience;
    const audience = new Publication.Audience();
    publication.setAudience(audience);
    for (const prop of ['name', 'receipt', 'receiptRejectedReason']) {
      if (src.hasOwnProperty(prop)) {
        audience['set' + prop[0].toUpperCase() + prop.substring(1)](src[prop]);
      }
    }
  }

  return publication;
}
