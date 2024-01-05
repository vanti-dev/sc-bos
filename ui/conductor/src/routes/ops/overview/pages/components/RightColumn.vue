<template>
  <div>
    <presence-card
        v-if="props.item.traits.showOccupancy"
        class="mb-5"
        :name="props.item.traits.showOccupancy"/>

    <environmental-card
        v-if="environmentalValues"
        class="mt-3"
        gauge-color="#ffc432"
        :name="environmentalValues.indoor"
        :external-name="environmentalValues.outdoor"/>

    <content-card
        v-if="props.item.traits.showEnergyConsumption"
        class="pb-0"
        style="min-height:385px;">
      <v-card-title class="text-h4 pl-4">Energy Consumption</v-card-title>
      <EnergyGraph
          classes="mt-n2 ml-n2 mr-1"
          color="#ffc432"
          color-middle="rgba(255, 196, 50, 0.35)"
          :hide-legends="true"
          :metered="props.item.traits.showEnergyConsumption"/>
    </content-card>
  </div>
</template>

<script setup>
import {computed} from 'vue';
import ContentCard from '@/components/ContentCard.vue';

import EnvironmentalCard from '@/routes/ops/overview/pages/widgets/environmental/EnvironmentalCard.vue';
import PresenceCard from '@/routes/ops/overview/pages/widgets/occupancy/PresenceCard.vue';
import EnergyGraph from '@/routes/ops/overview/pages/widgets/energyAndDemand/EnergyGraph.vue';

const props = defineProps({
  item: {
    type: Object,
    default: () => ({})
  }
});

const traits = computed(() => {
  return props.item.traits;
});

const environmentalValues = computed(() => {
  // Extracting indoor and outdoor values, defaulting to undefined if not present
  const indoor = traits.value.showEnvironment?.indoor;
  const outdoor = props.item.traits.showEnvironment?.outdoor;

  // Function to handle the value conversion
  const handleValue = (value) => {
    // If value is false or undefined, return an empty string
    if (value === false || value === undefined) {
      return '';
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

</script>
