<template>
  <UCard>
    <template #header>
      <h2>{{ segment.title }}</h2>
    </template>

    <div class="flex gap-4 items-center">
      <button
        @click="togglePlayPause"
        class="flex-shrink-0 w-16 h-16 rounded-full bg-blue-500 hover:bg-blue-600 text-white flex items-center justify-center transition-colors"
        :title="isPlaying ? 'Pause' : 'Play'"
      >
        <UIcon v-if="!isPlaying" name="i-heroicons-play-solid" class="w-6 h-6" />
        <UIcon v-else name="i-heroicons-pause-solid" class="w-6 h-6" />
      </button>
      <div ref="waveformContainer" class="flex-grow"></div>
    </div>

  </UCard>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import WaveSurfer from 'wavesurfer.js'

interface Segment {
  id: number
  title: string
  exercise_id: number
  spoken_text: string
  audio_file: string
}

const props = defineProps<{
  segment: Segment
}>()

const waveformContainer = ref<HTMLElement | null>(null)
const isPlaying = ref(false)
let wavesurfer: any

const togglePlayPause = () => {
  if (wavesurfer) {
    wavesurfer.playPause()
  }
}

onMounted(() => {
  if (!waveformContainer.value) 
    return
  wavesurfer = WaveSurfer.create({
    container: waveformContainer.value,
    waveColor: '#4f46e5',
    progressColor: '#818cf8',
    height: 96
  })
  
  wavesurfer.load(props.segment.audio_file)
  
  wavesurfer.on('play', () => {
    isPlaying.value = true
  })
  
  wavesurfer.on('pause', () => {
    isPlaying.value = false
  })
})

watch(() => props.segment.audio_file, (newAudioFile) => {
  if (wavesurfer && newAudioFile) {
    wavesurfer.load(newAudioFile)
    isPlaying.value = false
  }
})

onBeforeUnmount(() => {
  wavesurfer?.destroy()
})
</script>