<template>
  <div v-if="!confirmForget">
    <v-form @submit.prevent="onEnroll">
      <v-text-field
          v-model="addressInput"
          class="mx-8"
          :clearable="addressInput !== null"
          dense
          hide-details
          label="Component Address"
          outlined
          @click:clear="addressInput = null"/>
      <!-- Error label if the address is already enrolled -->
      <v-alert
          v-if="errorText"
          class="
        mx-8
        mt-2"
          color="error"
          dense
          max-width="400px"
          outlined
          type="error">
        <v-row class="pa-2 d-flex flex-row flex-nowrap">
          <span class="text-capitalize">{{ errorText.message }}</span>
          <v-spacer/>
          <status-alert
              v-if="errorText?.error"
              :resource="errorText.error"
              icon="mdi-alert-circle-outline"
              class="ml-2"/>
        </v-row>
      </v-alert>

      <v-card-actions class="d-flex flex-row justify-space-around mt-10">
        <v-btn
            class="mr-4 px-4"
            color="primary"
            :disabled="!address || isEnrolled"
            text
            type="submit"
            @click="onEnroll">
          Enroll Node
        </v-btn>
      </v-card-actions>
    </v-form>
  </div>
  <div v-else style="max-width: 600px">
    <v-card-text class="px-7 text-left text-subtitle-1 font-weight-regular">
      <p>
        Forgetting a node means
        <span class="font-weight-bold warning--text">it can no longer interact with other Smart Core nodes,
          and those nodes cannot interact with it.</span>
      </p>
      <p>
        Any automations that rely on inter-node communication with or from this node
        <span class="font-weight-bold error--text">will stop working!</span>
        This includes managing the node centrally via this app. You can re-enrol this node at any time.
      </p>
    </v-card-text>
    <v-card-actions class="d-flex flex-row justify-space-around mt-10">
      <v-btn
          class="mr-4 px-4"
          text
          @click="cancelAction">
        Cancel
      </v-btn>
      <v-btn
          class="px-4"
          color="error"
          text
          @click="forgetHubNode">
        Forget Node
      </v-btn>
    </v-card-actions>
  </div>
</template>

<script setup>
import StatusAlert from '@/components/StatusAlert.vue';
import {computed, ref, watch} from 'vue';
import {formatErrorMessage} from '@/util/error';

const emits = defineEmits([
  'update:dialogState',
  'update:address',
  'inspectHubNodeAction',
  'resetInspectHubNodeValue',
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
  inspectHubNodeValue: {
    type: Object,
    default: () => ({})
  },
  listItems: {
    type: Array,
    default: () => []
  },
  nodeQuery: {
    type: Object,
    default: () => ({})
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

// Enroll the hub node
const onEnroll = () => {
  if (!addressInput.value) {
    return;
  }

  emits('inspectHubNodeAction', addressInput.value);
};

const forgetHubNode = () => {
  if (!props.nodeQuery.address) {
    return;
  }

  emits('forgetHubNodeAction', props.nodeQuery.address);
  emits('update:dialogState', false);
  confirmForget.value = false;
};

const cancelAction = () => {
  return props.nodeQuery.isToForget ? emits('update:dialogState', false) : confirmForget.value = false;
};

// Display the correct dialog content depending on the confirmForget value
// If confirmForget is true, display the forget dialog
// If confirmForget is false, display the enroll/forget dialog
watch(() => props.nodeQuery, (newQuery) => {
  confirmForget.value = newQuery.isToForget;
}, {immediate: true, deep: true});

// Depending on input, check if the address is enrolled
// and update the isEnrolled value to enable the correct button
watch(addressInput, (newAddress, oldAddress) => {
  if (newAddress !== oldAddress) {
    emits('resetInspectHubNodeValue');
  }

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

const errorText = computed(() => {
  if (props.inspectHubNodeValue.error) {
    return {
      message: formatErrorMessage(props.inspectHubNodeValue.error.error.message),
      error: props.inspectHubNodeValue.error
    };
  } else if (!props.inspectHubNodeValue.error && isEnrolled.value) {
    return {
      message: 'This node is already enrolled',
      error: null
    };
  }

  return null;
});
</script>

