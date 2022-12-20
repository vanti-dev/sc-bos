import {convertProperties, fieldMaskFromObject, setProperties, timestampFromObject} from '@/api/convpb.js';
import {clientOptions} from '@/api/grpcweb.js';
import {pullResource, setCollection, trackAction} from '@/api/resource.js';
import {AlertApiPromiseClient} from '@sc-bos/ui-gen/proto/alerts_grpc_web_pb';
import {
  AcknowledgeAlertRequest,
  Alert,
  ListAlertsRequest,
  ListAlertsResponse,
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
  })
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
 * @param {AcknowledgeAlertRequest.AsObject} request
 * @param {ActionTracker<Alert.AsObject>} [tracker]
 * @return {Promise<Alert.AsObject>}
 */
export function acknowledgeAlert(request, tracker) {
  return trackAction('Alerts.acknowledgeAlert', tracker ?? {}, endpoint => {
    const api = client(endpoint);
    return api.acknowledgeAlert(acknowledgeAlertRequestFromObject(request));
  })
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
  })
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
  if (!obj) return undefined
  const dst = new ListAlertsRequest();
  setProperties(dst, obj, 'name', 'pageToken', 'pageSize')
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
 * @param {Alert.Query.AsObject} obj
 * @return {Alert.Query|undefined}
 */
function alertQueryFromObject(obj) {
  if (!obj) return undefined;
  const dst = new Alert.Query();
  setProperties(dst, obj, 'floor', 'zone', 'source', 'severityNotBefore', 'severityNotAfter', 'acknowledged');
  convertProperties(dst, obj, timestampFromObject, 'createdNotBefore', 'createdNotAfter');
  return dst;
}

function acknowledgeAlertRequestFromObject(obj) {
  if (!obj) return undefined;
  const dst = new AcknowledgeAlertRequest();
  setProperties(dst, obj, 'name', 'id', 'allowMissing', 'allowAcknowledged');
  return dst;
}
