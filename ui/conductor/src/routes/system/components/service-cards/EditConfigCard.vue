<template>
  <v-card flat tile>
    <v-subheader class="text-title-caps-large text-neutral-lighten-3">Config</v-subheader>
    <v-card-text class="px-4 pt-0 json-form">
      <div class="text-caption pb-1 text-neutral text--lighten-6">Export</div>
      <v-textarea
          v-model="config"
          auto-grow
          class="text-body-code"
          :disabled="blockSystemEdit"
          :error-messages="jsonError"
          filled
          full-width
          hide-details="auto"
          readonly/>
      <v-btn
          icon
          tile
          class="copy-btn"
          :disabled="blockSystemEdit"
          @click="copyConfig">
        <v-icon>mdi-content-copy</v-icon>
      </v-btn>
      <v-snackbar v-model="copyConfirm" timeout="2000" color="success" max-width="250" min-width="200">
        <span class="text-body-large align-baseline"><v-icon left>mdi-check-circle</v-icon>Config copied</span>
      </v-snackbar>
    </v-card-text>
  </v-card>
</template>

<script setup>
import useAuthSetup from '@/composables/useAuthSetup';
import {useSidebarStore} from '@/stores/sidebar';
import {computed, ref} from 'vue';

const {blockSystemEdit} = useAuthSetup();

const sidebar = useSidebarStore();

const jsonError = ref('');

const config = computed({
  get() {
    return sidebar.data.service?.configRaw ?? '';
  },
  set(value) {
    jsonError.value = '';
    try {
      sidebar.data.config = JSON.parse(value);
      /**
       * @param {Error} e
       */
    } catch (e) {
      jsonError.value = 'JSON error: ' + e.message;
    }
  }
});

const copyConfirm = ref(false);

/**
 *
 */
function copyConfig() {
  navigator.clipboard.writeText(sidebar.data.service.configRaw);
  copyConfirm.value = true;
}
</script>

<style scoped>
.json-form {
  position: relative;
}

.copy-btn {
  position: absolute;
  top: 28px;
  right: 18px;
}
</style>
