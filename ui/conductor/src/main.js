import Vue from 'vue'
import pinia from './plugins/pinia.js';
import router from './routes/router.js';
import vuetify from './plugins/vuetify.js';
import './style.css'
import App from './App.vue'

const app = new Vue({
  pinia,
  router,
  vuetify,
  render: h => h(App)
});
app.$mount('#app');
