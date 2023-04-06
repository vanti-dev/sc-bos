import App from '@/App.vue';
import pinia from '@/plugins/pinia.js';
import vuetify from '@/plugins/vuetify.js';
import router from '@/routes/router.js';
import '@/style.scss';
import Vue from 'vue';
// plugin components
import VueApexCharts from 'vue-apexcharts';

const app = new Vue({
  pinia,
  router,
  vuetify,
  render: (h) => h(App)
});
app.$mount('#app');

Vue.use(VueApexCharts);
Vue.component('Apexchart', VueApexCharts);
