import {defineAsyncComponent} from 'vue';

export const service = /** @type {ServicePlugin} */ {
  metadata: {
    appearance: {
      title: 'Mock Driver',
      description: 'The mock driver provides a way to setup named devices, implementing traits, without any real hardware.'
    }
  },
  slots: {
    nav: () => defineAsyncComponent(() => import('./MockNav.vue'))
  },
  defaultRoute: 'devices',
  routes: [
    {
      path: 'devices',
      name: 'mock-driver-devices',
      component: () => import('./MockDevices.vue')
    }
  ]
};

export const plugin =
    /** @type {Plugin} */
    {service};
