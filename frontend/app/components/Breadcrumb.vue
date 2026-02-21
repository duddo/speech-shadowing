<template>
  <nav class="flex items-center gap-2 text-sm">
    <template v-for="(link, index) in links" :key="index">
      <span v-if="index > 0" class="text-gray-400 mx-2">›</span>
      <NuxtLink 
        v-if="link.to"
        :to="link.to"
        class="text-blue-600 hover:underline"
      >
        {{ link.label }}
      </NuxtLink>
      <span v-else class="text-gray-700 font-semibold">
        {{ link.label }}
      </span>
    </template>
  </nav>
</template>

<script setup lang="ts">
defineProps<{
  links: Array<{ label: string; to?: string }>
}>()
</script>
