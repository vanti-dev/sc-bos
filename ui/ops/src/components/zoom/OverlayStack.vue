<template>
  <div class="stack">
    <slot/>
  </div>
</template>

<style scoped>
.stack {
  display: grid;
  align-items: stretch;
  justify-items: stretch;
}

.stack > * {
  grid-row: 1;
  grid-column: 1;
}
</style>
