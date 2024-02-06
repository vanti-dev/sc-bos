<template>
  <side-bar>
    <lights-config-card v-if="automationType === 'lights'"/>
    <edit-config-card/>
    <v-card-actions class="justify-end px-4 pt-0" v-if="false">
      <v-btn class="primary" :disabled="blockActions" @click="saveConfig">Save</v-btn>
    </v-card-actions>
    <v-snackbar v-model="saveConfirm" timeout="2000" color="success" max-width="250" min-width="200">
      <span class="text-body-large align-baseline"><v-icon left>mdi-content-save-check</v-icon>Config saved</span>
    </v-snackbar>
  </side-bar>
</template>

<script setup>
import {newActionTracker} from '@/api/resource';
import {configureService, ServiceNames as ServiceTypes} from '@/api/ui/services';
import SideBar from '@/components/SideBar.vue';
import {useErrorStore} from '@/components/ui-error/error';
import useAuthSetup from '@/composables/useAuthSetup';
import LightsConfigCard from '@/routes/automations/components/config-cards/LightsConfigCard.vue';
import EditConfigCard from '@/routes/system/components/service-cards/EditConfigCard.vue';
import {usePageStore} from '@/stores/page';
import {serviceName} from '@/util/proxy';
import {storeToRefs} from 'pinia';
import {computed, onMounted, onUnmounted, reactive, ref} from 'vue';

const {blockActions} = useAuthSetup();

const pageStore = usePageStore();
const {sidebarData, sidebarNode} = storeToRefs(pageStore);

const saveTracker = reactive(/** @type {ActionTracker<Service.AsObject>} */ newActionTracker());
const saveConfirm = ref(false);

const automationType = computed(() => {
  return sidebarData.value?.config?.type ?? '';
});

const node = computed(() => {
  return sidebarNode.value?.name;
});

// UI error handling
const errorStore = useErrorStore();
let unwatchError;
onMounted(() => {
  unwatchError = errorStore.registerTracker(saveTracker);
});
onUnmounted(() => {
  unwatchError();
});


/**
 *
 */
async function saveConfig() {
  const req = {
    name: serviceName(node.value, ServiceTypes.Automations),
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
