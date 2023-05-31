<template>
  <v-card elevation="0" tile>
    <WithLighting :name="props.name">
      <v-list tile class="ma-0 pa-0">
        <v-subheader class="text-title-caps-large neutral--text text--lighten-3">Lighting</v-subheader>
        <v-list-item class="py-1">
          <v-list-item-title class="text-body-small text-capitalize">Brightness</v-list-item-title>
          <v-list-item-subtitle class="text-capitalize">{{ props.brightness }}%</v-list-item-subtitle>
        </v-list-item>
      </v-list>
      <v-progress-linear
          height="34"
          class="mx-4 my-2"
          :value="props.brightness"
          background-color="neutral lighten-1"
          color="accent"/>
      <v-card-actions class="px-4">
        <v-btn small color="neutral lighten-1" elevation="0" @click="props.updateLight(100)">On</v-btn>
        <v-btn small color="neutral lighten-1" elevation="0" @click="props.updateLight(0)">Off</v-btn>
        <v-spacer/>
        <v-btn
            small
            color="neutral lighten-1"
            elevation="0"
            @click="props.updateLight(props.brightness+1)"
            :disabled="props.brightness >= 100">
          Up
        </v-btn>
        <v-btn
            small
            color="neutral lighten-1"
            elevation="0"
            @click="props.updateLight(props.brightness-1)"
            :disabled="props.brightness <= 0">
          Down
        </v-btn>
      </v-card-actions>
      <v-progress-linear color="primary" indeterminate :active="value.updateValue.loading"/>
    </WithLighting>
  </v-card>
</template>

<script setup>
import WithLighting from '../renderless/WithLighting.vue';

const props = defineProps({
  brightness: {
    type: [Number, String],
    default: ''
  },
  // unique name of the device
  name: {
    type: String,
    default: ''
  },
  updateLight: {
    type: Function,
    default: () => {}
  },
  updateValue: {
    type: Object,
    default: () => {}
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
