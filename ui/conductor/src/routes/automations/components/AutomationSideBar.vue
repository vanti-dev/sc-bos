<template>
  <side-bar>
    <lights-config-card v-if="automationType === 'lights'"/>
    <edit-config-card/>
    <v-card-actions class="justify-end px-4 pt-0">
      <v-btn class="primary" @click="saveConfig">Save</v-btn>
    </v-card-actions>
    <v-snackbar v-model="saveConfirm" timeout="2000" color="success" max-width="250" min-width="200">
      <span class="text-body-large align-baseline"><v-icon left>mdi-content-save-check</v-icon>Config saved</span>
    </v-snackbar>
  </side-bar>
</template>

<script setup>
import SideBar from '@/components/SideBar.vue';
import EditConfigCard from '@/routes/system/components/service-cards/EditConfigCard.vue';
import LightsConfigCard from '@/routes/automations/components/config-cards/LightsConfigCard.vue';
import {usePageStore} from '@/stores/page';
import {storeToRefs} from 'pinia';
import {computed, reactive, ref} from 'vue';
import {newActionTracker} from '@/api/resource';
import {configureService, ServiceNames as ServiceTypes} from '@/api/ui/services';

const pageStore = usePageStore();
const {sidebarData} = storeToRefs(pageStore);

const saveTracker = reactive(/** @type {ActionTracker<Service.AsObject>} */ newActionTracker());
const saveConfirm = ref(false);

const automationType = computed(() => {
  return sidebarData.value?.config?.type ?? '';
});

/**
 *
 */
async function saveConfig() {
  const req = {
    name: ServiceTypes.Automations,
    id: sidebarData.value.id,
    configRaw: JSON.stringify(sidebarData.value.config, null, 2)
  };

  sidebarData.value = await configureService(req, saveTracker);
  sidebarData.value.config = JSON.parse(sidebarData.value.configRaw ?? {});
  saveConfirm.value = true;
}

</script>

<style scoped>
</style>
