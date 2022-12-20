<template>
  <v-list class="pa-0" dense nav>
    <v-list-item to="/auth/users">
      <v-list-item-icon>
        <v-icon>mdi-account-group</v-icon>
      </v-list-item-icon>
      <v-list-item-content>Users</v-list-item-content>
    </v-list-item>
    <v-list-item to="/auth/third-party">
      <v-list-item-icon>
        <v-icon>mdi-key</v-icon>
      </v-list-item-icon>
      <v-list-item-content>Third-party accounts</v-list-item-content>
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
