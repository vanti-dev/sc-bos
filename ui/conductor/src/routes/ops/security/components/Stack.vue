<template>
  <div class="stack">
    <slot/>
  </div>
</template>

<style scoped>
.stack {
  display: grid;
}

.stack > * {
  grid-row: 1;
  grid-column: 1;
}
</style>
