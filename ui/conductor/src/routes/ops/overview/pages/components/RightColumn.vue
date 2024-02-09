<template>
  <div>
    <presence-card
        v-if="showWidget('showOccupancy')"
        class="mb-5"
        :name="widgets.showOccupancy"/>

    <environmental-card
        v-if="showWidget('showEnvironment')"
        class="mt-3"
        gauge-color="#ffc432"
        :name="environmentalValues.indoor"
        :external-name="environmentalValues.outdoor"/>

    <content-card
        v-if="showWidget('showEnergyConsumption')"
        class="pb-0"
        style="min-height:385px;">
      <EnergyGraph
          chart-title="Energy Consumption"
          classes="pb-2 mr-1"
          color="#ffc432"
          color-middle="rgba(255, 196, 50, 0.35)"
          :hide-legends="true"
          :metered="widgets.showEnergyConsumption"/>
    </content-card>
  </div>
</template>

<script setup>
import ContentCard from '@/components/ContentCard.vue';
import EnergyGraph from '@/routes/ops/overview/pages/widgets/energyAndDemand/EnergyGraph.vue';
import EnvironmentalCard from '@/routes/ops/overview/pages/widgets/environmental/EnvironmentalCard.vue';
import PresenceCard from '@/routes/ops/overview/pages/widgets/occupancy/PresenceCard.vue';
import {computed} from 'vue';

const props = defineProps({
  item: {
    type: Object,
    default: () => ({})
  }
});

const widgets = computed(() => {
  return props.item.widgets;
});

const environmentalValues = computed(() => {
  // Extracting indoor and outdoor values, defaulting to undefined if not present
  const indoor = widgets.value?.showEnvironment?.indoor;
  const outdoor = widgets.value?.showEnvironment?.outdoor;

  // Function to handle the value conversion
  const handleValue = (value) => {
    // If value is false or undefined, do not return anything
    if (value === false || value === undefined) {
      return;
    }
    // Otherwise, return the value as it is (which should be a string)
    return value;
  };

  // Applying handleValue function to indoor and outdoor
  const processedIndoor = handleValue(indoor);
  const processedOutdoor = handleValue(outdoor);

  // Building the return object based on the processed values
  return {
    indoor: processedIndoor,
    outdoor: processedOutdoor
  };
});

/**
 * Checks if a config has a trait enabled
 *
 * @param {string} trait
 * @return {boolean}
 */
const showWidget = (trait) => {
  return widgets.value[trait] !== false && widgets.value[trait] !== undefined;
};
</script>
