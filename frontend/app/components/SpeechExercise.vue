<template>
  <UCard class="speech-exercise">
    <div class="flex items-center justify-between">
      <div>
        <h2>{{ exercise.title }}</h2>
        <p class="text-sm text-gray-600 mt-1">{{ exercise.description }}</p>
      </div>
      <UButton icon="i-heroicons-chevron-right-20-solid" size="lg" variant="ghost" @click="startExercise">
        Start
      </UButton>
    </div>
  </UCard>
</template>

<script setup lang="ts">
import type { Exercise } from "~/types/models";

const router = useRouter();

const props = defineProps<{
  exercise: Exercise;
}>();

const startExercise = () => {
  console.log('Starting exercise with id:', props.exercise.id);
  router.push({
    path: `/speech-exercises/${props.exercise.id}`,
    query: { title: props.exercise.title }
  });
};
</script>
