<template>
  <div :class="[props.statusBarColor, 'status-bar d-flex flex-row rounded-t']" style="height: 44px">
    <v-card-title class="d-block ext-uppercase text-body-1 font-weight-medium pt-3 text-truncate" style="width: 350px">
      {{ props.title }}
    </v-card-title>
    <v-spacer/>
    <!-- <v-card-title class="text-body-1 font-weight-bold text-uppercase white--text pt-3">
    </v-card-title> -->
    <v-tooltip v-if="showClose" bottom>
      <template #activator="{ on, attrs }">
        <v-btn class="elevation-0 mr-1 mt-1" icon v-bind="attrs" v-on="on" @click="emit('click:close')">
          <v-icon color="white">mdi-close</v-icon>
        </v-btn>
      </template>
      <span>Close</span>
    </v-tooltip>
  </div>
</template>

<script setup>
const props = defineProps({
  title: {
    type: String,
    default: ''
  },
  statusBarColor: {
    type: String,
    default: ''
  },
  showClose: {
    type: Boolean,
    default: false
  }
});
const emit = defineEmits(['click:close']);
</script>
<style lang="scss" scoped>
.status-bar {
  transition: background-color 0.5s ease-in-out;
}
</style>
