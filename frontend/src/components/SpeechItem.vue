<script setup lang="ts">
import { ref } from 'vue'

let showStart = ref(true);
let showStop = ref(false);

function startSpeaking() {
  showStop.value = true;
  showStart.value = false;
}

function stopSpeaking() {
  showStart.value = true;
  showStop.value= false;
}
</script>

<template>
  <div>
    <button @click="startSpeaking"
            :hidden="!showStart">
            Start speaking</button>
    <button @click="stopSpeaking"
            :hidden="!showStop" >
            Stop</button>
    <audio id="audio"
           hidden="hidden" controls />
    <p id="result"
       hidden="hidden" />
    <button id="retry"
            hidden="hidden" @click="startSpeaking" >
            Retry</button>
  </div>
</template>

<style scoped>
div {
  border: 1px solid red;
  display: grid;
  grid-template-columns: auto;
  padding: 5px;
}

button {
  margin: 5px;
}

audio {
  margin: 5px;
}
</style>