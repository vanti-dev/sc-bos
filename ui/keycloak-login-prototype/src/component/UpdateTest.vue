<script setup>
import {ref} from "vue";
import {TestApiPromiseClient} from "@ew-auth-poc/ui-gen/src/test_grpc_web_pb";
import {Test, UpdateTestRequest} from "@ew-auth-poc/ui-gen/src/test_pb";

const props = defineProps(["token"])

const data = ref("")
const error = ref(null)

const client = new TestApiPromiseClient("http://localhost:8000", null, null)

async function updateTest() {
  const token = props.token
  if (token === undefined || token === null) {
    window.alert("get token first")
    return
  }

  const request = new UpdateTestRequest()
      .setTest(new Test().setData(data.value))

  try {
    await client.updateTest(request, {
      "Authorization": "Bearer " + token,
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
  <textarea v-model="data"></textarea>
  <div>
    <button @click="updateTest">TestApi.UpdateTest</button>
  </div>
  <div v-if="error !== null">
    <h3>Error</h3>
    <pre>{{ error }}</pre>
  </div>
</template>