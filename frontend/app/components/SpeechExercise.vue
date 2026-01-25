<!--
  SpeechExercise

  Represents a single speech exercise component.
  The user can listen to a speech sample, then the components allow them to record
  their own attempt. It sends the audio to a backend service to compare with original.

-->

<template>
  <UCard class="speech-exercise">
    <div class="flex items-center justify-between">
      <h2>{{ title }}</h2>
      <UButton icon="i-heroicons-chevron-right-20-solid" size="sm" variant="ghost" @click="startExercise">
        Start
      </UButton>
    </div>
    <!-- Additional content will go here -->
  </UCard>
</template>

<script setup lang="ts">
const router = useRouter();

const props = defineProps<{
  title: string;
  id: number;
}>();

const startExercise = () => {
  console.log('Starting exercise with id:', props.id);
  router.push({
    path: `/speech-exercises/${props.id}`,
    query: { title: props.title }
  });
};
</script>
