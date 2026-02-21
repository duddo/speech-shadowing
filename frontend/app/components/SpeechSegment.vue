<template>
  <UCard>
    <template #header>
      <h2>{{ segment.title }}</h2>
    </template>

    <div class="flex gap-4 items-center">
      <PlayPauseButton
        :is-playing="isPlaying"
        @click="togglePlayPause"
      />
      <AudioWaveform
        ref="audioWaveformRef"
        :audio-file="segment.audio_file"
        v-model:is-playing="isPlaying"
      />
    </div>

  </UCard>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { Segment } from '@/types/models'

const props = defineProps<{
  segment: Segment
}>()

const isPlaying = ref(false)
const audioWaveformRef = ref<any>(null)

const togglePlayPause = () => {
  audioWaveformRef.value?.playPause()
}
</script>