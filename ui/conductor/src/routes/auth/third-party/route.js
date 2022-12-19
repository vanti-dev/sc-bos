export default [
  {path: 'third-party', component: () => import('./AccountList.vue')},
  {path: 'third-party/:tenantId', component: () => import('./Account.vue')}
];
