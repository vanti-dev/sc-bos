<template>
  <v-navigation-drawer app clipped floating color="transparent">
    <v-list>
      <v-list-item to="/admin/summary">
        <v-list-item-content>
          Overview
        </v-list-item-content>
      </v-list-item>
      <v-subheader class="text-uppercase">Access</v-subheader>
      <v-list-item to="/admin/tenants">
        <v-list-item-icon><v-icon>mdi-account-details</v-icon></v-list-item-icon>
        <v-list-item-content>Tenants</v-list-item-content>
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
