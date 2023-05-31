<template>
  <v-card elevation="0" tile>
    <WithLighting :name="props.name">
      <template #lighting="{lightingData}">
        <v-list tile class="ma-0 pa-0">
          <v-subheader class="text-title-caps-large neutral--text text--lighten-3">Lighting</v-subheader>
          <v-list-item class="py-1">
            <v-list-item-title class="text-body-small text-capitalize">Brightness</v-list-item-title>
            <v-list-item-subtitle class="text-capitalize">{{ lightingData.brightness }}%</v-list-item-subtitle>
          </v-list-item>
        </v-list>
        <v-progress-linear
            height="34"
            class="mx-4 my-2"
            :value="lightingData.brightness"
            background-color="neutral lighten-1"
            color="accent"/>
        <v-card-actions class="px-4">
          <v-btn small color="neutral lighten-1" elevation="0" @click="lightingData.updateLight(100)">On</v-btn>
          <v-btn small color="neutral lighten-1" elevation="0" @click="lightingData.updateLight(0)">Off</v-btn>
          <v-spacer/>
          <v-btn
              small
              color="neutral lighten-1"
              elevation="0"
              @click="lightingData.updateLight(lightingData.brightness+1)"
              :disabled="lightingData.brightness >= 100">
            Up
          </v-btn>
          <v-btn
              small
              color="neutral lighten-1"
              elevation="0"
              @click="lightingData.updateLight(lightingData.brightness-1)"
              :disabled="lightingData.brightness <= 0">
            Down
          </v-btn>
        </v-card-actions>
        <v-progress-linear color="primary" indeterminate :active="lightingData.updateValue.loading"/>
      </template>
    </WithLighting>
  </v-card>
</template>

<script setup>
import WithLighting from '../renderless/WithLighting.vue';

const props = defineProps({
  // unique name of the device
  name: {
    type: String,
    default: ''
  }
});

</script>

<style scoped>
.v-list-item {
  min-height: auto;
}
.v-progress-linear {
  width: auto;
}
</style>
