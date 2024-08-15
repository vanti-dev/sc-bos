<template>
  <v-list id="error-view" tile>
    <v-list-item
        v-for="(item, index) in errorStore.errors"
        :key="index"
        style="height: 64px;">
      <v-snackbar
          :model-value="true"
          timeout="5000"
          color="error"
          absolute>
        <span class="error-name">{{ item.name }}</span>
        {{ statusCodeToString(item.source.code) }}: {{ item.source.message }}
        <template #actions="attrs">
          <v-btn variant="text" v-bind="attrs" @click="errorStore.clearError(item)">
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

.error-name {
  display: block;
  font-size: 0.8em;
}
</style>
