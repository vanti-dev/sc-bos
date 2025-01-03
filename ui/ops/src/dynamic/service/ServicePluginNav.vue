<template>
  <v-list nav density="compact" class="pa-0">
    <v-progress-linear v-if="showSpinner" indeterminate color="primary"/>
    <v-list-item v-if="props.parent" :to="props.parent" title="Back" prepend-icon="mdi-arrow-left"/>
    <component :is="navSlot" v-if="navSlot"/>
    <v-divider v-if="showSeparator" class="my-2"/>
    <v-list-item :to="toManualEdit" title="Manual Edit" prepend-icon="mdi-pencil" class="my-2">
      <template #append>
        <v-icon v-tooltip="'Danger area! Proceed at your own risk'"
                icon="mdi-alert"
                color="warning"/>
      </template>
    </v-list-item>
  </v-list>
</template>

<script setup>
import {usePullService} from '@/composables/services.js';
import {useServicePlugin} from '@/dynamic/plugins.js';
import {useServiceRouterLink} from '@/dynamic/route.js';
import {serviceName} from '@/util/gateway.js';
import {computed, markRaw, provide} from 'vue';

const props = defineProps({
  name: {
    type: String,
    default: ''
  },
  id: {
    type: String,
    default: ''
  },
  category: {
    type: /** @type {('automations' | 'drivers' | 'zones' | 'systems')} */ String,
    default: ''
  },
  parent: {
    type: /** @type {import('vue-router').RouteLocationRaw} */ [Object, String],
    default: undefined
  }
});

// todo: remove duplicate fetches between this and the ServicePluginParent
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

for (const [key, value] of Object.entries(pullService)) {
  provide(`service.${key}`, value);
}
for (const [key, value] of Object.entries(servicePlugin)) {
  provide(`plugin.${key}`, value);
}

const showSpinner = computed(() => Boolean(!loaded.value || pullService.loading.value));
const navSlot = computed(() => {
  const nav = plugin.value?.slots?.nav;
  return nav ? markRaw(nav()) : null;
});
const showSeparator = computed(() => Boolean(navSlot.value));

const {toManualEdit} = useServiceRouterLink(() => props.category, () => props.name, () => props.id);
</script>

<style scoped>

</style>
