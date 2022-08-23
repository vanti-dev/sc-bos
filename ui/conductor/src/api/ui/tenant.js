import {trackAction} from '@/api/resource.js';

const mockTenants = [
  {id: '1', title: 'Lebank', zones: ['L2-N']},
  {id: '2', title: 'Golden Games', zones: ['L3', 'L4']},
  {id: '3', title: 'Showbies', zones: ['L2-SE']},
]

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
    console.log('found tenant', tenantId, tenant)
    return {
      toObject() {
        return tenant;
      }
    };
  })
}
