import VueRouter from 'vue-router';
import admin from './admin/route.js';
import commission from './commission/route.js';
import design from './design/route.js';
import operate from './operate/route.js';
import start from './start/route.js';
import Vue from 'vue';

Vue.use(VueRouter);

function route(route) {
  if (Array.isArray(route)) {
    return route;
  }
  return [route];
}

const router = new VueRouter({
  mode: 'history',
  routes: [
    ...route(start),
    ...route(design),
    ...route(commission),
    ...route(operate),
    ...route(admin),
  ]
});

export default router;
