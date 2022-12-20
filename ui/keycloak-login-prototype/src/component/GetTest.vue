<script setup>
import {TestApiPromiseClient} from '@sc-bos/ui-gen/src/test_grpc_web_pb'
import {GetTestRequest} from '@sc-bos/ui-gen/src/test_pb';
import {reactive} from 'vue'

const props = defineProps(['token'])
const client = new TestApiPromiseClient('http://localhost:8000', null, null)

const result = reactive({
  data: null,
  error: null
})

async function getTest() {
  const token = props.token
  if (token === null) {
    window.alert("get token first")
    return
  }

  result.data = null

  try {
    result.data = await client.getTest(new GetTestRequest(),{
      "Authorization": "Bearer " + props.token
    })
    result.error = null
  } catch (e) {
    result.error = e
  }
}
</script>

<template>
  <h2>Get Test</h2>
  <button @click="getTest">TestApi.GetTest</button>
  <div v-if="result.data !== null">
    <h3>Response</h3>
    <pre>{{ result.data }}</pre>
  </div>
  <div v-if="result.error !== null">
    <h3>Error</h3>
    <pre>{{ result.error }}</pre>
  </div>
</template>
