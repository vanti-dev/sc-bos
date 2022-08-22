<template>
  <v-navigation-drawer app clipped floating color="transparent">
    <v-list>
      <v-list-item to="/admin/summary">Summary</v-list-item>
      <v-list-item to="/admin/tenants">Tenants</v-list-item>
    </v-list>
  </v-navigation-drawer>
</template>

<script>
export default {
  name: 'AdminNav'
}
</script>

<style scoped>

</style>
