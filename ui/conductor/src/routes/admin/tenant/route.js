export default [
  {path: 'tenants', component: () => import('./TenantList.vue')},
  {path: 'tenants/:tenantId', component: () => import('./Tenant.vue')}
]
