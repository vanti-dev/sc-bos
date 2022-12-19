<template>
  <v-dialog v-bind="$attrs" v-on="$listeners" max-width="500">
    <v-card>
      <v-card-title class="align-baseline justify-space-between flex-nowrap text-body-1 mb-2">
        <span>Are you sure you want to delete this secret?</span>
        <v-btn small icon @click="rollback" class="mr-n2">
          <v-icon small>mdi-close</v-icon>
        </v-btn>
      </v-card-title>

      <v-alert type="warning" text tile class="text-body-2">
        Any applications or scripts using this secret will no
        longer be able to access
        the building API. You cannot undo this action.
      </v-alert>

      <v-card-actions>
        <v-btn block color="error" text class="text-none" @click="commit">I understand, delete this secret</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup>
const emit = defineEmits(['commit', 'rollback']);

/**
 *
 */
function commit() {
  emit('commit');
}

/**
 *
 */
function rollback() {
  emit('rollback');
}
</script>

<style scoped>

</style>
