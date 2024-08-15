<template>
  <v-tooltip v-if="!acked" location="left" transition="slide-x-reverse-transition" color="neutral-lighten-4">
    <template #activator="{ props }">
      <v-btn
          v-bind="props"
          icon
          v-if="!acked"
          color="warning"
          :disabled="blockActions"
          size="small"
          @click.stop="$emit('acknowledge')">
        <v-icon>mdi-circle-outline</v-icon>
      </v-btn>
    </template>
    Acknowledge this notification
  </v-tooltip>
  <v-menu v-else location="bottom left" offset-y>
    <template #activator="{ props }">
      <v-avatar v-bind="props" color="neutral lighten-8" class="text--black" size="21">
        <template v-if="hasAuthor">{{ authorInitials }}</template>
        <v-icon v-else color="black" size="small">mdi-check</v-icon>
      </v-avatar>
    </template>
    <v-card min-width="300">
      <v-card-title>
        Acknowledged
        <v-spacer/>
        <v-icon end color="secondary">mdi-check</v-icon>
      </v-card-title>
      <v-card-subtitle>{{ ackTimeStr }}</v-card-subtitle>
      <v-card-text>
        <template v-if="hasAuthor">
          <template v-if="hasAuthorName"> By: {{ authorName }}<br></template>
          <template v-if="hasAuthorEmail"> Mail: {{ authorEmail }}<br></template>
        </template>
        <template v-else> Anonymous acknowledgement</template>
      </v-card-text>
      <v-card-actions>
        <v-btn @click="$emit('unacknowledge')" variant="text" block color="error" :disabled="blockActions">
          <v-icon start>mdi-close</v-icon>
          Clear Acknowledgement
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-menu>
</template>

<script setup>
import {timestampToDate} from '@/api/convpb';
import useAuthSetup from '@/composables/useAuthSetup';
import {computed} from 'vue';

const {blockActions} = useAuthSetup();

const props = defineProps({
  ack: {
    type: [Object],
    default: null
  }
});
defineEmits(['acknowledge', 'unacknowledge']);

/**
 * @param {Alert.Acknowledgement.Author.AsObject|undefined} author
 * @return {string}
 */
function authorToInitials(author) {
  if (!author) return '-';
  const name = author.displayName;
  if (name) {
    const names = name.trim().split(' ');
    if (names.length === 0) return '?';
    if (names.length === 1) return names[0][0].toUpperCase();
    return names[0][0].toUpperCase() + names[names.length - 1][0].toUpperCase();
  }
  const email = author.email;
  if (email) {
    const subject = email.substring(0, email.indexOf('@'));
    const names = subject.split(/[\s._-]/g);
    if (names.length === 0) return '?';
    if (names.length === 1) return names[0][0].toUpperCase();
    return names[0][0].toUpperCase() + names[names.length - 1][0].toUpperCase();
  }

  return '-';
}

const acked = computed(() => Boolean(props.ack?.acknowledgeTime));
const ackTimeStr = computed(() => timestampToDate(props.ack.acknowledgeTime).toLocaleString());
const hasAuthor = computed(() => Boolean(props.ack?.author));
const hasAuthorName = computed(() => Boolean(props.ack?.author?.displayName));
const hasAuthorEmail = computed(() => Boolean(props.ack?.author?.email));

const authorName = computed(() => props.ack?.author?.displayName);
const authorEmail = computed(() => props.ack?.author?.email);
const authorInitials = computed(() => authorToInitials(props.ack?.author));
</script>

<style scoped>
.v-avatar {
  font-size: 12px;
  color: var(--v-neutral-base);
}
</style>
