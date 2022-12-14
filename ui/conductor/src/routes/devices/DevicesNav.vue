<template>
  <v-navigation-drawer app clipped color="transparent" floating>
    <v-list dense nav>
      <v-list-item to="/devices/lighting">
        <v-list-item-icon>
          <v-icon>mdi-lightbulb</v-icon>
        </v-list-item-icon>
        <v-list-item-content>Lighting</v-list-item-content>
      </v-list-item>
    </v-list>
  </v-navigation-drawer>
</template>

<script>
export default {
  name: 'AdminNav'
}
</script>

<style scoped>
::v-deep(.v-list-item--active) {
  color: var(--v-primary-base);
}
</style>
