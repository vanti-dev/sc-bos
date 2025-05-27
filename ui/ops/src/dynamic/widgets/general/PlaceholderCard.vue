<template>
  <v-card disabled/>
</template>

<style scoped>
.v-card {
  --c1: rgba(var(--v-theme-surface), .4);
  --c2: rgba(var(--v-theme-error), 0.2);
  background-color: var(--c2);
  background-image: linear-gradient(45deg,
  var(--c1) 45%,
  var(--c2) 45%, var(--c2) 50%,
  var(--c1) 50%, var(--c1) 95%,
  var(--c2) 95%, var(--c2) 100%);
  background-repeat: repeat;
  background-attachment: fixed;
  background-size: 40px 40px;
}
</style>
