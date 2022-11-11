<template>
  <div class="mt-4">
    <v-layout row wrap>
      <v-flex xs12 sm4 md3>
        <SideBar />
      </v-flex>
      <v-flex xs12 sm8 md9>
        <Table />
      </v-flex>
    </v-layout>
  </div>
</template>

<script>
import Table from "@/components/Operate/Table.vue";
import SideBar from "@/components/Operate/SideBar.vue";
</script>

<style lang="scss" scoped></style>
