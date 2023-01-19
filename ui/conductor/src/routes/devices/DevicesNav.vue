<template>
  <v-list class="pa-0" dense nav>
    <v-list-item to="/devices/all">
      <v-list-item-icon>
        <v-icon>mdi-view-list</v-icon>
      </v-list-item-icon>
      <v-list-item-content>All</v-list-item-content>
    </v-list-item>
    <v-list-item to="/devices/lighting">
      <v-list-item-icon>
        <v-icon>mdi-lightbulb</v-icon>
      </v-list-item-icon>
      <v-list-item-content>Lighting</v-list-item-content>
    </v-list-item>
    <v-list-item to="/devices/hvac">
      <v-list-item-icon>
        <v-icon>mdi-thermometer</v-icon>
      </v-list-item-icon>
      <v-list-item-content>HVAC</v-list-item-content>
    </v-list-item>
  </v-list>
</template>

<script>
export default {
  name: 'AdminNav'
};
</script>

<style scoped>
::v-deep(.v-list-item--active) {
  color: var(--v-primary-base);
}
</style>
