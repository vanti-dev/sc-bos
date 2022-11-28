import ChildOnlyPage from "@/components/ChildOnlyPage.vue";
import lighting from "@/routes/operate/systems/lighting/route.js";
import {route } from "@/util/router.js";

export default [
    {
        path: "systems/lighting",
        component: () => import("./lighting/LightingTable.vue"),
    },

]
    