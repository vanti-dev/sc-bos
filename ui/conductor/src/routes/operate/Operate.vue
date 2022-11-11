<template>
  <div class="mt-4 mx-6">
    <v-layout row wrap>
      <v-flex xs12 sm4 md3>
        <SideBar />
      </v-flex>
      <v-flex xs12 sm8 md9>
        <v-layout column>
          <Header />
          <BulkAction />
          <Table />
        </v-layout>
      </v-flex>
    </v-layout>
  </div>
</template>

<script></script>

<style lang="scss" scoped></style>
