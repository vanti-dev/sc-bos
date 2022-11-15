export default [
  {
    path: "analytics/environment",
    component: () => import("./Environment.vue"),
  },
  {
    path: "analytics/usage",
    component: () => import("./Usage.vue"),
  },
  {
    path: "analytics/occupancy",
    component: () => import("./Occupancy.vue"),
  },
];
