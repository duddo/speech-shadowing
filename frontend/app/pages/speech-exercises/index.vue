<script setup lang="ts">
import SpeechExercise from "~/components/SpeechExercise.vue";

interface Exercise {
  id: number;
  title: string;
}

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
    <UButton to="/" class="mb-4">Home</UButton>

    <h1 class="text-3xl font-bold mb-4">Speech exercises</h1>

    <div v-if="loading">Loading exercises...</div>
    <div v-else>
      <SpeechExercise
        v-for="exercise in exercises"
        :key="exercise.id"
        :title="exercise.title"
        :id="exercise.id"
      />
    </div>
  </div>
</template>