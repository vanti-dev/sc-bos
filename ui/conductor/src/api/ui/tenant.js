import {trackAction} from '@/api/resource.js';
import {add} from 'date-fns';

const mockTenants = [
  {id: '1', title: 'Lebank', zones: ['L2-N']},
  {id: '2', title: 'Golden Games', zones: ['L3', 'L4']},
  {id: '3', title: 'Showbies', zones: ['L2-SE']},
];

const now = new Date();
const mockSecrets = [
  {
    id: '1',
    tenant: {id: '1', title: 'Lebank'},
    hash: null,
    note: 'Minimal Access',
    expirationTime: add(now, {days: -3}),
    firstUseTime: add(now, {days: -10}),
    lastUseTime: null,
    scopeNames: ['lights', 'energy:read']
  },
  {
    id: '2',
    tenant: {id: '2', title: 'Golden Games'},
    hash: null,
    note: 'Environmental',
    expirationTime: add(now, {days: 4}),
    firstUseTime: add(now, {days: -10}),
    lastUseTime: add(now, {days: -1}),
    scopeNames: ['lights', 'hvac']
  },
  {
    id: '3',
    tenant: {id: '3', title: 'Showbies'},
    hash: null,
    note: 'Full Access',
    expirationTime: add(now, {days: 22}),
    firstUseTime: null,
    lastUseTime: add(now, {days: -3}),
    scopeNames: ['lights', 'energy']
  },
  {
    id: '4',
    tenant: {id: '3', title: 'Showbies'},
    hash: null,
    note: 'Read-only',
    expirationTime: null,
    firstUseTime: null,
    lastUseTime: null,
    scopeNames: ['lights:read', 'energy:read']
  },
];

/**
 * @param {any} request
 * @param {ActionTracker<mockTenants>} tracker
 * @return {Promise<mockTenants>}
 */
export function listTenants(request, tracker) {
  return trackAction('Tenant.listTenants', tracker ?? {}, endpoint => {
    return {
      toObject() {
        return mockTenants
      }
    };
  })
}

export function getTenant(request, tracker) {
  const tenantId = String(request.tenantId);
  if (!tenantId) throw new Error('request.tenantId must be specified');
  return trackAction('Tenant.getTenant', tracker ?? {}, async endpoint => {
    const tenant = mockTenants.find(n => String(n.id) === tenantId);
    return {
      toObject() {
        return tenant;
      }
    };
  })
}

export function listSecrets(request, tracker) {
  const tenantId = String(request.tenantId);
  if (!tenantId) throw new Error('request.tenantId must be specified');
  return trackAction('Tenant.listSecrets', tracker ?? {}, async endpoint => {
    const secrets = mockSecrets.filter(n => String(n.tenant.id) === tenantId);
    return {
      toObject() {
        return secrets;
      }
    };
  })
}
