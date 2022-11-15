export default [
  {
    path: "reports/environment",
    component: () => import("./Environment.vue"),
  },
  {
    path: "reports/usage",
    component: () => import("./Usage.vue"),
  },
  {
    path: "reports/occupancy",
    component: () => import("./Occupancy.vue"),
  },
];
