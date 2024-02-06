<template>
  <content-card class="mb-8 d-flex flex-column pt-6 pl-6 pr-8">
    <h4 class="text-h4 py-0 pt-lg-2 pb-4 pl-4">Power</h4>
    <energy-graph
        :generated="props.generated"
        :metered="props.metered"/>
    <v-row class="d-flex flex-row justify-center mt-10 mb-1 ml-12">
      <v-col cols="auto">
        <v-row>
          <v-col cols="auto" class="text-h1 align-self-center" style="line-height: 0.3em;">
            <WithElectricDemand
                v-slot="{resource}"
                :name="props.generated">
              {{
                Math.abs(storeEnergyValues('generated', resource.value))
              }}<span style="font-size: 0.5em;">kW</span><br>
            </WithElectricDemand>
            <span class="pl-1 text-title orange--text" style="line-height: 0.35em;">Generated</span>
          </v-col>
          <v-col
              cols="1"
              class="text-h1 d-flex flex-row justify-space-around"
              style="line-height: 0.75em;">
            +
          </v-col>

          <v-col cols="auto" class="text-h1 align-self-center" style="line-height: 0.3em;">
            <WithElectricDemand
                v-slot="{resource}"
                :name="props.metered">
              {{ storeEnergyValues('metered', resource.value) }}<span style="font-size: 0.5em;">kW</span><br>
            </WithElectricDemand>
            <span class="pl-1 text-title primary--text" style="line-height: 0.35em;">Metered</span>
          </v-col>

          <v-col
              cols="1"
              class="text-h1 d-flex flex-row justify-space-around"
              style="line-height: 0.75em;">
            =
          </v-col>

          <v-col cols="auto" class="text-h1 align-self-center" style="line-height: 0.3em;">
            {{ energy.total }}<span style="font-size: 0.5em;">kW</span><br>
            <span class="pl-1 text-title" style="line-height: 0.35em;">Total</span>
          </v-col>
        </v-row>
      </v-col>
      <v-divider class="mr-4" vertical/>
      <v-col cols="auto">
        <v-row>
          <v-col cols="auto" class="text-h1 align-self-center" style="line-height: 0.3em;">
            {{
              energyIntensity
            }}<span style="font-size: 0.5em;">kW/person</span><br>
            <span class="text-title orange--text" style="line-height: 0.35em;">Energy Intensity</span>
          </v-col>
        </v-row>
      </v-col>
    </v-row>
  </content-card>
</template>

<script setup>
import ContentCard from '@/components/ContentCard.vue';
import useOccupancyTrait from '@/composables/traits/useOccupancyTrait';
import WithElectricDemand from '@/routes/devices/components/renderless/WithElectricDemand.vue';
import EnergyGraph from '@/routes/ops/overview/pages/widgets/energyAndDemand/EnergyGraph.vue';
import {computed, reactive} from 'vue';

const props = defineProps({
  generated: {
    type: String,
    default: ''
  },
  metered: {
    type: String,
    default: 'building'
  }
});

const energy = reactive({
  generated: 0,
  metered: 0,
  total: computed(() => (energy.metered - energy.generated).toFixed(2))
});
const {occupancyValue} = useOccupancyTrait({name: props.metered});

const energyIntensity = computed(() => {
  if (!occupancyValue?.value?.peopleCount) {
    return Math.abs(energy.total).toFixed(2);
  }

  return Math.abs(energy.total / occupancyValue?.value?.peopleCount).toFixed(2);
});

/**
 *
 * @param {string} type
 * @param {number} value
 * @return {number}
 */
function storeEnergyValues(type, value) {
  if (value) {
    energy[type] = Number((value.realPower / 1000).toFixed(2));
  }

  return energy[type];
}
</script>

<style scoped>

</style>
