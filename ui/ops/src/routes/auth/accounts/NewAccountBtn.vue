<template>
  <v-btn>
    New Account...
    <v-menu v-model="menuVisible" activator="parent" :close-on-content-click="false">
      <new-account-card :name="props.name" @save="onSave" @cancel="onCancel" ref="cardRef"/>
    </v-menu>
  </v-btn>
</template>

<script setup>
import NewAccountCard from '@/routes/auth/accounts/NewAccountCard.vue';
import {ref, watch} from 'vue';

const props = defineProps({
  name: {
    type: String,
    default: undefined
  }
});
const emit = defineEmits(['save', 'cancel']);

const menuVisible = ref(false);

const onSave = (e) => {
  emit('save', e);
  menuVisible.value = false;
};
const onCancel = () => {
  emit('cancel');
  menuVisible.value = false;
}
// reset the new account form once the menu is hidden
const cardRef = ref(null);
watch(menuVisible, (value) => {
  if (value) {
    setTimeout(() => {
      const form = cardRef.value;
      if (!form) return;
      form.reset();
    }, 250);
  }
});
</script>

<style scoped>

</style>