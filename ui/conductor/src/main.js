import Vue from 'vue'
import vuetify from './plugins/vuetify.js';
import './style.css'
import App from './App.vue'

const app = new Vue({
  vuetify,
  render: h => h(App)
});
app.$mount('#app');
