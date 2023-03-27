<template>
  <v-form
      @submit.prevent="saveConfig"
      ref="form">
    <v-card flat tile>
      <v-subheader class="text-title-caps-large neutral--text text--lighten-3">Config</v-subheader>
      <v-card-text class="px-4 pt-0 json-form">
        <div class="text-caption pb-1 neutral--text text--lighten-6">Import/Export</div>
        <v-textarea
            v-model="config"
            class="text-body-code"
            full-width
            filled
            auto-grow
            hide-details="auto"
            :error-messages="jsonError"
            @submit.prevent="saveConfig"/>
        <v-btn icon tile class="copy-btn" @click="copyConfig"><v-icon>mdi-content-copy</v-icon></v-btn>
        <v-snackbar v-model="copyConfirm" timeout="2000" color="success" max-width="250" min-width="200">
          <span class="text-body-large align-baseline"><v-icon left>mdi-check-circle</v-icon>Config copied</span>
        </v-snackbar>
        <v-snackbar v-model="saveConfirm" timeout="2000" color="success" max-width="250" min-width="200">
          <span class="text-body-large align-baseline"><v-icon left>mdi-content-save-check</v-icon>Config saved</span>
        </v-snackbar>
      </v-card-text>
      <v-card-actions class="justify-end px-4 pt-0">
        <v-btn class="primary" type="submit">Save</v-btn>
      </v-card-actions>
    </v-card>
  </v-form>
</template>

<script setup>
import {usePageStore} from '@/stores/page';
import {storeToRefs} from 'pinia';
import {computed, reactive, ref} from 'vue';
import {configureService, ServiceNames as ServiceTypes} from '@/api/ui/services';
import {newActionTracker} from '@/api/resource';

const pageStore = usePageStore();
const {sidebarData} = storeToRefs(pageStore);

const jsonError = ref('');

const config = computed({
  get() {
    return sidebarData.value.configRaw;
  },
  set(value) {
    jsonError.value = '';
    try {
      sidebarData.value.config = JSON.parse(value);
      /**
       * @param {Error} e
       */
    } catch (e) {
      jsonError.value = 'JSON error: '+e.message;
    }
  }
});

const saveTracker = reactive(/** @type {ActionTracker<Service.AsObject>} */ newActionTracker());
const saveConfirm = ref(false);

/**
 *
 */
async function saveConfig() {
  const req = {
    name: ServiceTypes.Automations,
    id: sidebarData.value.id,
    configRaw: JSON.stringify(sidebarData.value.config, null, 2)
  };

  await configureService(req, saveTracker);
  sidebarData.value.configRaw = req.configRaw;
  saveConfirm.value = true;
}

const copyConfirm = ref(false);
/**
 *
 */
function copyConfig() {
  navigator.clipboard.writeText(sidebarData.value.configRaw);
  copyConfirm.value = true;
}

</script>

<style scoped>
.json-form{
  position: relative;
}
.copy-btn {
  position: absolute;
  top: 28px;
  right: 18px;
}
</style>
