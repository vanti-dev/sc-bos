import App from "@/App.vue";
import pinia from "@/plugins/pinia.js";
import vuetify from "@/plugins/vuetify.js";
import router from "@/routes/router.js";
import "@/style.css";
import Vue from "vue";

const app = new Vue({
  pinia,
  router,
  vuetify,
  render: (h) => h(App),
});
app.$mount("#app");
