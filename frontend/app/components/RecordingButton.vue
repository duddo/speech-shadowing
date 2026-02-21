<template>
  <div class="flex gap-4 items-center">
    <button
      @click="$emit('click')"
      :disabled="disabled"
      class="flex-shrink-0 w-16 h-16 rounded-full flex items-center justify-center
       disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
      :class="isRecording 
        ? 'bg-red-500 hover:bg-red-600 text-white' 
        : 'bg-green-500 hover:bg-green-600 text-white'"
      :title="isRecording ? 'Stop Recording' : 'Start Recording'"
    >
      <UIcon v-if="!isRecording" name="i-heroicons-microphone-solid" class="w-6 h-6" />
      <UIcon v-else name="i-heroicons-stop-solid" class="w-6 h-6" />
    </button>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  isRecording: boolean
  disabled: boolean
}>()

defineEmits<{
  click: []
}>()
</script>
