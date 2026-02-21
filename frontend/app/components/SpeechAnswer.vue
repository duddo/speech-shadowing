<template>
  <UCard class="mt-6">
    <template #header>
      <div v-if="state === State.Idle">
        <h3>Record what you hear:</h3>
      </div>
      <div v-else-if="state === State.Recording">
        <h3>Recording...</h3>
      </div>
      <div v-else-if="state === State.Submitting">
        <h3>Submitting...</h3>
      </div>
      <div v-else-if="state === State.Result">
        <h3>Done</h3>
      </div>
      <div v-else-if="state === State.Error">
        <h3>Error</h3>
      </div>
    </template>

    <div class="space-y-4">
      <!-- Idle / Recording -->
      <div v-if="state === State.Idle || state === State.Recording" class="flex gap-4 items-center justify-center">
        <RecordingButton
          :is-recording="state === State.Recording"
          :disabled="false"
          @click="toggleRecording"
        />
        <div v-if="state === State.Recording" class="text-red-500 font-semibold animate-pulse">
          Recording... {{ recordingTime.toFixed(1) }}s
        </div>
      </div>

      <!-- Submitting / Result / Error: show waveform -->
      <template v-if="state === State.Submitting || state === State.Result || state === State.Error">
        
        <div class="flex gap-4 items-center">
          <PlayPauseButton
            :is-playing="isPlaying"
            @click="waveformRef?.playPause()"
          />
          <AudioWaveform
            ref="waveformRef"
            :audio-file="audioBlob!"
            v-model:is-playing="isPlaying"
          />
        </div>

        <!-- Submitting spinner -->
        <div v-if="state === State.Submitting" class="flex items-center gap-2 text-sm text-gray-500">
          <UIcon name="i-heroicons-arrow-path" class="w-4 h-4 animate-spin" />
          Analyzing...
        </div>

        <!-- Result -->
        <AnswerResult
          v-if="state === State.Result && submitResult"
          :rating="submitResult.rating"
          :transcript="submitResult.transcript"
        />

        <!-- Error -->
        <div v-if="state === State.Error" class="p-4 rounded-lg bg-red-50 border border-red-200">
          <p class="text-red-700 text-sm">{{ errorMessage }}</p>
        </div>

        <!-- Try again button (for both Result and Error states) -->
        <UButton
          v-if="state === State.Result || state === State.Error"
          @click="reset"
          icon="i-heroicons-arrow-path"
          variant="soft"
          size="lg"
          block
        >
          Try again
        </UButton>
      </template>
    </div>
  </UCard>
</template>

<script setup lang="ts">
import { ref, onBeforeUnmount, watch } from 'vue'
import type { SegmentAnswer, Segment } from '@/types/models'

enum State {
  Idle = 'idle',
  Recording = 'recording',
  Submitting = 'submitting',
  Result = 'result',
  Error = 'error',
}

const props = defineProps<{
  segment: Segment
  exerciseId: number
}>()

const state = ref<State>(State.Idle)
const recordingTime = ref(0)
const audioBlob = ref<Blob | null>(null)
const isPlaying = ref(false)
const submitResult = ref<SegmentAnswer | null>(null)
const errorMessage = ref('')
const waveformRef = ref<any>(null)

let mediaRecorder: MediaRecorder | null = null
let recordingInterval: ReturnType<typeof setInterval> | null = null

watch(() => props.segment.id, () => {
  reset()
})

const toggleRecording = async () => {
  if (state.value === State.Recording) {
    stopRecording()
  } else {
    await startRecording()
  }
}

const startRecording = async () => {
  try {
    recordingTime.value = 0
    audioBlob.value = null
    submitResult.value = null

    const stream = await navigator.mediaDevices.getUserMedia({ audio: true })
    mediaRecorder = new MediaRecorder(stream)

    const chunks: BlobPart[] = []
    mediaRecorder.ondataavailable = (e) => chunks.push(e.data)
    mediaRecorder.onstop = () => {
      audioBlob.value = new Blob(chunks, { type: 'audio/webm' })
      submitAnswer()
    }

    mediaRecorder.start()
    state.value = State.Recording

    recordingInterval = setInterval(() => {
      recordingTime.value += 0.1
    }, 100)
  } catch (err) {
    errorMessage.value = err instanceof Error ? err.message : 'Failed to access microphone'
    state.value = State.Error
  }
}

const stopRecording = () => {
  if (recordingInterval) {
    clearInterval(recordingInterval)
    recordingInterval = null
  }
  if (mediaRecorder) {
    mediaRecorder.stop()
    mediaRecorder.stream.getTracks().forEach(track => track.stop())
  }
}

const submitAnswer = async () => {
  if (!audioBlob.value) return

  state.value = State.Submitting

  try {
    const formData = new FormData()
    formData.append('audio', audioBlob.value, 'answer.webm')
    formData.append('payload', JSON.stringify({
      exercise_id: props.exerciseId,
      segment_id: props.segment.id,
      audio_format: 'webm'
    }))

    const response = await fetch('/api/do-exercise', {
      method: 'POST',
      body: formData
    })

    const data = await response.json()

    if (response.ok) {
      submitResult.value = data.answer
      state.value = State.Result
    } else {
      errorMessage.value = data.error || 'Failed to submit'
      state.value = State.Error
    }
  } catch (err) {
    errorMessage.value = err instanceof Error ? err.message : 'Network error'
    state.value = State.Error
  }
}

const reset = () => {
  audioBlob.value = null
  submitResult.value = null
  errorMessage.value = ''
  isPlaying.value = false
  state.value = State.Idle
}

onBeforeUnmount(() => {
  stopRecording()
})
</script>
