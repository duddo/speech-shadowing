<template>
  <div class="space-y-2 p-4 rounded-lg bg-green-50 border border-green-200">
    <div class="flex justify-center items-center gap-2">
      <span class="text-sm text-green-700">Accuracy:</span>
      <span class="text-lg font-bold text-green-600">{{ (rating * 100).toFixed(1) }}%</span>
    </div>
    <p class="text-sm text-green-700">
      <strong>Transcript:</strong> {{ transcript }}
    </p>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  rating: number
  transcript: string
}>()
</script>
