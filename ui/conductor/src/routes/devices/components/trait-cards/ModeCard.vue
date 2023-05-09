<template>
  <v-form @submit.prevent="saveModeValues">
    <v-card elevation="0" tile>
      <v-list tile class="ma-0 pa-0">
        <v-subheader class="text-title-caps-large neutral--text text--lighten-3">Modes</v-subheader>
        <v-list-item v-for="mode in modesDisplay" :key="mode.key">
          <v-list-item-content>
            <template v-if="mode.values">
              <v-select
                  :label="mode.title"
                  :items="mode.values"
                  :value="mode.value"
                  @input="updateMode(mode.key, $event, true)"
                  :disabled="loading"
                  outlined
                  dense
                  hide-details/>
            </template>
            <template v-else>
              <v-text-field
                  :label="mode.title"
                  :value="mode.value"
                  @input="updateMode(mode.key, $event)"
                  :disabled="loading"
                  outlined
                  dense
                  hide-details/>
            </template>
          </v-list-item-content>
        </v-list-item>
      </v-list>
      <v-card-actions class="px-4">
        <v-spacer/>
        <v-btn text type="submit" @click="saveModeValues" :disabled="updateValue.loading || !dirty">Save</v-btn>
      </v-card-actions>
      <v-progress-linear color="primary" indeterminate :active="updateValue.loading"/>
    </v-card>
  </v-form>
</template>

<script setup>

import {closeResource, newActionTracker, newResourceValue} from '@/api/resource';
import {describeModes, pullModeValues, updateModeValues} from '@/api/sc/traits/mode';
import {useErrorStore} from '@/components/ui-error/error';
import {computed, onMounted, onUnmounted, reactive, set, watch} from 'vue';

const props = defineProps({
  // unique name of the device
  name: {
    type: String,
    default: ''
  }
});

const modeValue = reactive(newResourceValue());
const updateValue = reactive(newActionTracker());
const modeInfo = reactive(newActionTracker());

const loading = computed(() => modeValue.loading || updateValue.loading || describeModes.loading);

// track edits so we can commit them all in one go, or show a dirty flag if needed
const edits = reactive({});
const updateData = computed(() => modes.value.map(([k, v]) => [k, edits[k] ?? v]));
const dirty = computed(() => {
  for (let i = 0; i < modes.value.length; i++) {
    const old = modes.value[i];
    const edit = updateData.value[i];
    if (old[0] !== edit[0] || old[1] !== edit[1]) {
      return true;
    }
  }
  return false;
});

// UI error handling
const errorStore = useErrorStore();
let unwatchModeError;
let unwatchUpdateError;
let unwatchDescribeError;
onMounted(() => {
  unwatchModeError = errorStore.registerValue(modeValue);
  unwatchUpdateError = errorStore.registerTracker(updateValue);
  unwatchDescribeError = errorStore.registerTracker(modeInfo);
});
onUnmounted(() => {
  if (unwatchModeError) unwatchModeError();
  if (unwatchUpdateError) unwatchUpdateError();
  if (unwatchDescribeError) unwatchDescribeError();
});

// if device name changes
watch(() => props.name, async (name) => {
  // close existing stream if present
  closeResource(modeValue);
  // create new stream
  if (name && name !== '') {
    // noinspection ES6MissingAwait - handled by tracker
    describeModes({name}, modeInfo);
    pullModeValues(name, modeValue);
  }
}, {immediate: true});

onUnmounted(() => {
  closeResource(modeValue);
});

const modes = computed(() => {
  if (modeValue && modeValue.value) {
    return modeValue.value.valuesMap;
  }
  return [];
});
const modesDisplay = computed(() => modes.value.map(m => modeDisplay(m)));

// like {"myMode": ["opt1", "opt2"]}
const modeInfoMap = computed(() => {
  const modesList = modeInfo.response?.availableModes?.modesList || [];
  const res = {};
  for (const mode of modesList) {
    res[mode.name] = mode.valuesList.map(v => v.name);
  }
  return res;
});

/**
 * @param {string[]} mode
 * @return {{key: string, value: string, values?: string[], title: string}}
 */
function modeDisplay([k, v]) {
  let items = undefined;
  if (modeInfoMap.value[k]) {
    items = modeInfoMap.value[k];
  }
  if (k === 'lighting.mode') {
    return {
      key: k,
      value: v,
      title: 'Lighting Mode',
      values: items || ['auto', 'normal', 'extended', 'night', 'maintenance', 'test']
    };
  } else {
    return {
      key: k,
      value: v,
      title: k,
      values: items || undefined
    };
  }
}

/**
 * @param {string} key
 * @param {string} value
 * @param {boolean} commit
 */
function updateMode(key, value, commit = false) {
  set(edits, key, value);
  if (commit) {
    saveModeValues();
  }
}

/**
 */
function saveModeValues() {
  /* @type {UpdateModeValuesRequest.AsObject} */
  const req = {
    name: props.name,
    modeValues: {
      valuesMap: modes.value.map(([k, v]) => [k, edits[k] ?? v])
    }
  };
  updateModeValues(req, updateValue);
}

</script>

<style scoped>
.v-list-item {
    min-height: auto;
}

.v-progress-linear {
    width: auto;
}
</style>
