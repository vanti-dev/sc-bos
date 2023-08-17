<template>
  <div :class="[`${props.statusBarColor}-status`, 'd-flex flex-row rounded-t']" style="height: 44px">
    <v-card-title class="d-block ext-uppercase text-body-1 font-weight-medium pt-3 text-truncate" style="width: 350px">
      {{ props.device.title }}
    </v-card-title>
    <v-spacer/>
    <!-- <v-card-title class="text-body-1 font-weight-bold text-uppercase white--text pt-3">
    </v-card-title> -->
    <v-tooltip v-if="showClose" bottom>
      <template #activator="{ on, attrs }">
        <v-btn class="elevation-0 mr-1 mt-1" icon v-bind="attrs" v-on="on" @click="showClose = false">
          <v-icon color="white">mdi-close</v-icon>
        </v-btn>
      </template>
      <span>Close</span>
    </v-tooltip>
  </div>
</template>

<script setup>
import {storeToRefs} from 'pinia';
import {useStatusBarStore} from '@/routes/ops/security/components/access-point-card/statusBarStore';

const props = defineProps({
  value: {
    type: Object,
    default: () => {}
  },
  device: {
    type: Object,
    default: () => {}
  },
  statusBarColor: {
    type: String,
    default: ''
  }
});

const statusBarStore = useStatusBarStore();
const {showClose} = storeToRefs(statusBarStore);
</script>
<style lang="scss" scoped>
.granted-status {
  background-color: green;
  transition: background-color 0.5s ease-in-out;
}
.denied-status,
.forced-status,
.failed-status {
  background-color: red;
}
.pending-status,
.aborted-status,
.tailgate-status {
  background-color: orange;
}
.grant_unknown-status {
  background-color: grey;
}

.transparent-status {
  background-color: transparent;
  transition: background-color 0.5s ease-in-out;
}
</style>
