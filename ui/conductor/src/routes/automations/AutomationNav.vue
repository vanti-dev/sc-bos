<template>
  <v-list class="pa-0" density="compact" nav>
    <v-list-item :disabled="hasNoAccess('/automations/all')" to="/automations/all" prepend-icon="mdi-view-list">
      <v-list-item-title class="text-capitalize">All</v-list-item-title>
    </v-list-item>
    <v-list-item
        v-for="automation of automationTypeList"
        :key="automation.type"
        :to="'/automations/' + encodeURIComponent(automation.type)"
        class="my-2"
        :disabled="hasNoAccess('/automations/' + automation.type)"
        :prepend-icon="icon[mapIconKey(automation.type)] ?? defaultIcon">
      <v-list-item-title class="text-capitalize text-truncate">
        {{ formatNaming(automation.type) }}
      </v-list-item-title>
    </v-list-item>
  </v-list>
</template>

<script setup>
import {ServiceNames} from '@/api/ui/services';
import {usePullServiceMetadata} from '@/composables/services.js';
import useAuthSetup from '@/composables/useAuthSetup';
import {useUserConfig} from '@/stores/userConfig.js';
import {computed, ref} from 'vue';

const {hasNoAccess} = useAuthSetup();

const userConfig = useUserConfig();
const serviceName = computed(() => userConfig.node?.name + '/' + ServiceNames.Automations);
const {value: serviceMd} = usePullServiceMetadata(serviceName);

// filter out automations that have no instances, and map to {type, number} obj
const automationTypeList = computed(() => {
  if (!serviceMd.value) return [];
  const list = [];
  serviceMd.value.typeCountsMap.forEach(([type, number]) => {
    if (number > 0) {
      list.push({type, number});
    }
  });
  list.sort();
  return list;
});

// map of icons to use for different automation sections
const icon = ref({
  bms: 'mdi-office-building-cog-outline',
  history: 'mdi-history',
  lightreport: 'mdi-file-chart-outline',
  lights: 'mdi-lightbulb',
  resetenterleave: 'mdi-account-multiple-remove-outline',
  statusalerts: 'mdi-alert-circle-outline',
  statusemail: 'mdi-email-newsletter',
  udmi: 'mdi-transit-connection-variant'
});
const defaultIcon = 'mdi-auto-mode';

const acronyms = ['bms', 'udmi'];
const suffixes = ['report', 'reports', 'alert', 'alerts', 'email', 'emails', 'reset', 'enter', 'leave'];

/**
 * @param {string} name
 * @return {string}
 */
const formatNaming = (name) => {
  // Split the name by "/" and format each part separately
  const parts = name.split('/').map(part => {
    // Function to dynamically split concatenated suffixes
    const splitConcatenatedWords = (str) => {
      for (const word of suffixes) {
        const index = str.lastIndexOf(word);
        if (index > 0 && index + word.length === str.length) {
          return splitConcatenatedWords(str.substring(0, index)) + ' ' + word;
        }
      }
      return str;
    };

    part = splitConcatenatedWords(part);

    for (const acronym of acronyms) {
      if (part === acronym) {
        return part.toUpperCase();
      }
    }

    return part;
  });

  // Join the formatted parts back together with "/"
  return parts.join('/');
};


/**
 * @param {string} name
 * @return {string}
 */
const mapIconKey = (name) => {
  // Check if the name includes certain keywords or phrases and map accordingly
  if (name.includes('lightreport')) {
    return 'lightreport';
  }
  // Add more checks and mappings as needed

  // Default mapping (no changes)
  return name;
};
</script>

<style scoped>
:deep(.v-list-item--active) {
  color: rgb(var(--v-theme-primary));
}
</style>
