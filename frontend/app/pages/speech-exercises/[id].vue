<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'

interface Segment {
  id: number
  title: string
  exercise_id: number
  spoken_text: string
  audio_file: string
}

const route = useRoute()
const exerciseId = Number(route.params.id)
const exerciseTitle = route.query.title as string || `Exercise ${exerciseId}`

const segments = ref<Segment[]>([])
const currentSegmentIndex = ref(0)
const loading = ref(true)
const error = ref<string | null>(null)

const currentSegment = computed(() => segments.value[currentSegmentIndex.value])
const totalSegments = computed(() => segments.value.length)
const isLastSegment = computed(() => currentSegmentIndex.value === totalSegments.value - 1)

const nextSegment = () => {
  if (!isLastSegment.value) {
    currentSegmentIndex.value++
  }
}

const prevSegment = () => {
  if (currentSegmentIndex.value > 0) {
    currentSegmentIndex.value--
  }
}

onMounted(async () => {
  try {
    const response = await fetch(`/api/exercises/${exerciseId}/segments`)
    if (!response.ok) throw new Error('Failed to fetch segments')
    const data = await response.json()
    segments.value = data.segments
    console.log("Total segments fetched:", segments.value.length);
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Unknown error'
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="container mx-auto p-4">
    <UButton to="/speech-exercises" class="mb-4">Back to Exercises</UButton>

    <h1 class="text-3xl font-bold mb-4">{{ exerciseTitle }}</h1>
    
    <div v-if="loading" class="text-center py-8">
      <p>Loading segments...</p>
    </div>

    <div v-else-if="error" class="text-red-500 py-4">
      <p>Error: {{ error }}</p>
    </div>

    <div v-else-if="segments.length === 0" class="text-gray-500 py-4">
      <p>No segments found for this exercise.</p>
    </div>

    <div v-else>
      <div class="mb-4 text-sm text-gray-600">
        Segment {{ currentSegmentIndex + 1 }} of {{ totalSegments }}
      </div>

      <SpeechSegment v-if="currentSegment" :segment="currentSegment" />

      <SpeechAnswer v-if="currentSegment" :segment="currentSegment" :exerciseId="exerciseId" />

      <div class="flex gap-2 mt-6">
        <UButton @click="prevSegment" :disabled="currentSegmentIndex === 0">
          Previous
        </UButton>
        <UButton @click="nextSegment" :disabled="isLastSegment">
          Next
        </UButton>
      </div>
    </div>
  </div>
</template>