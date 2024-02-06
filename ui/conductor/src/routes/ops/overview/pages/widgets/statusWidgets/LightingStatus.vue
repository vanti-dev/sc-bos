<template>
  <v-btn class="ma-0 pa-0 mt-4" text :to="navigateToDevices" :ripple="false">
    <content-card class="mt-0 pa-6">
      <h4 class="text-h4">
        Lighting
      </h4>
      <v-divider :class="[lightStatusBar, 'py-1 rounded mt-6 mb-n6']"/>
    </content-card>
  </v-btn>
</template>

<script setup>
import ContentCard from '@/components/ContentCard.vue';
import useLighting from '@/composables/traits/useLighting';
import useTraits from '@/composables/useTraits';
import {computed, onUnmounted, reactive, ref, watch} from 'vue';

const props = defineProps({
  floor: {
    type: String,
    default: undefined
  },
  item: {
    type: [Array, Object],
    default: () => ([])
  },
  zone: {
    type: String,
    default: undefined
  }
});

const {collectDeviceNames} = useTraits(props);
const lightingNames = computed(() => collectDeviceNames('Light'));
const lightingInstances = ref({});
const lightStatus = reactive({});

const lightStatusBar = computed(() => {
  const lightStatusValues = Object.values(lightStatus);

  const success = lightStatusValues.every((value) => value === false);
  const warning = lightStatusValues.some((value) => value === false);
  const error = lightStatusValues.every((value) => value === true);

  // If there are lights, and all are true, return success
  if (success) return 'success';

  // If there are lights, but not all are true, return warning
  if (warning) return 'warning';

  // If there are lights, and all are false, return error
  if (error) return 'error';

  return 'neutral lighten-2';
});

const navigateToDevices = computed(() => {
  if (props.floor && !props.zone) {
    return '/devices/lighting/floors/' + encodeURIComponent(props.floor);
  } else if (!props.floor && props.zone) {
    return '/devices/lighting/zones/' + encodeURIComponent(props.zone);
  } else if (props.floor && props.zone) {
    return '/devices/lighting/floors/' + encodeURIComponent(props.floor) + '/zones/' + encodeURIComponent(props.zone);
  }

  return '/devices/lighting';
});

// Watch for changes to the lighting names and update the lightStatus object
watch(lightingNames, (newNames, oldNames) => {
  // Cleanup old names
  if (oldNames) {
    oldNames.forEach(name => {
      delete lightStatus[name];
      delete lightingInstances.value[name]; // Cleanup old instances
    });
  }

  // Handle new names
  if (newNames) {
    newNames.forEach(name => {
      lightingInstances.value[name] = useLighting({name, paused: false});

      // Watching each lighting device for changes in the lightValue
      watch(() => lightingInstances.value[name].lightValue, (newValue) => {
        lightStatus[name] = !!newValue?.streamError;
      }, {immediate: true});
    });
  }
}, {immediate: true, deep: true});

onUnmounted(() => {
  // Cleanup
  lightingNames.value.forEach(name => {
    lightingInstances.value[name].cleanUp();
    delete lightingInstances.value[name];
    delete lightStatus[name];
  });
});
</script>

<style lang="scss" scoped>
@keyframes pulse {
  0% {
    opacity: 0.3;
  }

  50% {
    opacity: 1;
  }

  100% {
    opacity: 0.3;
  }
}

/** Pulse animation for the light status bar if errored */
.v-divider.error {
  animation: pulse 1s infinite;
}
</style>
