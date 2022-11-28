import ChildOnlyPage from "@/components/ChildOnlyPage.vue";
import systems from "@/routes/operate/systems/route.js";
import notifications from "@/routes/operate/status/notifications/route.js";
import reports from "@/routes/operate/reports/route.js";
import analytics from "@/routes/operate/analytics/route.js";

import { route } from "@/util/router.js";

export default {
  name: "operate",
  path: "/operate",
  redirect: "/operate/summary",
  components: {
    default: ChildOnlyPage,
    nav: () => import("./OperateNav.vue"),
  },
  children: [
    { path: "summary", component: () => import("./OperateSummary.vue") },
    ...route(systems),
    ...route(notifications),
    ...route(reports),
    ...route(analytics),
  ],
  meta: {
    title: "Operator",
  },
};
