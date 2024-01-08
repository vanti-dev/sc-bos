<script setup>
import {OnOffApiPromiseClient} from '@smart-core-os/sc-api-grpc-web/traits/on_off_grpc_web_pb';
import {OnOff, UpdateOnOffRequest} from '@smart-core-os/sc-api-grpc-web/traits/on_off_pb';
import {ref} from 'vue';

const props = defineProps(['token'])

const data = ref('')
const error = ref(null)

const client = new OnOffApiPromiseClient('http://localhost:8000', null, null)

async function updateTest() {
  const token = props.token
  if (token === undefined || token === null) {
    window.alert('get token first')
    return
  }

  const request = new UpdateOnOffRequest()
      .setOnOff(new OnOff().setState(OnOff.State.ON))

  try {
    await client.updateOnOff(request, {
      'Authorization': 'Bearer ' + token,
    })
    error.value = null
  } catch (e) {
    error.value = e
  }
}

</script>

<template>
  <h2>Update Test</h2>
  <h3>Value to Write</h3>
  <label>
    Turn On
    <input v-model="data" type="checkbox" value="false">
  </label>
  <div>
    <button @click="updateTest">OnOffApi.UpdateOnOff</button>
  </div>
  <div v-if="error !== null">
    <h3>Error</h3>
    <pre>{{ error }}</pre>
  </div>
</template>
