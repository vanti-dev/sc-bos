<template>
  <component :is="component" v-bind="attrs"/>
</template>

<script setup>
import ChildOnlyPage from '@/components/pages/ChildOnlyPage.vue';
import SidebarPage from '@/components/pages/SidebarPage.vue';
import {usePullService} from '@/composables/services.js';
import {useServicePlugin} from '@/dynamic/plugins.js';
import {usePluginRedirect, usePluginRoutes} from '@/dynamic/route.js';
import {serviceName} from '@/util/gateway.js';
import {computed, markRaw, provide} from 'vue';
import {VSkeletonLoader} from 'vuetify/components';

const props = defineProps({
  name: {
    type: String,
    default: ''
  },
  id: {
    type: String,
    required: true
  },
  category: {
    type: /** @type {('automations' | 'drivers' | 'zones' | 'systems')} */ String,
    required: true
  }
});

// todo: remove duplicate fetches between this and the ServicePluginNav
const pullService = usePullService(() => ({
  id: props.id,
  name: serviceName(props.name, props.category)
}));
const {value} = pullService;
const servicePlugin = useServicePlugin(() => ({
  name: props.name,
  type: value.value?.type,
  category: props.category
}));
const {plugin, loaded} = servicePlugin;

usePluginRoutes(plugin);
usePluginRedirect(plugin, () => Boolean(value.value && loaded.value));

provide('service.name', computed(() => serviceName(props.name, props.category)));
for (const [key, value] of Object.entries(pullService)) {
  provide(`service.${key}`, value);
}
for (const [key, value] of Object.entries(servicePlugin)) {
  provide(`plugin.${key}`, value);
}

const component = computed(() => {
  if (!loaded.value) return markRaw(VSkeletonLoader);
  if (!plugin.value) return markRaw(ChildOnlyPage);
  if (plugin.value.slots?.sidebar) return markRaw(SidebarPage);
  return markRaw(ChildOnlyPage);
});
const attrs = computed(() => {
  if (!loaded.value) return {'type': 'article'};
  return {};
});
</script>
