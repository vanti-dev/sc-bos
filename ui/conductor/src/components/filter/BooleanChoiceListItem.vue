<template>
  <v-list-item @click="_value = !_value">
    <v-list-item-icon>
      <v-icon>{{ props.icon }}</v-icon>
    </v-list-item-icon>
    <!-- all this flexing is to stop some strange layout issues caused by the default flex: row of content -->
    <v-list-item-content class="flex-column align-start">
      <template v-if="choiceText">
        <v-list-item-subtitle class="text-overline">{{ props.title }}</v-list-item-subtitle>
        <v-list-item-title class="align-self-auto">{{ choiceText }}</v-list-item-title>
      </template>
      <v-list-item-title class="align-self-auto" v-else>{{ props.title }}</v-list-item-title>
    </v-list-item-content>
    <v-list-item-action>
      <v-switch2 :indeterminate="_indeterminate" v-model="_value" style="pointer-events: none"/>
    </v-list-item-action>
    <v-list-item-action>
      <v-btn icon @click.stop="emits('clear')" class="mx-n2" :disabled="isDefault">
        <v-icon size="20">mdi-close</v-icon>
      </v-btn>
    </v-list-item-action>
  </v-list-item>
</template>

<script setup>
import VSwitch2 from '@/components/VSwitch2.vue';
import {computed, ref, watch} from 'vue';

const props = defineProps({
  icon: {
    type: String,
    default: ''
  },
  title: {
    type: String,
    default: ''
  },
  choice: {
    type: Object, // Choice
    default: null
  },
  defaultChoice: {
    type: Boolean, // Choice
    default: false
  }
});
const emits = defineEmits(['input', 'clear']);

const _out = ref(/** @type {boolean | null} */ null);
watch(_out, (newValue, oldValue) => {
  if (newValue === oldValue) return;
  emits('input', newValue);
});
watch(() => props.choice?.value,
    (newValue) => _out.value = newValue,
    {immediate: true});

const _indeterminate = computed(() => _out.value === null || _out.value === undefined);
const _value = computed({
  get: () => {
    if (_out.value === null) return false;
    return _out.value;
  },
  set: (value) => {
    // instead of cycling from false to true, cycle from false to unset.
    if (value && _out.value === false) {
      _out.value = null;
    } else {
      _out.value = value;
    }
  }
});

const isDefault = computed(() => props.defaultChoice);

const choiceText = computed(() => props.choice?.text);
</script>
