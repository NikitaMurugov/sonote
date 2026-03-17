<template>
  <div class="min-h-screen flex items-center justify-center bg-bg-base">
    <div class="text-center anim-fade-up">
      <h1 class="text-7xl font-bold text-text-tertiary/30 mb-4" style="font-family: var(--font-heading)">404</h1>
      <p class="text-text-secondary mb-6 text-sm">Страница не найдена</p>
      <router-link
        to="/"
        class="text-primary hover:text-primary-hover text-sm font-medium underline underline-offset-4 decoration-primary/30 hover:decoration-primary transition-colors"
      >
        На главную
      </router-link>
    </div>
  </div>
</template>
