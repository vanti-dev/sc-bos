<template>
  <v-list class="pa-0" dense nav>
    <v-list-item
        v-for="automation of automationTypeList"
        :key="automation.type"
        :to="'/automations/'+automation.type">
      <v-list-item-icon><v-icon>{{ icon[automation.type] }}</v-icon></v-list-item-icon>
      <v-list-item-content class="text-capitalize">{{ automation.type }}</v-list-item-content>
    </v-list-item>
  </v-list>
</template>

<script setup>

import {useAutomationsStore} from '@/routes/automations/store';
import {storeToRefs} from 'pinia';
import {onMounted, ref} from 'vue';

const automationsStore = useAutomationsStore();
const {automationTypeList} = storeToRefs(automationsStore);

// map of icons to use for different automation sections
const icon = ref({
  lights: 'mdi-lightbulb'
});

onMounted(() => automationsStore.refreshMetadata());

</script>

<style scoped>
:deep(.v-list-item--active) {
  color: var(--v-primary-base);
}
</style>
