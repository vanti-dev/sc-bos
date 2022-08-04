import Vue from 'vue'
import vuetify from './plugins/vuetify.js';
import './style.css'
import App from './App.vue'
import router from './routes/router.js';

const app = new Vue({
  router,
  vuetify,
  render: h => h(App)
});
app.$mount('#app');
