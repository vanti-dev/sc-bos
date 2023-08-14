<template>
  <v-card elevation="0" tile>
    <v-list tile class="ma-0 pa-0">
      <v-subheader class="text-title-caps-large neutral--text text--lighten-3">Lights</v-subheader>
      <v-list-item>
        <v-list-item-title>Unoccupied timeout</v-list-item-title>
        <v-text-field
            v-model="delayTimeout"
            hide-details
            dense
            outlined
            :rules="[delayRule]"
            readonly
            style="width: 100px;"/>
        <!-- todo: display error message -->
      </v-list-item>
    </v-list>
  </v-card>
</template>

<script setup>
import {computed} from 'vue';
import {usePageStore} from '@/stores/page';
import {storeToRefs} from 'pinia';

const pageStore = usePageStore();
const {sidebarData} = storeToRefs(pageStore);

const delayTimeout = computed({
  get() {
    return sidebarData.value?.config?.unoccupiedOffDelay ?? '0m';
  },
  set(value) {
    sidebarData.value.config.unoccupiedOffDelay = value;
  }
});


/**
 * @param {string} value
 * @return {boolean|string}
 */
function delayRule(value) {
  const pattern = /^\d*[smh]$/;
  return pattern.test(value) || 'Please specify: #[s|m|h] (e.g. 20m)';
}

</script>

<style>
</style>
