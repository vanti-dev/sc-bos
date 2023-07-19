import {convertProperties, fieldMaskFromObject, setProperties, timestampFromObject} from '@/api/convpb.js';
import {clientOptions} from '@/api/grpcweb.js';
import {pullResource, setCollection, setValue, trackAction} from '@/api/resource';
import {AlertApiPromiseClient} from '@sc-bos/ui-gen/proto/alerts_grpc_web_pb';
import {
  AcknowledgeAlertRequest,
  Alert,
  GetAlertMetadataRequest,
  ListAlertsRequest,
  PullAlertMetadataRequest,
  PullAlertsRequest
} from '@sc-bos/ui-gen/proto/alerts_pb';

/**
 * @param {ListAlertsRequest.AsObject} request
 * @param {ActionTracker<ListAlertsResponse.AsObject>} tracker
 * @return {Promise<ListAlertsResponse.AsObject>}
 */
export function listAlerts(request, tracker) {
  return trackAction('Alerts.listAlerts', tracker ?? {}, endpoint => {
    const api = client(endpoint);
    return api.listAlerts(listAlertsRequestFromObject(request));
  });
}

/**
 * @param {PullAlertsRequest.AsObject} request
 * @param {ResourceCollection<Alert.AsObject, Alert>} resource
 */
export function pullAlerts(request, resource) {
  pullResource('Alerts.pullAlerts', resource, endpoint => {
    const api = client(endpoint);
    const stream = api.pullAlerts(pullAlertsRequestFromObject(request));
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
 * @param {GetAlertMetadataRequest.AsObject} request
 * @param {ActionTracker<AlertMetadata.AsObject>}tracker
 * @return {Promise<AlertMetadata.AsObject>}
 */
export function getAlertMetadata(request, tracker) {
  console.debug('getAlertMetadata', request);
  return trackAction('Alerts.getAlertMetadata', tracker ?? {}, endpoint => {
    const api = client(endpoint);
    return api.getAlertMetadata(getAlertMetadataRequestFromObject(request));
  });
}

/**
 * @param {PullAlertMetadataRequest.AsObject} request
 * @param {ResourceValue<AlertMetadata.AsObject, PullAlertMetadataResponse>} resource
 */
export function pullAlertMetadata(request, resource) {
  console.debug('pullAlertMetadata', request);
  pullResource('Alerts.pullAlertMetadata', resource, endpoint => {
    const api = client(endpoint);
    const stream = api.pullAlertMetadata(pullAlertMetadataRequestFromObject(request));
    stream.on('data', msg => {
      const changes = msg.getChangesList();
      for (const change of changes) {
        setValue(resource, change.getMetadata().toObject());
      }
    });
    return stream;
  });
}

/**
 * @param {AcknowledgeAlertRequest.AsObject} request
 * @param {ActionTracker<Alert.AsObject>} [tracker]
 * @return {Promise<Alert.AsObject>}
 */
export function acknowledgeAlert(request, tracker) {
  return trackAction('Alerts.acknowledgeAlert', tracker ?? {}, endpoint => {
    const api = client(endpoint);
    return api.acknowledgeAlert(acknowledgeAlertRequestFromObject(request));
  });
}

/**
 * @param {AcknowledgeAlertRequest.AsObject} request
 * @param {ActionTracker<Alert.AsObject>} [tracker]
 * @return {Promise<Alert.AsObject>}
 */
export function unacknowledgeAlert(request, tracker) {
  return trackAction('Alerts.unacknowledgeAlert', tracker ?? {}, endpoint => {
    const api = client(endpoint);
    return api.unacknowledgeAlert(acknowledgeAlertRequestFromObject(request));
  });
}

/**
 * @param {string} endpoint
 * @return {AlertApiPromiseClient}
 */
function client(endpoint) {
  return new AlertApiPromiseClient(endpoint, null, clientOptions());
}

/**
 * @param {ListAlertsRequest.AsObject} obj
 * @return {ListAlertsRequest}
 */
function listAlertsRequestFromObject(obj) {
  if (!obj) return undefined;
  const dst = new ListAlertsRequest();
  setProperties(dst, obj, 'name', 'pageToken', 'pageSize');
  dst.setReadMask(fieldMaskFromObject(obj.readMask));
  dst.setQuery(alertQueryFromObject(obj.query));
  return dst;
}

/**
 * @param {PullAlertsRequest.AsObject} obj
 * @return {PullAlertsRequest|undefined}
 */
function pullAlertsRequestFromObject(obj) {
  if (!obj) return undefined;
  const dst = new PullAlertsRequest();
  setProperties(dst, obj, 'name');
  dst.setReadMask(fieldMaskFromObject(obj.readMask));
  dst.setQuery(alertQueryFromObject(obj.query));
  return dst;
}

/**
 * @param {GetAlertMetadataRequest.AsObject} obj
 * @return {GetAlertMetadataRequest|undefined}
 */
function getAlertMetadataRequestFromObject(obj) {
  if (!obj) return undefined;
  const dst = new GetAlertMetadataRequest();
  setProperties(dst, obj, 'name');
  dst.setReadMask(fieldMaskFromObject(obj.readMask));
  return dst;
}

/**
 * @param {PullAlertMetadataRequest.AsObject} obj
 * @return {PullAlertMetadataRequest|undefined}
 */
function pullAlertMetadataRequestFromObject(obj) {
  if (!obj) return undefined;
  const dst = new PullAlertMetadataRequest();
  setProperties(dst, obj, 'name', 'updatesOnly');
  dst.setReadMask(fieldMaskFromObject(obj.readMask));
  return dst;
}

/**
 * @param {Alert.Query.AsObject} obj
 * @return {Alert.Query|undefined}
 */
function alertQueryFromObject(obj) {
  if (!obj) return undefined;
  const dst = new Alert.Query();
  setProperties(dst, obj, 'floor', 'zone', 'severity', 'source',
      'severityNotBefore', 'severityNotAfter',
      'acknowledged', 'resolved');
  convertProperties(dst, obj, timestampFromObject,
      'createdNotBefore', 'createdNotAfter',
      'resolvedNotBefore', 'resolvedNotAfter');
  return dst;
}

/**
 *
 * @param {AcknowledgeAlertRequest.AsObject} obj
 * @return {AcknowledgeAlertRequest|undefined}
 */
function acknowledgeAlertRequestFromObject(obj) {
  if (!obj) return undefined;
  const dst = new AcknowledgeAlertRequest();
  setProperties(dst, obj, 'name', 'id', 'allowMissing', 'allowAcknowledged');
  dst.setAuthor(alertAuthorFromObject(obj.author));
  return dst;
}

/**
 * @param {Alert.Acknowledgement.Author.AsObject} obj
 * @return {undefined|Alert.Acknowledgement.Author}
 */
function alertAuthorFromObject(obj) {
  if (!obj) return undefined;
  const dst = new Alert.Acknowledgement.Author();
  setProperties(dst, obj, 'id', 'displayName', 'email');
  return dst;
}
