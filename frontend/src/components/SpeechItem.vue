<script setup lang="ts">
import { ref } from 'vue'
import {type RecordedAudio, type RecorderData, sendRecording, startRecording, stopRecording} from "@/recorder.ts";

enum State { Init, Recording, Sending, Stop }

let state = ref<State>(State.Init);
let buttonTag = ref("speaking");
let progress = ref("...");
let outcome = ref("");
let audioUrl = ref("");

let recorder : RecorderData | null;

if (navigator.mediaDevices === undefined) {
  alert("Your browser is not supported!");
}

async function onStart() {
  recorder = await startRecording(onDataAvailable);
  state.value = State.Recording;
  buttonTag.value = "again";
}

async function onStop() {
  if (recorder === null)
    return;

  const audio : RecordedAudio = stopRecording(recorder);
  recorder = null;
  audioUrl.value = audio.audioUrl;

  state.value = State.Sending;

  outcome.value = await sendRecording(audio);

  state.value = State.Stop;
}

function onDataAvailable() {
  progress.value = progress.value + ".";
  if (progress.value.length > 3)
      progress.value = ".";
}
</script>

<template>
  <div>

    <button @click="onStart"
            v-if="state !== State.Recording && state !== State.Sending">
            Start {{ buttonTag }}</button>

    <p v-if="state === State.Recording">
      Recording{{ progress }}</p>

    <button @click="onStop"
            v-if="state === State.Recording" >
            Stop</button>

    <audio v-if="state === State.Stop"
           :src="audioUrl"
           controls />

    <p v-if="state === State.Sending">Converting speech to text...</p>

    <p v-if="state === State.Stop" >{{ outcome }}</p>
  </div>
</template>

<style scoped>
div {
  border: 1px solid red;
  display: grid;
  grid-template-columns: auto;
  padding: 5px;
}

p {
  padding: 5px;
}

button {
  margin: 5px;
  padding: 7px;
}

audio {
  margin: 5px;
}
</style>