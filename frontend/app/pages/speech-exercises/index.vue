<script setup lang="ts">
import SpeechExercise from "~/components/SpeechExercise.vue";
import type { Exercise } from "~/types/models";

const exercises = ref<Exercise[]>([]);
const loading = ref(true);

onMounted(async () => {
  try {
    const response = await fetch("/api/exercises");
    const data = await response.json();
    exercises.value = data.exercises || [];
  } catch (error) {
    console.error("Failed to fetch exercises:", error);
  } finally {
    loading.value = false;
  }
});
</script>

<template>
  <div class="container mx-auto p-4">
    <h1 class="text-xl font-bold mb-4">Choose an exercise</h1>

    <div v-if="loading">Loading exercises...</div>
    <div v-else>
      <SpeechExercise
        v-for="exercise in exercises"
        :key="exercise.id"
        :exercise="exercise"
      />
    </div>
  </div>
</template>