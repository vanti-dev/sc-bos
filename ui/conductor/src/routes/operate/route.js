import ChildOnlyPageV from "@/components/ChildOnlyPage.vue";
import lighting from "@/routes/operate/systems/lighting/route.js";
import notifications from "@/routes/operate/status/notifications/route.js";
import reports from "@/routes/operate/reports/route.js";
import analytics from "@/routes/operate/analytics/route.js";

import { route } from "@/util/router.js";

export default {
  name: "operate",
  path: "/operate",
  redirect: "/operate/summary",
  components: {
    default: ChildOnlyPageV,
    nav: () => import("./OperateNav.vue"),
  },
  children: [
    { path: "summary", component: () => import("./OperateSummary.vue") },
    ...route(lighting),
    ...route(notifications),
    ...route(reports),
    ...route(analytics),
  ],
  meta: {
    title: "Operator",
  },
};
