<template>
  <v-card elevation="0" tile>
    <v-list tile class="ma-0 pa-0">
      <v-subheader class="text-title-caps-large neutral--text text--lighten-3">Electric</v-subheader>
      <v-list-item class="py-1" v-if="realPowerStr">
        <v-list-item-title class="text-body-small text-capitalize">Real Power</v-list-item-title>
        <v-list-item-subtitle class="text-end">{{ realPowerStr }}</v-list-item-subtitle>
      </v-list-item>
      <v-list-item class="py-1" v-if="apparentPowerStr">
        <v-list-item-title class="text-body-small text-capitalize">Apparent Power</v-list-item-title>
        <v-list-item-subtitle class="text-end">{{ apparentPowerStr }}</v-list-item-subtitle>
      </v-list-item>
      <v-list-item class="py-1" v-if="reactivePowerStr">
        <v-list-item-title class="text-body-small text-capitalize">Reactive Power</v-list-item-title>
        <v-list-item-subtitle class="text-end">{{ reactivePowerStr }}</v-list-item-subtitle>
      </v-list-item>
      <v-list-item class="py-1" v-if="powerFactorStr">
        <v-list-item-title class="text-body-small text-capitalize">Power Factor</v-list-item-title>
        <v-list-item-subtitle class="text-end">
          {{ powerFactorStr }}
          <!-- for alignment purposes -->
          <span style="visibility: hidden"> kW</span>
        </v-list-item-subtitle>
      </v-list-item>
    </v-list>
    <v-progress-linear color="primary" indeterminate :active="props.loading"/>
  </v-card>
</template>
<script setup>
import {computed} from 'vue';

const props = defineProps({
  value: {
    type: Object, // of ElectricDemand.AsObject
    default: () => {
    }
  },
  loading: {
    type: Boolean,
    default: false
  }
});

const realPower = computed(() => props.value?.realPower);
const realPowerStr = computed(() => powerStr(realPower.value));
const apparentPower = computed(() => props.value?.apparentPower);
const apparentPowerStr = computed(() => powerStr(apparentPower.value));
const reactivePower = computed(() => props.value?.reactivePower);
const reactivePowerStr = computed(() => powerStr(reactivePower.value));
const powerFactor = computed(() => props.value?.powerFactor);
const powerFactorStr = computed(() => {
  if (powerFactor.value === undefined) return '';
  return powerFactor.value.toFixed(2);
});

/**
 * @param {number} val
 * @return {string}
 */
function powerStr(val) {
  if (val === undefined) return '';
  return `${(val / 1000).toFixed(3)} kW`;
}
</script>

<style scoped>
.v-list-item {
  min-height: auto;
}
</style>
