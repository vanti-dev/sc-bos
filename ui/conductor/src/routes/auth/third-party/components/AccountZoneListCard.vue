<template>
  <v-card class="pa-0" flat tile>
    <v-subheader class="text-title-caps-large neutral--text text--lighten-3">
      Zones<v-spacer/>
      <v-btn outlined>Edit<v-icon right>mdi-pencil</v-icon></v-btn>
    </v-subheader>
    <v-card-text v-if="zoneList.length > 0" class="pt-1">
      {{ zoneList.join(', ') }}
    </v-card-text>
    <v-card-text v-else class="pt-1">
      <span class="font-italic">No zones associated with this account</span>
    </v-card-text>
  </v-card>
</template>

<script setup>
defineProps({
  zoneList: {
    type: Array,
    default: () => []
  }
});
</script>

<style scoped>

</style>
