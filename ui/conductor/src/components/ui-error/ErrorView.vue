<template>
  <v-list id="error-view" tile>
    <v-list-item
        v-for="(item, index) in errorStore.errors"
        :key="index"
        class="mt-2">
      <v-snackbar
          :value="true"
          timeout="5000"
          color="error"
          absolute>
        {{ statusCodeToString(item.source.code) }}: {{ item.source.message }}
        <template #action="{ attrs }">
          <v-btn text v-bind="attrs" @click="errorStore.clearError(item)">
            Dismiss
          </v-btn>
        </template>
      </v-snackbar>
    </v-list-item>
  </v-list>
</template>

<script setup>
import {useErrorStore} from '@/components/ui-error/error';
import {statusCodeToString} from '@/components/ui-error/util';

const errorStore = useErrorStore();

</script>

<style scoped>
#error-view {
  position: fixed;
  bottom: 0;
  left: 50%;
  background: transparent;
}
</style>
