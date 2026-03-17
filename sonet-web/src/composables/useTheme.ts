import { ref, watch } from 'vue'

const theme = ref<'light' | 'dark'>(
  (localStorage.getItem('sonote-theme') as 'light' | 'dark') || 'light'
)

function applyTheme(t: 'light' | 'dark') {
  if (t === 'dark') {
    document.documentElement.classList.add('dark')
  } else {
    document.documentElement.classList.remove('dark')
  }
}

// Apply on load
applyTheme(theme.value)

watch(theme, (newTheme) => {
  localStorage.setItem('sonote-theme', newTheme)
  applyTheme(newTheme)
})

export function useTheme() {
  const toggleTheme = () => {
    theme.value = theme.value === 'light' ? 'dark' : 'light'
  }

  return {
    theme,
    toggleTheme,
  }
}
