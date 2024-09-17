<template>
  <v-main>
    <v-container min-height="100%">
      <v-card color="transparent" :loading="loading">
        <v-card-title class="d-flex align-center">
          <span class="mr-auto">
            {{ blockActions ? 'Viewing' : 'Editing' }} "{{ serviceName }}" of type {{ serviceType }}
          </span>
          <v-fade-transition>
            <v-btn v-if="refreshVisible" color="warning" @click="doRefresh"
                   variant="text" class="mr-4"
                   v-tooltip="'New remote changes. Refreshing will undo any local changes!'">
              Refresh
            </v-btn>
          </v-fade-transition>
          <span v-tooltip="saveTooltipText">
            <v-btn color="primary" :disabled="saveDisabled" @click="doSave" :loading="saveTracker.loading">Save</v-btn>
          </span>
        </v-card-title>
        <v-expand-transition>
          <v-alert v-if="alertVisible" :type="alertType" :text="alertText"/>
        </v-expand-transition>
        <v-card-text>
          <json-editor v-model="editedJsonModel" :read-only="blockActions"/>
        </v-card-text>
      </v-card>
    </v-container>
  </v-main>
</template>

<script setup>
import {newActionTracker} from '@/api/resource.js';
import {configureService} from '@/api/ui/services.js';
import JsonEditor from '@/components/JsonEditor.vue';
import {statusCodeToString} from '@/components/ui-error/util.js';
import {usePullService} from '@/composables/services.js';
import useAuthSetup from '@/composables/useAuthSetup.js';
import {formatErrorMessage} from '@/util/error.js';
import deepEqual from 'fast-deep-equal';
import {computed, reactive, ref, watch} from 'vue';

const props = defineProps({
  // the device name hosting the service
  name: {
    type: String,
    default: ''
  },
  id: {
    type: String,
    required: true
  }
});

const {blockActions} = useAuthSetup();

// local edit state/props
const editedJsonModel = ref(/** @type {string|any} */ '');
const editedJsonParsed = computed(() => {
  if (typeof editedJsonModel.value === 'object') return {json: editedJsonModel.value};
  try {
    return {json: JSON.parse(editedJsonModel.value)};
  } catch (e) {
    return {error: e};
  }
});
const editedJsonUnfilteredRaw = computed(() => {
  const model = typeof editedJsonModel.value === 'object' ? editedJsonModel.value : JSON.parse(editedJsonModel.value);
  model.name = serviceName.value;
  model.type = serviceType.value;
  model.disabled = serviceDisabled.value;
  return JSON.stringify(model);
});
const editedIsChanged = computed(() => !deepEqual(editedJsonParsed.value?.json, filteredJson.value));

// remote config state/props
const pullServiceRequest = computed(() => {
  const req = {id: props.id};
  if (props.name) {
    req.name = props.name;
  }
  return req;
});
const {value, streamError, loading} = usePullService(pullServiceRequest);
const fetchedJsonRaw = computed(() => /** @type {string|undefined} */ value.value?.configRaw);
const fetchedJson = computed(() => JSON.parse(fetchedJsonRaw.value || '{}'));
const filteredJson = computed(() => {
  const unwanted = ['disabled', 'name', 'type'];
  return Object.fromEntries(Object.entries(fetchedJson.value).filter(([key]) => !unwanted.includes(key)));
});
watch(filteredJson, (newValue, oldValue) => {
  if (deepEqual(newValue, oldValue)) return;
  if (!oldValue || !editedJsonModel.value || deepEqual(oldValue, editedJsonParsed.value?.json)) {
    editedJsonModel.value = newValue;
    return;
  }
  refreshVisible.value = !deepEqual(newValue, editedJsonParsed.value?.json);
});

// refresh and save
const refreshVisible = ref(false);
const doRefresh = () => {
  editedJsonModel.value = filteredJson.value;
  refreshVisible.value = false;
};

const saveDisabled = computed(() => Boolean(!editedIsChanged.value || editedJsonParsed.value.error || refreshVisible.value || blockActions.value));
const saveTooltipText = computed(() => {
  if (blockActions.value) return 'You do not have permission to change the configuration';
  if (editedJsonParsed.value.error) return 'You have errors in your JSON configuration';
  if (!editedIsChanged.value) return 'No changes made';
  if (refreshVisible.value) return 'The remote configuration has changed. Refresh and reapply your changes to save.';
  return 'Save changes';
});
const saveTracker = reactive(newActionTracker());
const doSave = () => {
  const request = {
    name: props.name,
    id: props.id,
    configRaw: editedJsonUnfilteredRaw.value
  };
  configureService(request, saveTracker);
};

// display utils
const serviceName = computed(() => fetchedJson.value.name ?? props.id);
const serviceType = computed(() => fetchedJson.value.type ?? 'unknown');
const serviceDisabled = computed(() => fetchedJson.value.disabled);

// error display
const alertVisible = computed(() => streamError.value || saveTracker.error);
const alertType = computed(() => 'error');
const alertText = computed(() => {
  if (!alertVisible.value) return '';
  if (saveTracker.error) return errorString(saveTracker.error, 'Error saving service configuration');
  return errorString(streamError.value, 'Error fetching service information');
});
const errorString = (error, prefix = 'An error occurred') => {
  if (!error) return '';
  if (error.error) error = error.error;
  const parts = [prefix + ':'];
  if (error.code) parts.push(statusCodeToString(error.code));
  if (error.message) parts.push(formatErrorMessage(error.message));
  return parts.join(' ');
};
</script>
