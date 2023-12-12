export default [
  {
    path: 'overview',
    component: () => import('@/routes/ops/OpsHome.vue'),
    redirect: 'overview/building',
    children: [
      {
        name: 'building-overview',
        path: 'building', // Directly map to the building path
        component: () => import('@/routes/ops/overview/BuildingOverview.vue')
      },
      {
        name: 'dynamic-areas-overview',
        path: 'building/:pathMatch(.*)*', // Captures all segments after /building/
        component: () => import('@/routes/ops/overview/DynamicAreasOverview.vue'),
        props: route => ({pathSegments: route.params.pathMatch.split('/')}) // Splits segments into an array
      }
    ]
  }
];
