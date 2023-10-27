<template>
  <v-card class="pa-0" flat tile>
    <v-subheader class="text-title-caps-large neutral--text text--lighten-3">
      Zones
      <v-spacer/>
      <v-btn v-if="hasZoneChanges" color="primary" :disabled="blockActions" @click="emitZoneChanges">Save</v-btn>
    </v-subheader>
    <v-card-actions>
      <AccountZoneChooser :zones.sync="chooserZones" :max-result-size.sync="maxResultSize" class="mx-2 flex-grow-1"/>
    </v-card-actions>
    <v-card-text v-if="chooserZones.length === 0" class="pt-1 font-italic">
      {{ hasZoneChanges ? 'If there are no' : 'No' }}
      zones associated with this account, they will not be able to access any devices.
    </v-card-text>
  </v-card>
</template>

<script setup>
import {computed, ref} from 'vue';
import useAuthSetup from '@/composables/useAuthSetup';

import AccountZoneChooser from '@/routes/auth/third-party/components/AccountZoneChooser.vue';

const props = defineProps({
  zoneList: {
    type: Array, // string[]
    default: () => []
  }
});

const emit = defineEmits(['update:zone-list']);

const zoneChanges = ref(null);
const maxResultSize = ref(5);
const chooserZones = computed({
  get() {
    if (zoneChanges.value !== null) {
      return zoneChanges.value;
    } else {
      return props.zoneList;
    }
  },
  set(v) {
    zoneChanges.value = v;
  }
});

const hasZoneChanges = computed(() => {
  if (zoneChanges.value === null) {
    return false;
  }
  if (zoneChanges.value.length !== props.zoneList.length) {
    return true;
  }

  const zoneNames = props.zoneList.reduce((acc, zone) => {
    acc[zone] = true;
    return acc;
  }, {});
  return zoneChanges.value?.some(zone => !zoneNames[zone.name]) ?? false;
});

const emitZoneChanges = () => {
  emit('update:zone-list', zoneChanges.value.map(zone => zone.name));
};

// ------------------------------ //
// ----- Authentication settings ----- //

const {blockActions} = useAuthSetup();
</script>

<style scoped></style>
