<template>
  <div v-if="!confirmForget">
    <v-combobox
        class="mx-8"
        :clearable="address !== null"
        dense
        hide-details
        :items="props.listItems"
        item-text="address"
        item-value="address"
        label="Component Address"
        outlined
        :search-input.sync="addressInput"/>

    <v-card-actions class="d-flex flex-row justify-space-around mt-10">
      <v-btn
          class="mr-4 px-4"
          color="primary"
          :disabled="!address || isEnrolled"
          text
          @click="emits('inspectHubNodeAction', props.address)">
        Enroll
      </v-btn>
      <v-btn
          class="px-4"
          color="error"
          :disabled="!address || !isEnrolled"
          text
          @click="confirmForget = true">
        Forget
      </v-btn>
    </v-card-actions>
  </div>
  <div v-else>
    <v-card-title class="text-center text-h5 font-weight-bold">
      Are you sure you want to forget this component?
    </v-card-title>
    <v-card-text class="text-center text-body-1">
      You have to re-enroll the component to use it again.
    </v-card-text>
    <v-card-actions class="d-flex flex-row justify-space-around mt-10">
      <v-btn
          class="mr-4 px-4"
          color="primary"
          text
          @click="confirmForget = false">
        Cancel
      </v-btn>
      <v-btn
          class="px-4"
          color="error"
          text
          @click="forgetHubNode">
        Confirm
      </v-btn>
    </v-card-actions>
  </div>
</template>

<script setup>
import {computed, ref, watch} from 'vue';

const emits = defineEmits([
  'update:dialogState',
  'update:address',
  'inspectHubNodeAction',
  'forgetHubNodeAction'
]);
const props = defineProps({
  address: {
    type: String,
    default: null
  },
  dialogState: {
    type: Boolean,
    default: false
  },
  listItems: {
    type: Array,
    default: () => []
  }
});

const isEnrolled = ref(null);
const confirmForget = ref(false);
const addressInput = computed({
  get() {
    return props.address;
  },
  set(value) {
    emits('update:address', value);
  }
});

const forgetHubNode = () => {
  emits('forgetHubNodeAction', addressInput.value);
  emits('update:dialogState', false);
  confirmForget.value = false;
};

// Depending on input, check if the address is enrolled
// and update the isEnrolled value to enable the correct button
watch(addressInput, (newAddress) => {
  // If the address is empty, reset the isEnrolled value to disable all buttons
  if (!newAddress) {
    isEnrolled.value = null;
    return;
  }

  // If the address is not empty, check if it is enrolled
  const matchAddress = props.listItems.find(node => node.address === newAddress);

  // If the address is enrolled, enable the forget button
  if (matchAddress) {
    isEnrolled.value = true;

    // If the address is not enrolled, enable the enroll button
  } else if (!matchAddress) {
    isEnrolled.value = false;

    // Otherwise, disable all buttons
  } else {
    isEnrolled.value = null;
  }
}, {immediate: true, deep: true});

// Reset the address value when the dialog is closed
watch(() => props.dialogState, (newState) => {
  if (!newState) {
    addressInput.value = null;
  }
}, {immediate: true, deep: true});
</script>

