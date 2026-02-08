<template>
  <main class="flex flex-col items-center justify-center min-h-[70vh] text-center px-4">
    <h1 class="text-3xl font-bold mb-4">Speech Shadowing</h1>
    <p class="max-w-md text-gray-600 mb-2">
      A small tool to practice pronunciation.
      Listen to a sentence, repeat it, and see how close you got.
    </p>
    <p class="max-w-md text-gray-500 text-sm mb-8">
      Based on the
      <a href="https://en.wikipedia.org/wiki/Speech_shadowing" target="_blank" class="underline hover:text-gray-700">speech shadowing</a>
      technique — basically repeating what you hear, as you hear it.
    </p>
    <UButton to="/speech-exercises" size="lg">Get started</UButton>
  </main>
</template>
