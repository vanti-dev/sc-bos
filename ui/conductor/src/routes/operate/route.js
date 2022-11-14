import ChildOnlyPageV from "@/components/ChildOnlyPage.vue";
import lighting from "@/routes/operate/lighting/route.js";
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
  ],
  meta: {
    title: "Operator",
  },
};
