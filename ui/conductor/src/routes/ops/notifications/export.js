import {alertToObject, listAlerts} from '@/api/ui/alerts';
import {severityData} from '@/composables/notifications.js';
import {downloadCSVRows} from '@/util/downloadCSV';

const dateTimeProp = (obj, prop) => {
  const v = obj?.[prop];
  if (!v) return '';
  return `${v.toLocaleDateString()} ${v.toLocaleTimeString()}`;
};

const alertCSVHeaders = [
  {title: 'Create Time', val: (a) => dateTimeProp(a, 'createTime')},
  {title: 'Source', val: (a) => a.source},
  {title: 'Floor', val: (a) => a.floor},
  {title: 'Zone', val: (a) => a.zone},
  {title: 'Severity', val: (a) => severityData(a.severity).text},
  {title: 'Description', val: (a) => a.description},
  {title: 'Resolve Time', val: (a) => dateTimeProp(a, 'resolveTime')},
  {title: 'Acknowledged', val: (a) => a.acknowledgement ? 'Yes' : 'No'},
  {title: 'Acknowledged Time', val: (a) => dateTimeProp(a.acknowledgement, 'acknowledgeTime')},
  {title: 'Acknowledged By', val: (a) => a.acknowledgement?.author?.displayName ?? ''}
];

/**
 * Causes the browser to download a CSV file containing alerts matching the given list request.
 *
 * @param {Partial<ListAlertsRequest.AsObject>} request
 * @return {Promise<void>}
 */
export async function downloadCSV(request) {
  request.pageSize = 1000;
  const csvRows =
      /** @type {string[][]} */
      [alertCSVHeaders.map(h => h.title)];
  while (true) {
    const page = await listAlerts(request);
    for (let alert of page.alertsList) {
      alert = alertToObject(alert);
      csvRows.push(alertCSVHeaders.map(h => h.val(alert)));
    }
    if (!page.nextPageToken) break;
    request.pageToken = page.nextPageToken;
  }

  const now = new Date();
  const filename = `alerts_${now.getFullYear()}-${now.getMonth() + 1}-${now.getDate()}.csv`;
  downloadCSVRows(filename, csvRows);
}
