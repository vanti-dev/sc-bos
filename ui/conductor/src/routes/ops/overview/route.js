export default [
  {
    path: 'overview',
    component: () => import('@/routes/ops/OpsHome.vue'),
    children: [
      {
        name: 'building-overview',
        path: 'building', // Directly map to the building path
        component: () => import('@/routes/ops/overview/BuildingOverview.vue')
      },
      {
        name: 'areas-overview',
        path: 'building/areas/:area', // Update the path to include the building segment
        component: () => import('@/routes/ops/overview/AreaOverview.vue'),
        props: true
      },
      {
        name: 'floors-overview',
        path: 'building/floors/:floor', // Update the path to include the building floor segment
        component: () => import('@/routes/ops/overview/FloorOverview.vue'),
        props: true
      },
      {
        name: 'floors-zone-overview',
        path: 'building/floors/:floor/zones/:zone', // Update the path to include the building floor zone segment
        component: () => import('@/routes/ops/overview/FloorZoneOverview.vue'),
        props: true
      }
    ]
  }
];
