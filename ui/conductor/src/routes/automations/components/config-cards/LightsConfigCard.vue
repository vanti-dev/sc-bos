<template>
  <v-card elevation="0" tile>
    <v-list tile class="ma-0 pa-0">
      <v-subheader class="text-title-caps-large text-neutral-lighten-3">Lights</v-subheader>
      <v-list-item>
        <v-list-item-title>Unoccupied timeout</v-list-item-title>
        <v-text-field
            v-model="delayTimeout"
            hide-details
            :disabled="blockActions"
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
import useAuthSetup from '@/composables/useAuthSetup';
import {useSidebarStore} from '@/stores/sidebar';
import {computed} from 'vue';

const {blockActions} = useAuthSetup();

const sidebar = useSidebarStore();

const delayTimeout = computed({
  get() {
    return sidebar.data?.config?.unoccupiedOffDelay ?? '0m';
  },
  set(value) {
    sidebar.data.config.unoccupiedOffDelay = value;
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
