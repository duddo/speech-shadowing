<template>
  <div ref="waveformContainer" class="flex-grow"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import WaveSurfer from 'wavesurfer.js'

const props = defineProps<{
  audioFile: string | Blob
}>()

const emit = defineEmits<{
  'update:isPlaying': [value: boolean]
}>()

const waveformContainer = ref<HTMLElement | null>(null)
let wavesurfer: any
let blobUrl: string | null = null

const loadAudio = (source: string | Blob) => {
  if (!wavesurfer) return
  revokeBlobUrl()
  if (source instanceof Blob) {
    blobUrl = URL.createObjectURL(source)
    wavesurfer.load(blobUrl)
  } else {
    wavesurfer.load(source)
  }
}

const revokeBlobUrl = () => {
  if (blobUrl) {
    URL.revokeObjectURL(blobUrl)
    blobUrl = null
  }
}

const playPause = () => {
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
  
  loadAudio(props.audioFile)
  
  wavesurfer.on('play', () => {
    emit('update:isPlaying', true)
  })
  
  wavesurfer.on('pause', () => {
    emit('update:isPlaying', false)
  })
})

watch(() => props.audioFile, (newAudioFile) => {
  if (wavesurfer && newAudioFile) {
    loadAudio(newAudioFile)
    emit('update:isPlaying', false)
  }
})

onBeforeUnmount(() => {
  revokeBlobUrl()
  wavesurfer?.destroy()
})

defineExpose({
  playPause
})
</script>
