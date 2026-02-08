<template>
  <button
    @click="$emit('click')"
    class="flex-shrink-0 w-16 h-16 rounded-full bg-blue-500 hover:bg-blue-600
           text-white flex items-center justify-center transition-colors"
    :title="isPlaying ? 'Pause' : 'Play'" >
        <UIcon v-if="!isPlaying" name="i-heroicons-play-solid" class="w-6 h-6" />
        <UIcon v-else name="i-heroicons-pause-solid" class="w-6 h-6" />
  </button>
</template>

<script setup lang="ts">
defineProps<{
  isPlaying: boolean
}>()

defineEmits<{
  click: []
}>()
</script>
